package routers

import (
	"edumatch/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCourseRouter(router *gin.RouterGroup, userHandler handlers.CourseHandlerInterface) {
	courseGroup := router.Group("")
	courseGroup.POST("/", userHandler.CreateCourse)
}
