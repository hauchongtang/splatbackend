package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://splatbackend.herokuapp.com"},
		AllowMethods:     []string{"POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://splatbackend.herokuapp.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(middleware.Authentication())

	router.GET("/splat/api", func(c *gin.Context) {
		c.JSON(
			200,
			gin.H{"success": "Access granted"},
		)
	})

	router.Run(":" + port)
}
