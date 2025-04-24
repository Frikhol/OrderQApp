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

func (c *OrderController) CreateOrder() {
	location := c.GetString("location")
	startTimeStr := c.GetString("start_time")
	timeBuffer, _ := c.GetInt("time_buffer")

	if location == "" || startTimeStr == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Location and start time are required"}
		c.ServeJSON()
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid start time format"}
		c.ServeJSON()
		return
	}

	userID := c.GetSession("user_id").(int)
	o := orm.NewOrm()

	order := models.Order{
		Client:        &models.User{Id: userID},
		Location:      location,
		Status:        "pending",
		StartTime:     startTime,
		QueuePosition: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Calculate price based on time buffer
	order.Price = calculatePrice(timeBuffer)

	_, err = o.Insert(&order)
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
