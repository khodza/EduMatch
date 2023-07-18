package routers

import (
	"edumatch/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(router *gin.RouterGroup, authHandler handlers.AuthHandlerInterface) {
	authGroup := router.Group("")
	authGroup.POST("/signup", authHandler.SignUp)
	authGroup.POST("/login", authHandler.Login)
	// authGroup.GET("/signup", userHandler.GetUser)
}
