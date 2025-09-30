package proxy

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// 反向代理实现

var (
	proxy *Proxy
)

type Proxy struct {
	balancers map[string]*LoadBalancer // 每个服务对应一个负载均衡器
	mu        sync.Mutex               // 互斥锁
}

func init() {
	proxy = &Proxy{
		balancers: make(map[string]*LoadBalancer),
	}
}

// 创建服务的负载均衡器
func CreateService(name string, targets []string) error {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()

	lb, err := NewLoadBalancer(targets)
	if err != nil {
		return err
	}
	proxy.balancers[name] = lb
	return nil
}

// 查找服务的负载均衡器
func GetLoadBalancer(name string) *LoadBalancer {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()
	return proxy.balancers[name]
}

func ReverseProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		balancer := GetLoadBalancer("test")
		if balancer == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no such service"})
			return
		}
		backend := balancer.getNextAvailableBackend()
		if backend == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no healthy backend"})
			return
		}
		backend.Proxy.ServeHTTP(c.Writer, c.Request)
	}
}
