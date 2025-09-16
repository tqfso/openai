package test

import (
	"common"
	"openserver/client/resource"
	"openserver/client/resource/service"
	"openserver/rest"

	"github.com/gin-gonic/gin"
)

type Test1Handler struct {
	rest.Handler[service.StatusRequest]
}

func (h *Test1Handler) Handle() {

	resp := service.StatusResponse{}
	if err := resource.Get("/service/status", h.Request, &resp); err != nil {
		h.SetError(common.InnerAccessError, err.Error())
		return
	}

	h.SetResponseData(resp)
}

func NewTest1Handler() gin.HandlerFunc {
	h := &Test1Handler{}
	h.SetTaskHandler(h)
	return h.OnRequest
}
