package key

import (
	"common"
	"openserver/rest"
	"openserver/service"
	"time"

	"github.com/gin-gonic/gin"
)

// 查找API密钥，用于API网关调用

type FindHandler struct {
	rest.Handler[FindRequest]
}

type FindRequest struct {
	ID               string `form:"id" binding:"required"`
	WithServiceLimit bool   `form:"withWorkspaceLimit"`
}

type FindResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userID"`
	WorkspaceID string     `json:"workspaceID"`
	Description string     `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`

	ServiceLimits []ServiceLimit `json:"serviceLimits,omitempty"`
}

type ServiceLimit struct {
	ServiceID    string `json:"serviceID"`
	RequestLimit int64  `json:"requestLimit"`
	TokenLimit   int64  `json:"tokenLimit"`
}

func NewFindHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &FindHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *FindHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()

	apiKey, err := service.ApiKey().FindByID(ctx, req.ID)
	if err != nil {
		h.SetError(common.GetErrorCode(err, common.Failure), err.Error())
		return
	}

	response := FindResponse{
		ID:          apiKey.ID,
		UserID:      apiKey.UserID,
		WorkspaceID: apiKey.WorkspaceID,
		ExpiresAt:   apiKey.ExpiresAt,
		Description: apiKey.Description,
	}

	if req.WithServiceLimit {
		usageLimits, err := service.Workspace().ListUsageLimits(ctx, apiKey.WorkspaceID)
		if err != nil {
			h.SetError(common.GetErrorCode(err, common.Failure), err.Error())
			return
		}

		for _, usageLimit := range usageLimits {
			response.ServiceLimits = append(response.ServiceLimits, ServiceLimit{
				ServiceID:    usageLimit.ServiceID,
				RequestLimit: usageLimit.RequestLimit,
				TokenLimit:   usageLimit.TokenLimit,
			})
		}
	}

	h.SetResponseData(response)

}
