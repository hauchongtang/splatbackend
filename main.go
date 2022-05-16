package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/middleware"
	"github.com/hauchongtang/splatbackend/routes"
)

func main() {
	port := os.Getenv("PORT")

	var w http.ResponseWriter

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.AuthRoutes(router, &w)
	routes.UserRoutes(router, &w)

	router.Use(middleware.Authentication(&w))

	router.GET("/splat/api", func(c *gin.Context) {
		c.JSON(
			200,
			gin.H{"success": "Access granted"},
		)
	})

	router.Run(":" + port)
}
