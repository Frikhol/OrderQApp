package controllers

import (
	"encoding/json"
	"orderqapp/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ErrorLogController struct {
	beego.Controller
}

// @Title LogError
// @Description Log client-side error
// @Param	body		body 	models.ErrorLog	true		"Error log data"
// @Success 200 {object} models.ErrorLog
// @Failure 400 Bad request
// @Failure 500 Internal server error
// @router /logs/error [post]
func (c *ErrorLogController) LogError() {
	var errorLog models.ErrorLog
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &errorLog); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.ServeJSON()
		return
	}

	// Parse timestamp from ISO string
	timestamp, err := time.Parse(time.RFC3339, errorLog.Timestamp.Format(time.RFC3339))
	if err != nil {
		timestamp = time.Now()
	}
	errorLog.Timestamp = timestamp

	o := orm.NewOrm()

	// Insert the error log
	_, err = o.Insert(&errorLog)
	if err != nil {
		beego.Error("Failed to save error log:", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Failed to save error log"}
		c.ServeJSON()
		return
	}

	// Log to Beego's logger
	beego.Error("Client Error:", errorLog.Context, "-", errorLog.Error)
	if errorLog.Stack != "" {
		beego.Error("Stack trace:", errorLog.Stack)
	}
	beego.Error("User Agent:", errorLog.UserAgent)
	beego.Error("URL:", errorLog.Url)

	c.Data["json"] = errorLog
	c.ServeJSON()
}

// @Title GetErrorLogs
// @Description Get error logs with pagination
// @Param	page		query 	int	false		"Page number"
// @Param	limit		query 	int	false		"Number of logs per page"
// @Success 200 {object} []models.ErrorLog
// @Failure 500 Internal server error
// @router /logs/error [get]
func (c *ErrorLogController) GetErrorLogs() {
	page, _ := c.GetInt("page", 1)
	limit, _ := c.GetInt("limit", 20)
	offset := (page - 1) * limit

	o := orm.NewOrm()
	var errorLogs []models.ErrorLog
	_, err := o.QueryTable(new(models.ErrorLog)).
		OrderBy("-created_at").
		Limit(limit, offset).
		All(&errorLogs)

	if err != nil {
		beego.Error("Failed to get error logs:", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Failed to get error logs"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = errorLogs
	c.ServeJSON()
}
