package test

import (
	"common"
	"openserver/client/resource"
	"openserver/rest"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	rest.Handler[resource.StatusRequest]
}

func (h *TestHandler) Handle() {

	resp := resource.StatusResponse{}
	if err := resource.Get("/service/status", h.Request, &resp); err != nil {
		h.SetError(common.InnerAccessError, err.Error())
		return
	}

	h.SetResponseData(resp)
}

func NewTestHandler() gin.HandlerFunc {
	h := &TestHandler{}
	h.SetTaskHandler(h)
	return h.OnRequest
}
