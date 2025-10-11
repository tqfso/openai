package model

import "time"

type ApiKey struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userID"`
	WorkspaceID string     `json:"workspaceID"`
	Description string     `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	UpdatedAt   time.Time  `json:"updateAt"`
	CreatedAt   time.Time  `json:"createAt"`
}

type ApiKeyEx struct {
	ApiKey        `json:",inline"`
	WorkspaceName string `json:"workspaceName,omitempty"`
}
