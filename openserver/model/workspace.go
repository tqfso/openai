package model

import "time"

type Workspace struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userID"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updateAt"`
	CreatedAt time.Time `json:"createAt"`
}

type UsageLimit struct {
	WorkspaceID  string    `json:"workspaceID"`
	ServiceID    string    `json:"serviceID"`
	RequestLimit int64     `json:"requestLimit"`
	TokenLimit   int64     `json:"tokenLimit"`
	UpdatedAt    time.Time `json:"updateAt"`
	CreatedAt    time.Time `json:"createAt"`
}

const (
	MaxWorkspaceCount int = 10
)
