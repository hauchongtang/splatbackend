package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/controllers"
	"github.com/hauchongtang/splatbackend/middleware"
)

// get routes for user signup and login
func TaskRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/tasks", controllers.GetAllActivity(), middleware.Authentication())
	incomingRoutes.GET("/tasks/:id", controllers.GetTasksById(), middleware.Authentication())
	incomingRoutes.PUT("/tasks/:id", controllers.UpdateHiddenStatus(), middleware.Authentication())
	incomingRoutes.POST("/tasks", controllers.AddTask(), middleware.Authentication())
}
