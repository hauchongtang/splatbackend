package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/controllers"
	"github.com/hauchongtang/splatbackend/middleware"
)

// get routes for user signup and login
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controllers.GetUsers(), middleware.Authentication())
	incomingRoutes.GET("/users/:id", controllers.GetUserById(), middleware.Authentication())
	incomingRoutes.GET("/cached/users/:id", controllers.GetCachedUserById(), middleware.Authentication())
	incomingRoutes.GET("/cached/users", controllers.GetCachedUsers(), middleware.Authentication())
	incomingRoutes.PUT("/users/:id", controllers.IncreasePoints(), middleware.Authentication())
	incomingRoutes.PUT("/users/update/:id", controllers.ModifyParticulars(), middleware.Authentication())
	incomingRoutes.PUT("/users/modules/:id", controllers.UpdateModuleImportLink(), middleware.Authentication())
	incomingRoutes.DELETE("/users/:id", controllers.DeleteUserById(), middleware.Authentication())
}
