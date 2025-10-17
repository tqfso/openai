package api_key

import (
	"common"
	"openserver/model"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	rest.Handler[ListRequest]
}

type ListRequest struct {
	PageIndex int `form:"pageIndex"`
	PageSize  int `form:"pageSize"`
}

type ListResponse struct {
	TotalCount int               `json:"totalCount"`
	PageIndex  int               `json:"pageIndex"`
	PageSize   int               `json:"pageSize"`
	Keys       []*model.ApiKeyEx `json:"keys,omitempty"`
}

func NewListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &ListHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *ListHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()

	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	apiKeys, total, err := service.ApiKey().ListByUser(ctx, h.GetFromUser(), req.PageIndex, req.PageSize)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	response := ListResponse{
		TotalCount: total,
		PageIndex:  req.PageIndex,
		PageSize:   req.PageSize,
		Keys:       apiKeys,
	}

	h.SetResponseData(response)

}
