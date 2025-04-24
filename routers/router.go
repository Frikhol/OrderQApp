package routers

import (
	"orderqapp/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// Root route
	beego.Router("/", &controllers.AuthController{}, "get:LoginPage")

	// Auth routes
	beego.Router("/auth/login", &controllers.AuthController{}, "get:LoginPage")
	beego.Router("/auth/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/auth/register", &controllers.AuthController{}, "get:RegisterPage")
	beego.Router("/auth/register", &controllers.AuthController{}, "post:Register")
	beego.Router("/auth/logout", &controllers.AuthController{}, "get:Logout")

	// Dashboard route
	beego.Router("/dashboard", &controllers.DashboardController{}, "get:Index")

	// Order routes (protected)
	beego.Router("/orders/new", &controllers.OrderController{}, "post:CreateOrder")
	beego.Router("/orders/:id/position", &controllers.OrderController{}, "put:UpdateQueuePosition")

	// Add authentication middleware to protected routes
	beego.InsertFilter("/dashboard", beego.BeforeRouter, controllers.AuthMiddleware)
	beego.InsertFilter("/orders", beego.BeforeRouter, controllers.AuthMiddleware)
	beego.InsertFilter("/orders/:id/position", beego.BeforeRouter, controllers.AuthMiddleware)
}
