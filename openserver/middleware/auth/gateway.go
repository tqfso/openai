package auth

import (
	"common"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ZGatewayAuthHander() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, common.Response{Code: common.AuthError, Msg: "Authorization required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(auth, "Bearer") {
			c.JSON(http.StatusUnauthorized, common.Response{Code: common.AuthError, Msg: "Bearer required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
