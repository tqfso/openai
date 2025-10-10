package middleware

import (
	"common"
	"common/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		logger.Info("HTTP Access",
			logger.String("clientIP", clientIP),
			logger.String("method", method),
			logger.String("path", path),
			logger.String("userAgent", c.Request.UserAgent()),
		)

		c.Next() // 处理请求
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
			"code":    common.InnerServerError,
			"message": "Internal Server Error",
		})
	})
}
