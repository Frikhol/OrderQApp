package controllers

import (
	"encoding/json"
	"orderqapp/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type OrderController struct {
	beego.Controller
}

// @Title GetOrders
// @Description Get all orders for the current user
// @Success 200 {object} []models.Order
// @Failure 500 Internal server error
// @router /api/orders [get]
func (c *OrderController) GetOrders() {
	userID := c.GetSession("user_id").(int)
	o := orm.NewOrm()

	var orders []models.Order
	_, err := o.QueryTable("orders").
		Filter("client_id", userID).
		OrderBy("-created_at").
		All(&orders)

	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error getting orders"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = orders
	c.ServeJSON()
}

// @Title GetCurrentOrder
// @Description Get the current active order for the user
// @Success 200 {object} models.Order
// @Failure 404 Not found
// @Failure 500 Internal server error
// @router /api/orders/current [get]
func (c *OrderController) GetCurrentOrder() {
	userID := c.GetSession("user_id").(int)
	o := orm.NewOrm()

	var order models.Order
	err := o.QueryTable("orders").
		Filter("client_id", userID).
		Filter("status__in", "pending", "in_progress").
		OrderBy("-created_at").
		One(&order)

	if err == orm.ErrNoRows {
		beego.Info("No active order found for user:", userID)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "No active order found"}
		c.ServeJSON()
		return
	}

	if err != nil {
		beego.Error("Error getting current order:", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error getting current order: " + err.Error()}
		c.ServeJSON()
		return
	}

	// Load the client data
	if err := o.Read(order.Client); err != nil {
		beego.Error("Error loading client:", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error loading client data: " + err.Error()}
		c.ServeJSON()
		return
	}

	beego.Info("Current order found:", order)
	c.Data["json"] = order
	c.ServeJSON()
}

// @Title CreateOrder
// @Description Create a new order
// @Param	body		body 	models.Order	true		"Order data"
// @Success 200 {object} models.Order
// @Failure 400 Bad request
// @Failure 500 Internal server error
// @router /api/orders [post]
func (c *OrderController) CreateOrder() {
	var orderData struct {
		Location   string `json:"location"`
		StartTime  string `json:"start_time"`
		TimeBuffer int    `json:"time_buffer"`
	}

	// Log the raw request body
	beego.Info("Raw request body:", string(c.Ctx.Input.RequestBody))

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &orderData); err != nil {
		beego.Error("Error unmarshaling request body:", err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid request body: " + err.Error()}
		c.ServeJSON()
		return
	}

	beego.Info("Parsed order data:", orderData)

	if orderData.Location == "" || orderData.StartTime == "" {
		beego.Error("Missing required fields:", orderData)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Location and start time are required"}
		c.ServeJSON()
		return
	}

	startTime, err := time.Parse(time.RFC3339, orderData.StartTime)
	if err != nil {
		beego.Error("Error parsing start time:", err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid start time format: " + err.Error()}
		c.ServeJSON()
		return
	}

	userID := c.GetSession("user_id").(int)
	o := orm.NewOrm()

	// Load the user
	user := models.User{Id: userID}
	if err := o.Read(&user); err != nil {
		beego.Error("Error loading user:", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error loading user data: " + err.Error()}
		c.ServeJSON()
		return
	}

	order := models.Order{
		Client:        &user,
		Location:      orderData.Location,
		Status:        "pending",
		StartTime:     startTime,
		QueuePosition: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Calculate price based on time buffer
	order.Price = calculatePrice(orderData.TimeBuffer)

	beego.Info("Creating order:", order)

	_, err = o.Insert(&order)
	if err != nil {
		beego.Error("Error inserting order:", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error creating order: " + err.Error()}
		c.ServeJSON()
		return
	}

	// Load the order with related data
	if err := o.Read(&order); err != nil {
		beego.Error("Error loading order:", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error loading order data: " + err.Error()}
		c.ServeJSON()
		return
	}

	// Load the client data
	if err := o.Read(order.Client); err != nil {
		beego.Error("Error loading client:", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error loading client data: " + err.Error()}
		c.ServeJSON()
		return
	}

	beego.Info("Order created successfully:", order)
	c.Data["json"] = order
	c.ServeJSON()
}

// @Title UpdateQueuePosition
// @Description Update the queue position of an order
// @Param	id		path 	int	true		"Order ID"
// @Param	position		query 	int	true		"New queue position"
// @Success 200 {object} models.Order
// @Failure 404 Not found
// @Failure 500 Internal server error
// @router /api/orders/:id/position [put]
func (c *OrderController) UpdateQueuePosition() {
	orderID, _ := c.GetInt(":id")
	newPosition, _ := c.GetInt("position")

	o := orm.NewOrm()
	order := models.Order{Id: orderID}
	if err := o.Read(&order); err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "Order not found"}
		c.ServeJSON()
		return
	}

	order.QueuePosition = newPosition
	order.UpdatedAt = time.Now()

	queueUpdate := models.QueueUpdate{
		Order:    &order,
		Position: newPosition,
	}

	o.Begin()
	_, err := o.Update(&order)
	if err != nil {
		o.Rollback()
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error updating order"}
		c.ServeJSON()
		return
	}

	_, err = o.Insert(&queueUpdate)
	if err != nil {
		o.Rollback()
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error creating queue update"}
		c.ServeJSON()
		return
	}

	o.Commit()
	c.Data["json"] = order
	c.ServeJSON()
}

// @Title CancelOrder
// @Description Cancel an order
// @Param	id		path 	int	true		"Order ID"
// @Success 200 {object} models.Order
// @Failure 404 Not found
// @Failure 500 Internal server error
// @router /api/orders/:id/cancel [post]
func (c *OrderController) CancelOrder() {
	orderID, _ := c.GetInt(":id")
	o := orm.NewOrm()

	order := models.Order{Id: orderID}
	if err := o.Read(&order); err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "Order not found"}
		c.ServeJSON()
		return
	}

	// Check if the order belongs to the current user
	userID := c.GetSession("user_id").(int)
	if order.Client.Id != userID {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = map[string]string{"error": "Not authorized to cancel this order"}
		c.ServeJSON()
		return
	}

	// Only allow cancellation of pending orders
	if order.Status != "pending" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Only pending orders can be cancelled"}
		c.ServeJSON()
		return
	}

	order.Status = "cancelled"
	order.UpdatedAt = time.Now()

	_, err := o.Update(&order)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error cancelling order"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = order
	c.ServeJSON()
}

func calculatePrice(timeBuffer int) float64 {
	// Basic price calculation logic
	basePrice := 10.0
	timeMultiplier := float64(timeBuffer) / 30.0 // 30 minutes as base unit
	weatherMultiplier := 1.0                     // This would be fetched from a weather API in production

	return basePrice * timeMultiplier * weatherMultiplier
}
