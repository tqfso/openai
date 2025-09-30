package proxy

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL   *url.URL               // 目标服务地址
	Alive bool                   // 目标服务是否存活
	mu    sync.RWMutex           // 读写锁，保护Alive字段
	Proxy *httputil.ReverseProxy // 反向代理
}

func (b *Backend) String() string {
	return b.URL.String()
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

func NewBackend(rawURL string) (*Backend, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &Backend{
		URL:   parsedURL,
		Alive: true,
		Proxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}, nil
}
