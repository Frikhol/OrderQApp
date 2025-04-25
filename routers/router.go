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
	beego.Router("/api/orders", &controllers.OrderController{}, "get:GetOrders")
	beego.Router("/api/orders/current", &controllers.OrderController{}, "get:GetCurrentOrder")
	beego.Router("/api/orders", &controllers.OrderController{}, "post:CreateOrder")
	beego.Router("/api/orders/:id/position", &controllers.OrderController{}, "put:UpdateQueuePosition")
	beego.Router("/api/orders/:id/cancel", &controllers.OrderController{}, "post:CancelOrder")

	// Add authentication middleware to protected routes
	beego.InsertFilter("/dashboard", beego.BeforeRouter, controllers.AuthMiddleware)
	beego.InsertFilter("/api/orders", beego.BeforeRouter, controllers.AuthMiddleware)
	beego.InsertFilter("/api/orders/*", beego.BeforeRouter, controllers.AuthMiddleware)

	// Error log routes
	beego.Router("/api/logs/error", &controllers.ErrorLogController{}, "post:LogError")
	beego.Router("/api/logs/error", &controllers.ErrorLogController{}, "get:GetErrorLogs")
}
