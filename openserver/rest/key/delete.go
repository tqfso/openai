package key

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

func (h *DeleteHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()
	userId := h.GetFromUser()
	if err := service.ApiKey().Delete(ctx, req.ID, userId); err != nil {
		h.SetError(common.GetErrorCode(err, common.Failure), err.Error())
		return
	}
}

func NewDeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &DeleteHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}
