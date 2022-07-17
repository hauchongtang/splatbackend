package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/controllers"
	"github.com/hauchongtang/splatbackend/middleware"
)

// get routes for user signup and login
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:id", controllers.GetUserById())
	incomingRoutes.PUT("/users/:id", controllers.IncreasePoints())
	incomingRoutes.PUT("/users/update/:id", controllers.ModifyParticulars())
	incomingRoutes.PUT("/users/modules/:id", controllers.UpdateModuleImportLink())
	incomingRoutes.DELETE("/users/:id", controllers.DeleteUserById())
}
