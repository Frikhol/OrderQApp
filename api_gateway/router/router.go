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
	web.Router("/auth/logout", &controllers.AuthController{AuthClient: authClient}, "get:Logout")

	// Auth post routes
	web.Router("/auth/login", &controllers.AuthController{AuthClient: authClient}, "post:Login")
	web.Router("/auth/register", &controllers.AuthController{AuthClient: authClient}, "post:Register")

	// Order web routes - protected with JWT authentication
	web.InsertFilter("/orders", web.BeforeRouter, middleware.JWTAuthMiddleware(authClient))
	web.Router("/orders", &controllers.OrderController{OrderClient: orderClient}, "get:GetOrdersPage")

	web.InsertFilter("/agent/orders", web.BeforeRouter, middleware.JWTAuthMiddleware(authClient))
	web.Router("/agent/orders", &controllers.AgentController{OrderClient: orderClient}, "get:GetOrdersPage")

	// Order API routes - protected with JWT authentication
	web.InsertFilter("/api/orders/*", web.BeforeRouter, middleware.JWTAuthMiddleware(authClient))
	web.Router("/api/orders/create", &controllers.OrderController{OrderClient: orderClient}, "post:CreateOrder")
	web.Router("/api/orders/list", &controllers.OrderController{OrderClient: orderClient}, "get:GetOrdersList")
	web.Router("/api/orders/:id", &controllers.OrderController{OrderClient: orderClient}, "get:GetOrderById")
	web.Router("/api/orders/:id/cancel", &controllers.OrderController{OrderClient: orderClient}, "post:CancelOrder")
	web.Router("/api/orders/:id/complete", &controllers.OrderController{OrderClient: orderClient}, "post:CompleteOrder")

	web.InsertFilter("/api/agent/*", web.BeforeRouter, middleware.JWTAuthMiddleware(authClient))
	web.Router("/api/orders/start_search", &controllers.AgentController{OrderClient: orderClient}, "post:StartSearch")
	web.Router("/api/orders/stop_search", &controllers.AgentController{OrderClient: orderClient}, "post:StopSearch")
	web.Router("/api/orders/:id/accept", &controllers.AgentController{OrderClient: orderClient}, "post:AcceptOrder")
	web.Router("/api/orders/:id/decline", &controllers.AgentController{OrderClient: orderClient}, "post:DeclineOrder")
	web.Router("/api/orders/:id/join", &controllers.AgentController{OrderClient: orderClient}, "post:JoinOrderQueue")
}
