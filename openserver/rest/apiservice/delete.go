package apiservice

import (
	"common"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type DeleteHandler struct {
	rest.Handler[DeleteRequest]
}

type DeleteRequest struct {
	ID string `form:"id" binding:"required"`
}

func NewDeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &DeleteHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *DeleteHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()
	if err := service.ApiService().Delete(ctx, req.ID); err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}
}
