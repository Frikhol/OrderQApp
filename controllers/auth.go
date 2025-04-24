package controllers

import (
	"orderqapp/models"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
)

type AuthController struct {
	beego.Controller
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (c *AuthController) Login() {
	// Get form values directly
	email := c.GetString("email")
	password := c.GetString("password")

	// Validate email and password
	if email == "" || password == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Email and password are required"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	user := models.User{Email: email}
	if err := o.Read(&user, "Email"); err != nil {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Invalid credentials"}
		c.ServeJSON()
		return
	}

	if !user.CheckPassword(password) {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Invalid credentials"}
		c.ServeJSON()
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("jwtSecret")))
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error generating token"}
		c.ServeJSON()
		return
	}

	// Set session
	c.SetSession("user_id", user.Id)
	c.SetSession("user_role", user.Role)
	c.SetSession("is_authenticated", true)

	c.Data["json"] = map[string]interface{}{
		"token":   tokenString,
		"role":    user.Role,
		"user_id": user.Id,
	}
	c.ServeJSON()
}

func (c *AuthController) Register() {
	var req RegisterRequest

	// Get form values directly
	req.Email = c.GetString("email")
	req.Password = c.GetString("password")
	req.Role = c.GetString("role")

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

	// Validate role
	if req.Role != "client" && req.Role != "agent" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid role. Must be either 'client' or 'agent'"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// Check if user already exists
	existingUser := models.User{Email: req.Email}
	err := o.Read(&existingUser, "Email")
	if err == nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "User with this email already exists"}
		c.ServeJSON()
		return
	}

	user := models.User{
		Email: req.Email,
		Role:  req.Role,
	}

	if err := user.HashPassword(req.Password); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error hashing password"}
		c.ServeJSON()
		return
	}

	_, err = o.Insert(&user)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Error creating user"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"message": "User registered successfully"}
	c.ServeJSON()
}

func (c *AuthController) LoginPage() {
	c.TplName = "auth/login.tpl"
}

func (c *AuthController) RegisterPage() {
	c.TplName = "auth/register.tpl"
}

func (c *AuthController) Logout() {
	// Clear session
	c.DelSession("user_id")
	c.DelSession("user_role")
	c.DelSession("is_authenticated")

	// Redirect to login page
	c.Redirect("/auth/login", 302)
}

// AuthMiddleware проверяет аутентификацию пользователя
func AuthMiddleware(ctx *context.Context) {
	if ctx.Input.Session("is_authenticated") == nil {
		ctx.Redirect(302, "/auth/login")
		return
	}
}
