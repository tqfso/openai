package auth

import (
	"common"
	"net/http"
	"openserver/config"
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

		if auth[7:] != config.GetZdan().ApiServerKey {
			c.JSON(http.StatusUnauthorized, common.Response{Code: common.AuthError, Msg: "Invalid Authorization"})
			c.Abort()
		}

		c.Next()
	}
}
