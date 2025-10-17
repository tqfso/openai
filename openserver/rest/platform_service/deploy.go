package platform_service

import (
	"common"
	resource_service "openserver/client/resource/service"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type DeploylHandler struct {
	rest.Handler[DeployRequest]
}

type DeployRequest struct {
	Name      string                    `json:"name" binding:"required"`
	TopoID    uint64                    `json:"topoID" binding:"required"`
	ModelName string                    `json:"modelName" binding:"required"`
	EipInfo   *resource_service.EipInfo `json:"eipInfo,omitempty"`
}

type DeployResponse struct {
	ID string `json:"id"`
}

func NewDeployHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &DeploylHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *DeploylHandler) Handle() {
	req := h.Request
	svc := service.PlatformService()
	svc.EipInfo = req.EipInfo

	id, err := svc.Create(h.GetContext(), req.TopoID, req.Name, req.ModelName)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	h.SetResponseData(DeployResponse{ID: id})
}
