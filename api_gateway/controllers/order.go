package controllers

import (
	"api_gateway/proto/order_service"
	"encoding/json"
	"time"

	"github.com/beego/beego/v2/server/web"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderController struct {
	web.Controller
	OrderClient order_service.OrderServiceClient
}

func (c *OrderController) GetOrdersPage() {
	c.TplName = "orders.tpl"
}

func (c *OrderController) GetOrdersList() {
	userID := c.Ctx.Input.GetData("user_id").(string)

	orders, err := c.OrderClient.GetUserOrders(c.Ctx.Request.Context(), &order_service.GetUserOrdersRequest{
		UserId: userID,
	})
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = orders
	c.ServeJSON()
}

func (c *OrderController) CreateOrder() {
	// Define a struct to unmarshal the JSON data with string timestamps
	type OrderRequest struct {
		UserId        string `json:"user_id"`
		OrderAddress  string `json:"order_address"`
		OrderLocation string `json:"order_location"`
		OrderDate     string `json:"order_date"`
		OrderTimeGap  string `json:"order_time_gap"`
	}

	var jsonReq OrderRequest

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &jsonReq)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Create the protobuf request
	req := &order_service.CreateOrderRequest{
		UserId:        c.Ctx.Input.GetData("user_id").(string),
		OrderAddress:  jsonReq.OrderAddress,
		OrderLocation: jsonReq.OrderLocation,
	}

	// Parse and convert the order_date string to a timestamppb.Timestamp
	if jsonReq.OrderDate != "" {
		orderDate, err := time.Parse(time.RFC3339, jsonReq.OrderDate)
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Invalid order_date format: " + err.Error()}
			c.ServeJSON()
			return
		}
		req.OrderDate = timestamppb.New(orderDate)
	}

	// Parse and convert the order_time_gap string to a durationpb.Duration
	if jsonReq.OrderTimeGap != "" {
		orderTimeGap, err := time.ParseDuration(jsonReq.OrderTimeGap)
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Invalid order_time_gap format: " + err.Error()}
			c.ServeJSON()
			return
		}
		req.OrderTimeGap = durationpb.New(orderTimeGap)
	}

	_, err = c.OrderClient.CreateOrder(c.Ctx.Request.Context(), req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"success": "Order created successfully"}
	c.ServeJSON()
}

func (c *OrderController) GetOrderById() {
	orderID := c.Ctx.Input.Param(":id")

	order, err := c.OrderClient.GetOrderById(c.Ctx.Request.Context(), &order_service.GetOrderByIdRequest{
		OrderId: orderID,
	})
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = order
	c.ServeJSON()
}

func (c *OrderController) CancelOrder() {
	orderID := c.Ctx.Input.Param(":id")

	_, err := c.OrderClient.CancelOrder(c.Ctx.Request.Context(), &order_service.CancelOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"success": "Order cancelled successfully"}
	c.ServeJSON()
}

func (c *OrderController) CompleteOrder() {
	orderID := c.Ctx.Input.Param(":id")

	_, err := c.OrderClient.CompleteOrder(c.Ctx.Request.Context(), &order_service.CompleteOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"success": "Order completed successfully"}
	c.ServeJSON()
}
