package main

import (
	"log"
	_ "orderq/web/routers" // Import routers to initialize them

	"github.com/beego/beego/v2/server/web"
)

func main() {
	// Initialize Beego web server
	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.AutoRender = true

	// Template configuration
	web.BConfig.WebConfig.ViewsPath = "web/views"
	web.BConfig.WebConfig.TemplateLeft = "{{"
	web.BConfig.WebConfig.TemplateRight = "}}"

	// Enable debug mode
	web.BConfig.RunMode = "dev"

	// Start the server
	log.Println("Starting API Gateway on port 8080")
	log.Println("Views path:", web.BConfig.WebConfig.ViewsPath)
	web.Run(":8080")
}
