package controllers

import (
	"github.com/astaxie/beego"
)

type DashboardController struct {
	beego.Controller
}

func (c *DashboardController) Index() {
	// Get user role from session
	role := c.GetSession("user_role").(string)

	// Set template based on role
	if role == "client" {
		c.TplName = "dashboard/client.tpl"
	} else if role == "agent" {
		c.TplName = "dashboard/agent.tpl"
	}

	// Pass user data to template
	c.Data["IsAuthenticated"] = true
	c.Data["Role"] = role
}
