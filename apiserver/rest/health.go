package rest

import (
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	Handler[any]
}

func (h *HealthHandler) Handle() {
}

func NewHealthHandler() gin.HandlerFunc {
	h := &HealthHandler{}
	h.SetTaskHandler(h)
	return h.OnRequest
}
