package main

import (
	"flag"
	"fmt"

	"common/logger"
	"openserver/config"
	"openserver/middleware/auth"
	"openserver/rest/key"
	"openserver/rest/user"
	"openserver/rest/workspace"

	"openserver/middleware"
	"openserver/repository"
	"openserver/rest"

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

	// 初始化数据库
	if err := repository.Init(); err != nil {
		logger.Error("failed to connect to database:", logger.Err(err))
		return
	}

	defer repository.Close()

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
	SetUserRoute(r)
	SetWorkspaceRoute(r)
	SetApiKeyRoute(r)
}

func SetUserRoute(r *gin.Engine) {
	u := r.Group("/v1/user", auth.ZUserAuthHander())
	{
		u.POST("/create", user.NewCreateHandler())
	}
}

func SetWorkspaceRoute(r *gin.Engine) {
	u := r.Group("/v1/workspace", auth.ZUserAuthHander())
	{
		u.POST("/create", workspace.NewCreateHandler())
		u.POST("/delete", workspace.NewDeleteHandler())
	}
}

func SetApiKeyRoute(r *gin.Engine) {
	u := r.Group("/v1/key", auth.ZUserAuthHander())
	{
		u.POST("/create", key.NewCreateHandler())
	}
}
