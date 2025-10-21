package user

type UsageLog struct {
	Timestampt   int64
	ServiceID    string
	Status       UsageStatus
	InputTokens  int64
	OutputTokens int64
	ResponseTime int64
}

type UsageStatus int

const (
	UsageSuccess UsageStatus = 0
	UsageFailed  UsageStatus = 1
)

type UsageLogs map[string]*UsageLogInfo

type UsageLogInfo struct {
	UsageLogs []*UsageLog
}

func (u UsageLogs) Add(keyID string, usageLog *UsageLog) {
	found := u[keyID]
	if found == nil {
		found = &UsageLogInfo{}
	}

	found.UsageLogs = append(found.UsageLogs, usageLog)
}

func (u UsageLogs) Clone() UsageLogs {
	usageLogs = make(UsageLogs)
	for k, v := range u {
		usageLog := *v
		usageLogs[k] = &usageLog
	}
	return usageLogs
}
