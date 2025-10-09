package model

type Workspace struct {
	ID        uint64 `json:"id"`
	UserID    string `json:"userID"`
	Name      string `json:"name"`
	Desc      string `json:"desc,omitempty"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updateAt"`
	CreatedAt string `json:"createAt"`
}

type UsageLimit struct {
	WorkspaceID  uint64 `json:"workspaceID"`
	ServiceID    string `json:"serviceID"`
	RequestLimit int64  `json:"requestLimit"`
	TokenLimit   int64  `json:"tokenLimit"`
	UpdatedAt    string `json:"updateAt"`
	CreatedAt    string `json:"createAt"`
}
