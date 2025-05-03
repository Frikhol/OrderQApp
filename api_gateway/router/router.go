package router

import (
	"api_gateway/controllers"
	"api_gateway/middleware"
	"api_gateway/proto/auth_service"
	"api_gateway/proto/order_service"

	"github.com/beego/beego/v2/server/web"
)

func InitRoutes(authClient auth_service.AuthServiceClient, orderClient order_service.OrderServiceClient) {
	// Root route
	web.Router("/", &controllers.GatewayController{}, "get:GetIndex")

	// Auth get routes
	web.Router("/auth/login", &controllers.AuthController{AuthClient: authClient}, "get:GetLoginPage")
	web.Router("/auth/register", &controllers.AuthController{AuthClient: authClient}, "get:GetRegisterPage")

	// Auth post routes
	web.Router("/auth/login", &controllers.AuthController{AuthClient: authClient}, "post:Login")
	web.Router("/auth/register", &controllers.AuthController{AuthClient: authClient}, "post:Register")
	web.Router("/auth/validate", &controllers.AuthController{AuthClient: authClient}, "get:ValidateToken")
	web.Router("/auth/logout", &controllers.AuthController{AuthClient: authClient}, "get:Logout")

	// Order get routes - allow page access but protect API endpoints
	web.InsertFilter("/orders", web.BeforeRouter, middleware.JWTAuthMiddleware(authClient))
	web.Router("/orders", &controllers.OrderController{OrderClient: orderClient}, "get:GetOrdersPage")

	// Order API routes - protected with JWT authentication
	web.InsertFilter("/api/orders*", web.BeforeRouter, middleware.JWTAuthMiddleware(authClient))
	web.Router("/api/orders/create", &controllers.OrderController{OrderClient: orderClient}, "post:CreateOrder")
}
