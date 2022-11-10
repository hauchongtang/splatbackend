package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// get routes for swagger
func DocsRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
