package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsController struct {
	web.Controller
}

func (c *MetricsController) GetMetrics() {
	// Serve prometheus metrics
	promhttp.Handler().ServeHTTP(c.Ctx.ResponseWriter, c.Ctx.Request)
}
