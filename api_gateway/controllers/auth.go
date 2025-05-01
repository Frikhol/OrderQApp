package controllers

import (
	"context"
	"encoding/json"
	"strings"

	"api_gateway/proto/auth_service"

	"github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	web.Controller
	AuthClient auth_service.AuthServiceClient
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) GetLoginPage() {
	c.TplName = "login.tpl"
}

func (c *AuthController) GetRegisterPage() {
	c.TplName = "register.tpl"
}

func (c *AuthController) Login() {
	var req AuthRequest

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid JSON request"}
		c.ServeJSON()
		return
	}

	// Validate email and password
	if req.Email == "" || req.Password == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Email and password are required"}
		c.ServeJSON()
		return
	}

	// Create login request for auth service
	loginReq := &auth_service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call the auth service
	resp, err := c.AuthClient.Login(context.Background(), loginReq)
	if err != nil {
		// Log the full error to help with debugging
		c.Ctx.Output.SetStatus(401)
		errorMsg := err.Error()
		c.Data["json"] = map[string]string{
			"error": "Authentication failed: " + errorMsg,
			"email": req.Email,
		}
		c.ServeJSON()
		return
	}

	tokenString := resp.Token

	// Set session
	c.SetSession("user_email", req.Email)
	c.SetSession("is_authenticated", true)
	c.SetSession("token", tokenString) // Also store the token in the session

	c.Data["json"] = map[string]interface{}{
		"token":      tokenString,
		"user_email": req.Email,
	}
	c.ServeJSON()
}

func (c *AuthController) Register() {
	var req AuthRequest

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid JSON request"}
		c.ServeJSON()
		return
	}

	// Validate email format
	if req.Email == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Email is required"}
		c.ServeJSON()
		return
	}

	// Basic email format validation
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid email format"}
		c.ServeJSON()
		return
	}

	// Validate password
	if req.Password == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Password is required"}
		c.ServeJSON()
		return
	}

	// Remove the password hashing here - let the auth service handle it
	registerReq := &auth_service.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := c.AuthClient.Register(context.Background(), registerReq)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error registering user: " + err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"message": resp.Message}
	c.ServeJSON()
}

func (c *AuthController) IsLoggedIn() {
	token := c.GetSession("token")

	req := auth_service.ValidateTokenRequest{Token: token.(string)}

	resp, err := c.AuthClient.ValidateToken(context.Background(), &req)
	if !resp.Success || err != nil {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Invalid token"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"Message": resp.Message,
	}
	c.ServeJSON()
}

func (c *AuthController) Logout() {
	c.DelSession("token")
	c.Redirect("/auth/login", 302)
}
