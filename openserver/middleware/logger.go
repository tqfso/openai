package middleware

import (
	"openserver/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpLogger() gin.HandlerFunc {
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
