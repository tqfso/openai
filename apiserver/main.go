package main

import (
	"apiserver/config"
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

}
