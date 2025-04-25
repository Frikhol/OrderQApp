package main

import (
	_ "orderqapp/routers"
	"os"

	"orderqapp/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "orderqapp"
	}

	// Database configuration
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres",
		"host="+dbHost+
			" port="+dbPort+
			" user="+dbUser+
			" password="+dbPassword+
			" dbname="+dbName+
			" sslmode=disable")

	// Register models
	orm.RegisterModel(new(models.User))
	orm.RegisterModel(new(models.Order))
	orm.RegisterModel(new(models.QueueUpdate))
	orm.RegisterModel(new(models.ErrorLog))

	// Create tables
	orm.RunSyncdb("default", false, true)
}

func main() {
	// Development configuration
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// Template configuration
	beego.SetViewsPath("views")
	beego.BConfig.WebConfig.ViewsPath = "views"
	beego.BConfig.WebConfig.TemplateLeft = "{{"
	beego.BConfig.WebConfig.TemplateRight = "}}"

	// Run the application
	beego.Run()
}
