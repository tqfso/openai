package model

import (
	"apiserver/client/openserver"
	"common/logger"
	"context"
	"sync"
	"time"
)

// 模型名称对应模型服务列表
type Models map[string]*Services

type Manager struct {
	mutex sync.Mutex
	modes Models
}

var (
	manager Manager
)

func init() {
	manager = Manager{modes: make(Models)}
}

func (m *Manager) Refresh(models Models) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.modes = models
}

func (m *Manager) SelectTarget(modelName string) *openserver.ServiceTarget {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	found := m.modes[modelName]
	if found == nil {
		return nil
	}

	return found.SelectTarget()
}

// 选择转发模型
func SelectTarget(modelName string) *openserver.ServiceTarget {
	return manager.SelectTarget(modelName)
}

// 加载模型服务任务
func LoadServicesTask(ctx context.Context) {

	logger.Info("Model background task start")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	LoadServices(ctx)

	for {
		select {
		case <-ticker.C:
			LoadServices(ctx)
		case <-ctx.Done():
			goto end
		}
	}

end:
	logger.Info("Model background task final")
}

func LoadServices(ctx context.Context) {

	logger.Debug("Load Services")

	resp, err := openserver.FindModelServices(ctx)
	if err != nil {
		logger.Error("FindModelServices", logger.Err(err))
		return
	}

	models := make(Models)
	for _, s := range resp {
		service := &Service{
			ID:      s.ID,
			Power:   s.Power,
			Load:    s.Load,
			Targets: s.Targets,
		}

		found := models[s.ModelName]
		if found == nil {
			found = &Services{}
			models[s.ModelName] = found
		}

		found.Services = append(found.Services, service)

	}

	manager.Refresh(models)
}
