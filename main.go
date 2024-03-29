package main

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/hauchongtang/splatbackend/docs"
	"github.com/hauchongtang/splatbackend/middleware"
	"github.com/hauchongtang/splatbackend/routes"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "PUT")
		/*
		   c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		   c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		   c.Writer.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		   c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")
		*/

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// @title     SplatApp Backend API
// @version 1.0
// @description This is the backend service for splatapp at https://github.com/hauchongtang/splatbackend

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @query.collection.format multi
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.TaskRoutes(router)
	routes.StatsRoutes(router)
	routes.DocsRoutes(router)

	router.GET("/splat/api", middleware.Authentication(), func(c *gin.Context) {
		c.JSON(
			200,
			gin.H{"success": "Access granted"},
		)
	})

	router.Run(":" + port)
}
