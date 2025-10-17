package model

import "time"

type UsageLimit struct {
	WorkspaceID  string    `json:"workspaceID" binding:"required"`
	ModelName    string    `json:"modelName" binding:"required"`
	RequestLimit int64     `json:"requestLimit" binding:"required"`
	TokenLimit   int64     `json:"tokenLimit" binding:"required"`
	UpdatedAt    time.Time `json:"-"`
	CreatedAt    time.Time `json:"-"`
}
