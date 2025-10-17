package api_key

import (
	"common"
	"openserver/rest"
	"openserver/service"
	"time"

	"github.com/gin-gonic/gin"
)

// 查找API密钥，用于API网关调用

type InfoHandler struct {
	rest.Handler[InfoRequest]
}

type InfoRequest struct {
	ID             string `form:"id" binding:"required"`
	WithUsageLimit bool   `form:"withUsageLimit"`
}

type InfoResponse struct {
	ID          string       `json:"id"`
	UserID      string       `json:"userID"`
	WorkspaceID string       `json:"workspaceID"`
	Description string       `json:"description,omitempty"`
	ExpiresAt   *time.Time   `json:"expiresAt,omitempty"`
	UsageLimits []UsageLimit `json:"usageLimits,omitempty"`
}

type UsageLimit struct {
	ModelName    string `json:"modelName"`
	RequestLimit int64  `json:"requestLimit"`
	TokenLimit   int64  `json:"tokenLimit"`
}

func NewInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &InfoHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *InfoHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()

	apiKey, err := service.ApiKey().FindByID(ctx, req.ID)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	response := InfoResponse{
		ID:          apiKey.ID,
		UserID:      apiKey.UserID,
		WorkspaceID: apiKey.WorkspaceID,
		ExpiresAt:   apiKey.ExpiresAt,
		Description: apiKey.Description,
	}

	if req.WithUsageLimit {
		usageLimits, err := service.Workspace().ListUsageLimits(ctx, apiKey.WorkspaceID)
		if err != nil {
			h.SetErrorWithDefaultCode(err, common.Failure)
			return
		}

		for _, usageLimit := range usageLimits {
			response.UsageLimits = append(response.UsageLimits, UsageLimit{
				ModelName:    usageLimit.ModelName,
				RequestLimit: usageLimit.RequestLimit,
				TokenLimit:   usageLimit.TokenLimit,
			})
		}
	}

	h.SetResponseData(response)

}
