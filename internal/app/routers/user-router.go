package routers

import (
	"edumatch/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(router *gin.RouterGroup, userHandler handlers.UserHandlerInterface) {
	userGroup := router.Group("")
	userGroup.GET("/", userHandler.GetUsers)
	userGroup.GET("/:id", userHandler.GetUser)
	userGroup.PATCH("/:id", userHandler.UpdateUser)
	userGroup.DELETE("/:id", userHandler.DeleteUser)
	// userGroup.POST("/", userHandler.CreateUser)
}
