package platform_service

import (
	"common"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type ReleaseHandler struct {
	rest.Handler[ReleaseRequest]
}

type ReleaseRequest struct {
	ID string `form:"id" binding:"required"`
}

func NewReleaseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &ReleaseHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *ReleaseHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()
	if err := service.PlatformService().Release(ctx, req.ID); err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}
}
