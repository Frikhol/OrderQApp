package controllers

import (
	"api_gateway/proto/order_service"

	"github.com/beego/beego/v2/server/web"
)

type AgentController struct {
	web.Controller
	OrderClient order_service.OrderServiceClient
}

func (c *AgentController) GetOrdersPage() {
	c.TplName = "agent_orders.tpl"
}

func (c *AgentController) StartSearch() {
	orders, err := c.OrderClient.GetAvailableOrders(c.Ctx.Request.Context(), &order_service.GetAvailableOrdersRequest{})
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = orders
	c.ServeJSON()
}

func (c *AgentController) StopSearch() {
	//TODO: implement
}

func (c *AgentController) AcceptOrder() {
	//TODO: implement
}

func (c *AgentController) DeclineOrder() {
	//TODO: implement
}

func (c *AgentController) JoinOrderQueue() {
	//TODO: implement
}
