package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/controllers"
	"github.com/hauchongtang/splatbackend/middleware"
)

// get routes for user authentication
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", middleware.Authentication(), controllers.GetUsers())
	incomingRoutes.GET("/users/:id", middleware.Authentication(), controllers.GetUserById())
	incomingRoutes.GET("/cached/users/:id", middleware.Authentication(), controllers.GetCachedUserById())
	incomingRoutes.GET("/cached/users", middleware.Authentication(), controllers.GetCachedUsers())
	incomingRoutes.PUT("/users/:id", middleware.Authentication(), controllers.IncreasePoints())
	incomingRoutes.PUT("/users/update/:id", middleware.Authentication(), controllers.ModifyParticulars())
	incomingRoutes.PUT("/users/modules/:id", middleware.Authentication(), controllers.UpdateModuleImportLink())
	incomingRoutes.DELETE("/users/:id", middleware.Authentication(), controllers.DeleteUserById())
}
