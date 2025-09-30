package proxy

import (
	"common/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type LoadBalancer struct {
	backends []*Backend // 目标服务列表
	index    int        // 当前选择的目标服务索引
	mu       sync.Mutex // 互斥锁
}

func NewLoadBalancer(targets []string) (*LoadBalancer, error) {
	backends := make([]*Backend, 0)
	for _, addr := range targets {
		u, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		b := &Backend{
			URL:   u,
			Alive: true,
			Proxy: httputil.NewSingleHostReverseProxy(u),
		}
		backends = append(backends, b)
	}
	return &LoadBalancer{backends: backends}, nil
}

func (lb *LoadBalancer) GetBackends() []*Backend {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.backends
}

func (lb *LoadBalancer) getNextAvailableBackend() *Backend {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	total := len(lb.backends)
	for range total {
		lb.index = (lb.index + 1) % total
		if lb.backends[lb.index].IsAlive() {
			return lb.backends[lb.index]
		}
	}
	return nil
}

// 周期性健康检查
func (lb *LoadBalancer) StartHealthCheck(interval time.Duration) {
	for {
		backends := lb.GetBackends()
		for _, backend := range backends {
			resp, err := http.Get(backend.URL.String() + "/health")

			if err != nil {
				logger.Error("backend is down", logger.String("url", backend.String()), logger.Err(err))
				backend.SetAlive(false)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				logger.Error("backend is down", logger.String("url", backend.String()))
				backend.SetAlive(false)
			} else {
				backend.SetAlive(true)
			}
		}
		time.Sleep(interval)
	}
}
