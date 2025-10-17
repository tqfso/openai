package main

import (
	"flag"
	"fmt"

	"common/logger"
	"openserver/config"
	"openserver/middleware/auth"
	"openserver/rest/api_key"
	"openserver/rest/api_service"
	"openserver/rest/platform_model"
	"openserver/rest/platform_service"
	"openserver/rest/user"
	"openserver/rest/workspace"

	"openserver/middleware"
	"openserver/repository"
	"openserver/rest"

	"github.com/gin-gonic/gin"
)

func main() {

	// 解析参数
	configFileName := flag.String("config", "config/config.yaml", "config file name")
	host := flag.String("host", "", "listen ip")
	port := flag.Int("port", 8080, "listen port")

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
	SetRouter(r)

	// 未找到路由
	r.NoRoute(rest.NewNotFoundHandler())

	serverAddress := fmt.Sprintf("%s:%d", *host, *port)
	r.Run(serverAddress)
}

func SetRouter(r *gin.Engine) {
	SetUserRouter(r)
	SetGatewayRouter(r)
	SetCloudRouter(r)
}

func SetUserRouter(r *gin.Engine) {
	u := r.Group("/v1/user", auth.ZUserAuthHander())
	{
		u.POST("/create", user.NewCreateHandler())
	}

	u = r.Group("/v1/workspace", auth.ZUserAuthHander())
	{
		u.POST("/create", workspace.NewCreateHandler())
		u.POST("/delete", workspace.NewDeleteHandler())

		u.POST("/grant_model", workspace.NewGrantModelHandler())
		u.POST("/cancel_model", workspace.NewCancelModelHandler())
	}

	{
		r.POST("/v1/key/create", auth.ZUserAuthHander(), api_key.NewCreateHandler())
		r.POST("/v1/key/delete", auth.ZUserAuthHander(), api_key.NewDeleteHandler())
		r.GET("/v1/key/list", auth.ZUserAuthHander(), api_key.NewListHandler())
	}
}

func SetGatewayRouter(r *gin.Engine) {

	r.GET("/v1/key/info", auth.ZGatewayAuthHander(), api_key.NewInfoHandler())
}

func SetCloudRouter(r *gin.Engine) {
	u := r.Group("/v1/apiservice", auth.ZCloudAuthHander())
	{
		u.POST("/create", api_service.NewCreateHandler())
		u.POST("/delete", api_service.NewDeleteHandler())
	}

	u = r.Group("/v1/pm", auth.ZCloudAuthHander())
	{
		u.POST("/create", platform_model.NewCreateHandler())
		u.POST("/delete", platform_model.NewDeleteHandler())
		u.GET("/list", platform_model.NewListHandler())
	}

	u = r.Group("/v1/ps", auth.ZCloudAuthHander())
	{
		u.POST("/deploy", platform_service.NewDeployHandler())
		u.POST("/release", platform_service.NewReleaseHandler())
	}
}
