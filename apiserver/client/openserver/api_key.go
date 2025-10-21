package openserver

import (
	"context"
	"time"
)

// 查询API密钥

type KeyInfoRequest struct {
	ID             string `form:"id"`
	WithUsageLimit bool   `form:"withUsageLimit"`
}

type KeyInfoResponse struct {
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

func FindApiKey(ctx context.Context, id string) (*KeyInfoResponse, error) {
	request := KeyInfoRequest{ID: id, WithUsageLimit: true}
	var resp KeyInfoResponse
	if err := Get(ctx, "/v1/gateway/key/info", request, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
