package routers

import (
	"edumatch/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCourseRouter(router *gin.RouterGroup, courseHandler handlers.CourseHandlerInterface) {
	courseGroup := router.Group("")
	courseGroup.POST("/", courseHandler.CreateCourse)
	courseGroup.GET("/:id", courseHandler.GetCourse)
	courseGroup.GET("/", courseHandler.GetAllCourses)
	courseGroup.PUT("/", courseHandler.UpdateCourse)
	courseGroup.DELETE("/:id", courseHandler.DeleteCourse)
}
