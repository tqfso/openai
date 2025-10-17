package gateway

import (
	"common"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type ModelServicesHandler struct {
	rest.Handler[ModelServicesRequest]
}

type ModelServicesRequest struct {
	ID string `form:"id" binding:"required"`
}

func NewModelServicesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &ModelServicesHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *ModelServicesHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()
	services, err := service.PlatformService().ListByGateway(ctx, req.ID)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	h.SetResponseData(services)
}
