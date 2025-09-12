package middleware

import (
	"common"
	"net/http"
	"openserver/logger"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func HttpRecovery() gin.HandlerFunc {
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
