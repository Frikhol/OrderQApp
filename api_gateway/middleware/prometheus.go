package middleware

import (
	"strconv"
	"time"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

// PrometheusFilter returns a middleware function that collects prometheus metrics
func PrometheusFilter() func(ctx *context.Context) {
	return func(ctx *context.Context) {
		start := time.Now()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		method := ctx.Input.Method()
		endpoint := ctx.Input.URI()
		status := strconv.Itoa(ctx.Output.Status)

		httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}
