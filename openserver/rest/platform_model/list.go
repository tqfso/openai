package platform_model

import (
	"common"
	"openserver/model"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	rest.Handler[model.PlatformModelSearchParam]
}

type ListResponse struct {
	TotalCount int                    `json:"totalCount"`
	PageIndex  int                    `json:"pageIndex"`
	PageSize   int                    `json:"pageSize"`
	Models     []*model.PlatformModel `json:"models,omitempty"`
}

func NewListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &ListHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *ListHandler) Handle() {
	req := &h.Request
	ctx := h.GetContext()

	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	models, total, err := service.PlatformModel().List(ctx, req)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	response := ListResponse{
		TotalCount: total,
		PageIndex:  req.PageIndex,
		PageSize:   req.PageSize,
		Models:     models,
	}

	h.SetResponseData(response)

}
