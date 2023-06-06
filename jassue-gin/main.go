package main

import (
	"jassue-gin/bootstrap"
	"jassue-gin/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	bootstrap.InitializeConfig()

	// 初始化日志
	global.App.Log = bootstrap.InitializeLog()
	global.App.Log.Info("log init success!")

	r := gin.Default()

	// 测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 启动服务器
	r.Run(":" + global.App.Config.App.Port)
}
