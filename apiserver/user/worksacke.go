package user

import "apiserver/client/openserver"

type Workspaces map[string]*WorkspaceInfo

type WorkspaceInfo struct {
	UsageLimits []openserver.UsageLimit
}
