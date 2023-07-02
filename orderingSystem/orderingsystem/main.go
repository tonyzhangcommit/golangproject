package main

import (
	"fmt"
	"net/http"
	"orderingsystem/bootstrap"
	"orderingsystem/global"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	r.Run(":" + global.App.Config.App.Port)
}
