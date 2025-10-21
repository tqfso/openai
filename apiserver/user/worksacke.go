package user

import "apiserver/client/openserver"

type Workspaces map[string]*WorkspaceInfo

type WorkspaceInfo struct {
	UsageLimits []openserver.UsageLimit
}

func (w Workspaces) Set(id string, info *WorkspaceInfo) {
	w[id] = info
}

func (w Workspaces) Get(id string) *WorkspaceInfo {
	return w[id]
}
