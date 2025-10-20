package main

import (
	"apiserver/config"
	"apiserver/middleware"
	"apiserver/rest"
	"common/logger"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	// 解析参数
	configFileName := flag.String("config", "config/config.yaml", "config from file")
	host := flag.String("host", "", "listen ip")
	port := flag.Int("port", 8000, "listen port")
	flag.Parse()

	// 初始化配置
	if err := config.Load(*configFileName); err != nil {
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
	SetRouter(r)

	// 未找到路由
	r.NoRoute(rest.NewNotFoundHandler())

	// 启动服务
	serverAddress := fmt.Sprintf("%s:%d", *host, *port)
	r.Run(serverAddress)

}

func SetRouter(r *gin.Engine) {

	r.GET("/health", rest.NewHealthHandler())

	SetProxyRouter(r)
}

func SetProxyRouter(r *gin.Engine) {
}
