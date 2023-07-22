package main

import (
	"fmt"
	"orderingsystem/bootstrap"
	"orderingsystem/global"
)

func main() {
	fmt.Println("main is beginning...")

	bootstrap.InitializeConfig()
	bootstrap.InitializeLog()
	global.App.Log.Info("log init success!")
	global.App.DB = bootstrap.InitializeDatabase()
	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()

	bootstrap.InitializeValidator()
	global.App.Redis = bootstrap.InitializeRedis()
	// r := gin.Default()
	// r.GET("/test", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "hello world")
	// })

	// r.Run(":" + global.App.Config.App.Port)
	bootstrap.RunServer()
}
