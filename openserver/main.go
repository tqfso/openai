package main

import (
	"flag"
	"net/http"
	"runtime/debug"
	"time"

	"openserver/config"
	"openserver/logger"

	"github.com/gin-gonic/gin"
)

func main() {

	// 解析参数
	cfgfile := flag.String("config", "config/config.yaml", "config from file")
	logfile := flag.String("logs", "logs/app.log", "log to file")
	flag.Parse()

	// 初始化日志
	defer logger.Sync()
	logger.Init("info", *logfile)
	logger.Info("Application started")
	gin.DefaultWriter = logger.GinWriter()

	// 初始化配置
	if err := config.Load(*cfgfile); err != nil {
		logger.Error(err.Error())
		return
	}

	// 启动HTTP服务

	r := gin.New()
	r.Use(GinLogger(), GinRecovery())
	r.Run(config.GetServer().ListenAddress())
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // 处理请求

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		logger.Info("HTTP Access",
			logger.String("clientIP", clientIP),
			logger.String("method", method),
			logger.Int("status", statusCode),
			logger.String("path", path),
			logger.Int64("latency_ms", latency.Milliseconds()),
			logger.String("userAgent", c.Request.UserAgent()),
		)
	}
}

func GinRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered",
			logger.Any("error", recovered),
			logger.String("clientIP", c.ClientIP()),
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("stack", string(debug.Stack())), // 记录堆栈
		)

		// 返回友好错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	})
}
