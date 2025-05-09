package controllers

import (
	"api_gateway/proto/order_service"

	"github.com/beego/beego/v2/server/web"
)

type AgentController struct {
	web.Controller
	OrderClient order_service.OrderServiceClient
	//AgentClient agent_service.AgentServiceClient
}

func (c *AgentController) GetOrdersPage() {
	c.TplName = "agent_orders.tpl"
}

func (c *AgentController) StartSearch() {
	//TODO: implement
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
