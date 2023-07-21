package handlers

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CourseHandlerInterface interface {
	CreateCourse(c *gin.Context)
	// CreateUser(c *gin.Context)
}
type CourseHandler struct {
	courseService services.CourseServiceInterface
	logger        *zap.Logger
}

func NewCourseHandler(courseService services.CourseServiceInterface, logger *zap.Logger) CourseHandlerInterface {
	return &CourseHandler{
		courseService: courseService,
		logger:        logger,
	}
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course models.Course
	if err := HandleJSONBinding(c, &course, h.logger); err != nil {
		return
	}

	createdUser, err := h.courseService.CreateCourse(course)

	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "CreateCourse", h.logger)

	c.JSON(http.StatusCreated, createdUser)
}
