package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/middleware"
	"github.com/hauchongtang/splatbackend/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.Use(middleware.Authentication())

	router.GET("/splat/api", func(c *gin.Context) {
		c.JSON(
			200,
			gin.H{"success": "Access granted"},
		)
	})

	router.Run(":" + port)
}
