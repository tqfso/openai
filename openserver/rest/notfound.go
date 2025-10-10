package rest

import (
	"common"

	"github.com/gin-gonic/gin"
)

type NotFoundHandler struct {
	Handler[any]
}

func (h *NotFoundHandler) Handle() {
	h.SetError(common.HandlerNotFound, "not found handler")
}

func NewNotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &NotFoundHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}
