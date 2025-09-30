package main

import (
	"apiserver/config"
	"apiserver/middleware"
	"apiserver/proxy"
	"apiserver/rest"
	"common/logger"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	// 解析参数
	cfgfile := flag.String("config", "config/config.yaml", "config from file")
	flag.Parse()

	// 初始化配置
	if err := config.Load(*cfgfile); err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	// 初始化日志
	defer logger.Sync()
	logger.Init(*config.GetLog())
	logger.Info("Application started", logger.Any("config", config.GetConfig()))
	gin.DefaultWriter = logger.GetWriter()

	// HTTP服务

	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery())

	// 分组路由
	SetRoute(r)

	// 未找到路由
	r.NoRoute(rest.NewNotFoundHandler())

	// 启动服务
	serverconfig := config.GetServer()
	serverAddress := fmt.Sprintf("%s:%d", serverconfig.Host, serverconfig.Port)
	r.Run(serverAddress)

}

func SetRoute(r *gin.Engine) {
	SetProxyRoute(r)
}

func SetProxyRoute(r *gin.Engine) {

	// 测试
	proxy.CreateService("test", []string{"http://172.21.21.98:8000"})

	r.Any("/*proxyPath", proxy.ReverseProxyHandler())
}
