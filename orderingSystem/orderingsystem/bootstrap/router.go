package bootstrap

import (
	"context"
	"log"
	"net/http"
	"orderingsystem/global"
	"orderingsystem/routers"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")
	routers.SetApiGroupRouters(apiGroup)
	return router
}

func RunServer() {
	r := setupRouter()
	// r.Run(":" + global.App.Config.App.Port)
	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.App.Log.Info("start err!")
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<- quit
	log.Println("Shutdonw Server ......")

	ctx ,cancel := context.WithTimeout(context.Background(),5 * time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx);err != nil {
		log.Fatal("Server Shutdown:",err)
	}
	log.Println("Server Exiting")
}