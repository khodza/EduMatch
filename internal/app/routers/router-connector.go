package routers

import (
	"edumatch/cmd/docs"
	"edumatch/internal/dependencies"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func ConnectRoutersToHandlers(router *gin.Engine, h dependencies.Handlers) {
	api := router.Group("/api")
	docs.SwaggerInfo.BasePath = "/api"
	//auth
	api.POST("/auth/signup", h.AuthHandler.SignUp)
	api.POST("/auth/login", h.AuthHandler.Login)

	//users
	api.GET("/users/", h.AuthHandler.ProtectedEndpoint, h.UserHandler.GetUsers)
	api.GET("/users/:id", h.UserHandler.GetUser)
	api.PATCH("users/:id", h.AuthHandler.ProtectedEndpoint, h.UserHandler.UpdateUser)
	api.DELETE("users/:id", h.AuthHandler.ProtectedEndpoint, h.UserHandler.DeleteUser)

	//eduCenters
	api.POST("/educenters/", h.AuthHandler.ProtectedEndpoint, h.EduCenterHandler.CreateEduCenter)
	api.GET("/educenters/", h.AuthHandler.ProtectedEndpoint, h.EduCenterHandler.GetEduCenters)
	api.GET("/educenters/:id", h.EduCenterHandler.GetEduCenter)
	api.PATCH("educenters/:id", h.AuthHandler.ProtectedEndpoint, h.EduCenterHandler.UpdateEduCenter)
	api.DELETE("educenters/:id", h.AuthHandler.ProtectedEndpoint, h.EduCenterHandler.DeleteEduCenter)

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
