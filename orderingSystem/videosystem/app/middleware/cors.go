package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Core() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS"}
	config.ExposeHeaders = []string{"New-Token", "New-Expires-In", "Content-Disposition"}
	return cors.New(config)
}
