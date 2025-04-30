package controllers

import (
	"log"

	"github.com/beego/beego/v2/server/web"
)

type GatewayController struct {
	web.Controller
}

func (c *GatewayController) GetIndex() {
	log.Println("Handling index request")
	c.TplName = "index.tpl"
	c.Data["title"] = "OrderQ API Gateway"
}
