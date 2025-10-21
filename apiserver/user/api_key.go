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
	expiredAt := time.Now().Add(time.Second * 180)
	if info.ExpiresAt == nil {
		info.ExpiresAt = &expiredAt
	} else if info.ExpiresAt.After(expiredAt) {
		info.ExpiresAt = &expiredAt
	}
	keys[id] = info
}

func (keys ApiKeys) Del(id string) {
	delete(keys, id)
}
