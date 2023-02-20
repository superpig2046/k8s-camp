package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"k8s-camp/work1/middleware"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	// 初始化日志
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})
	log.SetOutput(os.Stdout)

	// 设置日志级别
	log.SetLevel(log.DebugLevel)
	logLevel, exist := os.LookupEnv("LOG_LEVEL")
	if exist {
		switch logLevel {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		}
	}
}

// Health 健康检查默认返回 200， 并且输出服务的环境参数
func Health(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	// 读取环境配置
	envs := make([]string, 0)
	for _, k := range []string{"VERSION", "LOG_LEVEL", "GIN_MODE"} {
		v, exist := os.LookupEnv(k)
		if exist {
			envs = append(envs, fmt.Sprintf("%s=%s", k, v))
		}
	}

	ctx.JSON(http.StatusOK,
		map[string]string{
			"status": "ok",
			"envs":   strings.Join(envs, ";"),
		})
	return
}

// SetHeaderHandler 把操作系统环境参数写入Header，并且返回原来的Header。返回数据包含所有信息
func SetHeaderHandler(ctx *gin.Context) {
	result := map[string]string{}
	for k, v := range ctx.Request.Header {
		result[k] = strings.Join(v, ";")
	}

	version, exist := os.LookupEnv("VERSION")
	if exist {
		ctx.Header("VERSION", version)
		result["VERSION"] = version
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, result)
	return
}

// ExportWithDebugLog 输出一个Debug日志在console
func ExportWithDebugLog(ctx *gin.Context) {
	log.Debugf("here is some debug log")
	ctx.JSON(http.StatusOK, map[string]string{"result": "ok"})
	return
}

// ExportWithWarnLog 输出一个Warn日志在console
func ExportWithWarnLog(ctx *gin.Context) {
	log.Warnf("here is some debug log")
	ctx.JSON(http.StatusOK, map[string]string{"result": "ok"})
	return
}

// RandomSleep 随机执行
func RandomSleep(ctx *gin.Context) {
	time.Sleep(time.Second * time.Duration(rand.Float64()*2))
	ctx.JSON(http.StatusOK, map[string]string{"result": "ok"})
	return
}

// LogMiddleware 修改Gin默认日志输出格式 IP [时间] 耗时 方法+请求路径
func LogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(p gin.LogFormatterParams) string {
		return fmt.Sprintf("%s\t [%s]|\t%d|\t%s|%s\t %s\n",
			p.ClientIP,
			p.TimeStamp.Format(time.RFC3339),
			p.StatusCode,
			p.Latency,
			p.Method,
			p.Path,
		)
	})
}

func Metric(handler http.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func main() {
	// 从环境中读取 GIN_MODE 通过 configMap 引入到环境变量
	ginMode, exist := os.LookupEnv("GIN_MODE")
	if exist {
		gin.SetMode(ginMode)
	}

	router := gin.New()
	router.Use(LogMiddleware(), middleware.MetricMiddle(), gin.Recovery())

	router.GET("/set-header", SetHeaderHandler)
	router.GET("/healthz", Health)
	router.GET("/debug-log", ExportWithDebugLog)
	router.GET("/warn-log", ExportWithWarnLog)
	// 一份代码两个入口，通过设置Version环境参数不同
	router.GET("/business", Health)
	router.GET("/sale", Health)

	// 随机返回
	router.GET("/random-sleep", RandomSleep)

	// metric
	router.GET("/metrics", Metric(promhttp.InstrumentMetricHandler(
		middleware.MetricRegistry, promhttp.HandlerFor(middleware.MetricRegistry, promhttp.HandlerOpts{}),
	)))

	_ = router.Run(":8888")
}
