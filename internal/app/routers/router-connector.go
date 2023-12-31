package routers

import (
	"edumatch/internal/app/docs"
	"edumatch/internal/app/models"
	"edumatch/internal/dependencies"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func ConnectRoutersToHandlers(router *gin.Engine, h dependencies.Handlers) {
	api := router.Group("/api")
	docs.SwaggerInfo.BasePath = ""
	//auth
	api.POST("/auth/signup", h.AuthHandler.SignUp)
	api.POST("/auth/login", h.AuthHandler.Login)
	api.POST("/auth/refresh", h.AuthHandler.RefreshToken)

	//users
	api.GET("/users/", h.AuthHandler.ProtectedEndpoint(models.AdminRole), h.UserHandler.GetUsers)
	api.PATCH("/users/:id", h.AuthHandler.ProtectedEndpoint(), h.UserHandler.UpdateUser)
	api.GET("/users/:id", h.UserHandler.GetUser)
	api.DELETE("users/:id", h.AuthHandler.ProtectedEndpoint(), h.UserHandler.DeleteUser)

	//eduCenters
	api.GET("/educenters/", h.EduCenterHandler.GetAllEduCenters)
	api.GET("/educenters/:id", h.EduCenterHandler.GetEduCenter)
	api.POST("/educenters/", h.AuthHandler.ProtectedEndpoint(), h.EduCenterHandler.CreateEduCenter)
	api.POST("/educenters/rating", h.AuthHandler.ProtectedEndpoint(), h.EduCenterHandler.GiveRating)
	api.PATCH("/educenters/:id", h.AuthHandler.ProtectedEndpoint(), h.EduCenterHandler.UpdateEduCenter)
	api.DELETE("/educenters/:id", h.AuthHandler.ProtectedEndpoint(), h.EduCenterHandler.DeleteEduCenter)
	api.POST("/educenters/location", h.EduCenterHandler.GetEduCenterByLocation)

	//courses
	api.GET("/courses/", h.CourseHandler.GetAllCourses)
	api.GET("/courses/:id", h.CourseHandler.GetCourse)
	api.POST("/courses/", h.AuthHandler.ProtectedEndpoint(), h.CourseHandler.CreateCourse)
	api.POST("/courses/rating", h.AuthHandler.ProtectedEndpoint(), h.CourseHandler.GiveRating)
	api.PATCH("/courses/", h.AuthHandler.ProtectedEndpoint(), h.CourseHandler.UpdateCourse)
	api.DELETE("/courses/:id", h.AuthHandler.ProtectedEndpoint(), h.CourseHandler.DeleteCourse)

	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
