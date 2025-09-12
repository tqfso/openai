package main

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.rightchaintech.com/openai/openserver/logger"
)

func main() {

	logger.Init("info", "logs/app.log")

	logger.Info("Application started")	

	defer logger.Sync()

	// GIN日志使用自定义logger
	gin.DefaultWriter = logger.GinWriter()

	r := gin.New()
	r.Use(GinLogger(), GinRecovery())
	r.Run(":8080")
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
