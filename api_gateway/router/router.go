package router

import (
	"api_gateway/controllers"
	"api_gateway/proto/auth_service"

	"github.com/beego/beego/v2/server/web"
)

func InitRoutes(authClient auth_service.AuthServiceClient) {
	// Root route
	web.Router("/", &controllers.GatewayController{}, "get:GetIndex")

	web.Router("/auth/login", &controllers.AuthController{AuthClient: authClient}, "get:GetLoginPage")
	web.Router("/auth/register", &controllers.AuthController{AuthClient: authClient}, "get:GetRegisterPage")
	web.Router("/auth/login/check", &controllers.AuthController{AuthClient: authClient}, "get:IsLoggedIn")

	web.Router("/auth/login", &controllers.AuthController{AuthClient: authClient}, "post:Login")
	web.Router("/auth/register", &controllers.AuthController{AuthClient: authClient}, "post:Register")
}
