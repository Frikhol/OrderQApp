package routers

import (
	"api_gateway/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Root route
	web.Router("/", &controllers.GatewayController{}, "get:GetIndex")
}
