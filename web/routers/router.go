package routers

import (
	"orderq/web/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Root route
	web.Router("/", &controllers.GatewayController{}, "get:GetIndex")
}
