package controllers

import (
	"orderqapp/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type OrderController struct {
	beego.Controller
}

type CreateOrderRequest struct {
	Location   string    `json:"location"`
	StartTime  time.Time `json:"start_time"`
	TimeBuffer int       `json:"time_buffer"` // in minutes
}

func (c *OrderController) CreateOrder() {
	var req CreateOrderRequest
	if err := c.ParseForm(&req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid request"}
		c.ServeJSON()
		return
	}

	userID := c.GetSession("user_id").(int)
	o := orm.NewOrm()

	order := models.Order{
		Client:        &models.User{Id: userID},
		Location:      req.Location,
		Status:        "pending",
		StartTime:     req.StartTime,
		QueuePosition: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Calculate price based on time, location, and other factors
	order.Price = calculatePrice(req.TimeBuffer)

	_, err := o.Insert(&order)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error creating order"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = order
	c.ServeJSON()
}

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

func calculatePrice(timeBuffer int) float64 {
	// Basic price calculation logic
	basePrice := 10.0
	timeMultiplier := float64(timeBuffer) / 30.0 // 30 minutes as base unit
	weatherMultiplier := 1.0                     // This would be fetched from a weather API in production

	return basePrice * timeMultiplier * weatherMultiplier
}
