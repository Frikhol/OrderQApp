package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type GatewayController struct {
	web.Controller
}

func (c *GatewayController) GetIndex() {
	c.TplName = "index.tpl"
}
