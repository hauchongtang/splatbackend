package routes

import (
	"net/http"

	controller "github.com/hauchongtang/splatbackend/controllers"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func AuthRoutes(incomingRoutes *gin.Engine, w *http.ResponseWriter) {
	incomingRoutes.POST("/users/signup", controller.SignUp(w))
	incomingRoutes.POST("/users/login", controller.Login(w))
}
