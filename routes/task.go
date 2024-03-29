package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/controllers"
	"github.com/hauchongtang/splatbackend/middleware"
)

// get routes for user signup and login
func TaskRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/tasks", middleware.Authentication(), controllers.GetAllActivity())
	incomingRoutes.GET("/cached/tasks", middleware.Authentication(), controllers.GetCachedAllActivity())
	incomingRoutes.GET("/tasks/:id", middleware.Authentication(), controllers.GetTasksByUserId())
	incomingRoutes.GET("cached/tasks/:id", middleware.Authentication(), controllers.GetCachedTasksByUserId())
	incomingRoutes.PUT("/tasks/:id", middleware.Authentication(), controllers.UpdateHiddenStatus())
	incomingRoutes.POST("/tasks", middleware.Authentication(), controllers.AddTask())
}
