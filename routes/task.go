package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/controllers"
	"github.com/hauchongtang/splatbackend/middleware"
)

// get routes for user signup and login
func TaskRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/tasks", controllers.GetAllActivity())
	incomingRoutes.GET("/tasks/:id", controllers.GetTasksById())
	incomingRoutes.POST("/tasks", controllers.AddTask())
}
