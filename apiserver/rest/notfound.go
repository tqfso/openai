package rest

import (
	"common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotFoundHandler struct {
	Handler[any]
}

func NewNotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &NotFoundHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *NotFoundHandler) Handle() {
	h.SetStatusCode(http.StatusNotFound)
	h.SetError(common.HandlerNotFound, "not found handler")
}
