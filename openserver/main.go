package main

import (
	"flag"

	"openserver/config"
	"openserver/logger"
	"openserver/middleware/auth"
	"openserver/rest/test"

	"openserver/middleware"
	"openserver/rest"

	"github.com/gin-gonic/gin"
)

func main() {

	// 解析参数
	cfgfile := flag.String("config", "config/config.yaml", "config from file")
	logfile := flag.String("logs", "logs/app.log", "log to file")
	flag.Parse()

	// 初始化日志
	defer logger.Sync()
	logger.Init("debug", *logfile)
	logger.Info("Application started")
	gin.DefaultWriter = logger.GetWriter()

	// 初始化配置
	if err := config.Load(*cfgfile); err != nil {
		logger.Error(err.Error())
		return
	}

	// HTTP服务

	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery())

	// 分组路由
	SetRoute(r)

	// 未找到路由
	r.NoRoute(rest.NewNotFoundHandler())

	r.Run(config.GetServer().ListenAddress())
}

func SetRoute(r *gin.Engine) {
	SetTestRoute(r)
}

func SetTestRoute(r *gin.Engine) {
	t := r.Group("/v1/test")
	t.GET("/1", auth.ZCloudAuthHander(), test.NewTest1Handler())
}
