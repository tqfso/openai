package user

import (
	"apiserver/client/openserver"
	"common"
	"context"
	"sync"
)

var mutex sync.Mutex

var (
	apiKeys ApiKeys
)

func init() {
	apiKeys = make(ApiKeys)
}

func FindKey(ctx context.Context, id string) (*ApiKeyInfo, error) {
	mutex.Lock()
	found := apiKeys.Find(id)
	if found == nil {
		mutex.Unlock()
		resp, err := openserver.FindApiKey(ctx, id)
		if err != nil {
			if !common.IsErrorCode(err, common.ApiKeyNotFound) {
				return nil, err
			}
		}

		found = &ApiKeyInfo{}
		if resp != nil {
			found.WorkspaceInfo = &WorkspaceInfo{UsageLimits: resp.UsageLimits}
			found.ExpiresAt = resp.ExpiresAt
		}

		mutex.Lock()
		apiKeys.Add(id, found)
	}

	defer mutex.Unlock()

	if found.WorkspaceInfo == nil {
		return nil, nil
	}

	return found, nil
}
