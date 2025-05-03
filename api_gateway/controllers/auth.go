package controllers

import (
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
	resp, err := c.AuthClient.Login(c.Ctx.Request.Context(), loginReq)
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

	// Set the token cookie
	c.Ctx.SetCookie("token", resp.Token, 86400, "/", "", false, true)

	// Check if it's an AJAX request
	if c.Ctx.Input.IsAjax() {
		c.Data["json"] = map[string]string{
			"message":  "Login successful",
			"token":    resp.Token,
			"redirect": "/orders",
		}
		c.ServeJSON()
		return
	}

	// For non-AJAX requests, redirect directly
	c.Redirect("/orders", 302)
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

	_, err = c.AuthClient.Register(c.Ctx.Request.Context(), registerReq)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error registering user: " + err.Error()}
		c.ServeJSON()
		return
	}

	// Successful registration response
	c.Ctx.Output.SetStatus(201) // Created status
	c.Data["json"] = map[string]string{
		"message":  "User registered successfully",
		"redirect": "/auth/login", // Provide redirect URL
	}
	c.ServeJSON()
}

func (c *AuthController) Logout() {
	// Clear the token cookie
	c.Ctx.SetCookie("token", "", -1, "/", "", false, true)
	c.DelSession("token")
	c.Redirect("/auth/login", 302)
}

// ValidateToken checks if the token is valid
func (c *AuthController) ValidateToken() {
	// Get token from header or cookie
	token := c.Ctx.Input.Header("Authorization")
	if token != "" && strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	} else {
		token = c.Ctx.GetCookie("token")
	}

	if token == "" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "No token provided"}
		c.ServeJSON()
		return
	}

	// Validate token with auth service
	validateReq := &auth_service.ValidateTokenRequest{
		Token: token,
	}

	resp, err := c.AuthClient.ValidateToken(c.Ctx.Request.Context(), validateReq)
	if err != nil || !resp.Success {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Token is valid
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{
		"valid":   true,
		"user_id": resp.UserId,
	}
	c.ServeJSON()
}
