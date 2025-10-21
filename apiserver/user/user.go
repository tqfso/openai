package user

import (
	"apiserver/client/openserver"
	"common"
	"context"
	"sync"
	"time"
)

var (
	mutex     sync.Mutex
	apiKeys   ApiKeys
	usageLogs UsageLogs
)

func init() {
	apiKeys = make(ApiKeys)
	usageLogs = make(UsageLogs)
}

// 查找API密钥
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

// 记录使用量
func AddUsageLog(keyID string, usageLog *UsageLog) {
	usageLog.Timestampt = time.Now().UnixMilli()
	mutex.Lock()
	defer mutex.Unlock()
	usageLogs.Add(keyID, usageLog)
}

// 报告使用量
func ReportUsageLog() error {
	mutex.Lock()
	// report := usageLogs
	usageLogs = make(UsageLogs)
	mutex.Unlock()

	return nil
}
