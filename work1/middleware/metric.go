package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"time"
)

var (
	MetricRegistry  *prometheus.Registry
	MetricRegister  prometheus.Registerer
	requestTotal    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
)

const (
	LabelHost      = "host"
	LabelHttpState = "httpState"
	LabelUrl       = "url"
)

func init() {
	MetricRegistry = prometheus.NewRegistry()
	hostname, _ := os.Hostname()
	MetricRegister = prometheus.WrapRegistererWith(prometheus.Labels{LabelHost: hostname}, MetricRegistry)

	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_call_count",
			Help: "total count for http calls",
		}, []string{LabelHttpState, LabelUrl})
	MetricRegister.MustRegister(requestTotal)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_call_duration",
			Help:    "call milliseconds ",
			Buckets: []float64{10, 50, 100, 200, 400, 800, 1500, 2000},
		}, []string{LabelHttpState, LabelUrl})
	MetricRegister.MustRegister(requestDuration)
}

func MetricMiddle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/metrics" {
			ctx.Next()
			return
		}
		t0 := time.Now()
		ctx.Next()

		requestDuration.With(prometheus.Labels{
			LabelHttpState: fmt.Sprintf("%d", ctx.Writer.Status()),
			LabelUrl:       ctx.Request.URL.Path,
		}).Observe(float64(time.Since(t0).Milliseconds()))

		requestTotal.With(prometheus.Labels{
			LabelHttpState: fmt.Sprintf("%d", ctx.Writer.Status()),
			LabelUrl:       ctx.Request.URL.Path,
		}).Inc()
	}
}
