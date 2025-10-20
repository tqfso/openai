package user

import "time"

type ApiKeys map[string]*ApiKeyInfo

type ApiKeyInfo struct {
	WorkspaceInfo *WorkspaceInfo // 可能为空
	ExpiresAt     *time.Time     // 到期时间
}

func (keys ApiKeys) Find(id string) *ApiKeyInfo {
	info := keys[id]
	if info == nil {
		return nil
	}

	if info.ExpiresAt.Before(time.Now()) {
		delete(keys, id)
		return nil
	}

	return info
}

func (keys ApiKeys) Add(id string, info *ApiKeyInfo) {
	if info.ExpiresAt == nil {
		expiredAt := time.Now().Add(time.Second * 60)
		info.ExpiresAt = &expiredAt
	}
	keys[id] = info
}

func (keys ApiKeys) Del(id string) {
	delete(keys, id)
}
