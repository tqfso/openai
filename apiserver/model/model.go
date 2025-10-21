package model

import (
	"apiserver/client/openserver"
	"common/logger"
	"context"
	"sync"
	"time"
)

type Service struct {
	ModelName string
	Power     uint64
	Load      uint64
	Targets   []*openserver.ModelServiceTarget
}

type Model struct {
	mutex    sync.Mutex
	Services map[string]*Service
}

var (
	model Model
)

func init() {
	model = Model{Services: make(map[string]*Service)}
}

func (m *Model) Refresh(services map[string]*Service) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Services = services
}

func (m *Model) FindService(id string) *Service {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Services[id]
}

// 定时加载自己负责的模型服务

func BackgroudLoad(ctx context.Context) {

	logger.Info("Model background task start")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	Load(ctx)

	for {
		select {
		case <-ticker.C:
			Load(ctx)
		case <-ctx.Done():
			goto end
		}
	}

end:
	logger.Info("Model background task final")
}

func Load(ctx context.Context) {

	logger.Debug("Load Services")

	resp, err := openserver.FindModelServices(ctx)
	if err != nil {
		logger.Error("FindModelServices", logger.Err(err))
		return
	}

	services := make(map[string]*Service)
	for _, s := range resp {
		service := &Service{
			ModelName: s.ModelName,
			Power:     s.Power,
			Load:      s.Load,
			Targets:   s.Targets,
		}

		services[s.ID] = service
	}

	model.Refresh(services)
}

func FindService(id string) *Service {
	return model.FindService(id)
}
