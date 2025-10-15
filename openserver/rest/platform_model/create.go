package platform_model

import (
	"common"
	"openserver/model"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type CreateHandler struct {
	rest.Handler[CreateRequest]
}

type CreateRequest struct {
	Name             string            `json:"name" binding:"required"`
	Provider         uint64            `json:"provider" binding:"required"`
	Classes          []uint64          `json:"classes,omitempty" binding:"required"`
	Abilities        []uint64          `json:"abilities,omitempty"`
	MaxContextLength uint64            `json:"maxContextLength"`
	DeployInfo       *model.DeployInfo `json:"deployInfo,omitempty"`
	Description      string            `json:"description,omitempty"`
}

func NewCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &CreateHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *CreateHandler) Handle() {
	req := h.Request
	pm := &model.PlatformModel{
		Name:             req.Name,
		Provider:         req.Provider,
		Classes:          req.Classes,
		Abilities:        req.Abilities,
		MaxContextLength: req.MaxContextLength,
		DeployInfo:       req.DeployInfo,
		Description:      req.Description,
	}

	if err := service.PlatformModel().Create(h.GetContext(), pm); err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}
}
