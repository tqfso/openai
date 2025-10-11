package rest

import (
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	Handler[any]
}

func NewHealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &HealthHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *HealthHandler) Handle() {
}
