package middleware

import (
	"common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ZDanAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("ZCookie")
		if token == "" {
			c.JSON(http.StatusUnauthorized, common.Response{Code: common.AuthError, Msg: "ZCookie required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
