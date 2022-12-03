package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

// Health 健康检查默认返回 200
func Health(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
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

func main() {
	router := gin.New()
	router.Use(LogMiddleware(), gin.Recovery())

	router.GET("/set-header", SetHeaderHandler)
	router.GET("/healthz", Health)

	_ = router.Run(":8888")
}
