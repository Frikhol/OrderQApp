package middleware

import (
	"net/http"
	"strings"

	"api_gateway/proto/auth_service"

	"github.com/beego/beego/v2/server/web/context"
)

func JWTAuthMiddleware(authClient auth_service.AuthServiceClient) func(ctx *context.Context) {
	return func(ctx *context.Context) {
		// Get the Authorization header
		authHeader := ctx.Input.Header("Authorization")

		var token string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Otherwise try to get token from cookie
			token = ctx.GetCookie("token")
		}

		if token == "" {
			ctx.Output.SetStatus(http.StatusUnauthorized)
			ctx.Output.Body([]byte("No token provided"))
			return
		}

		// Call the auth service to validate the token
		validateReq := &auth_service.ValidateTokenRequest{
			Token: token,
		}

		// Handle validation gracefully
		validateResp, err := authClient.ValidateToken(ctx.Request.Context(), validateReq)
		if err != nil {
			ctx.Output.SetStatus(http.StatusUnauthorized)
			ctx.Output.Body([]byte(err.Error()))
			return
		}

		if !validateResp.Success {
			ctx.Output.SetStatus(http.StatusUnauthorized)
			ctx.Output.Body([]byte("Invalid or expired token"))
			return
		}

		ctx.Input.SetData("user_id", validateResp.UserId)
	}
}
