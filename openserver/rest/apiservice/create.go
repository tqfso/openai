package apiservice

import (
	"common"
	reservice "openserver/client/resource/service"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type CreateHandler struct {
	rest.Handler[CreateRequest]
}

type CreateRequest struct {
	Name    string             `json:"name" binding:"required"`
	TopoID  uint64             `json:"topoID" binding:"required"`
	EipInfo *reservice.EipInfo `json:"eipInfo" binding:"required"`
}

type CreateResponse struct {
	ID string `json:"id"`
}

func NewCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &CreateHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *CreateHandler) Handle() {
	req := &h.Request
	apiSerivce := service.ApiService()

	found, err := apiSerivce.FindByTopoID(h.GetContext(), req.TopoID)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	if found != nil {
		h.SetResponseData(CreateResponse{ID: found.ID})
		return
	}

	id, err := apiSerivce.Create(h.GetContext(), uint64(req.TopoID), req.Name, req.EipInfo)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	h.SetResponseData(CreateResponse{ID: id})
}
