package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hauchongtang/splatbackend/controllers"
	"github.com/hauchongtang/splatbackend/middleware"
)

// get routes for user signup and login
func UserRoutes(incomingRoutes *gin.Engine, w *http.ResponseWriter) {
	incomingRoutes.Use(middleware.Authentication(w))
	incomingRoutes.GET("/users", controllers.GetUsers(w))
	incomingRoutes.GET("/users/:id", controllers.GetUserById(w))
}
