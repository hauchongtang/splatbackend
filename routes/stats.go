package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/controllers"
	"github.com/hauchongtang/splatbackend/middleware"
)

// get routes for user signup and login
func StatsRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/stats/mostpopular", controllers.GetMostPopularModule())
}
