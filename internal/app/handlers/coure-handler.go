package handlers

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CourseHandlerInterface interface {
	CreateCourse(c *gin.Context)
	UpdateCourse(c *gin.Context)
	GetCourse(c *gin.Context)
	GetAllCourses(c *gin.Context)
	DeleteCourse(c *gin.Context)
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

	c.JSON(http.StatusOK, createdUser)
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	var newCourse models.Course
	if err := HandleJSONBinding(c, &newCourse, h.logger); err != nil {
		c.Error(err)
		return
	}

	course, err := h.courseService.UpdateCourse(newCourse)
	if err != nil {
		c.Error(err)
		return
	}

	LoggingResponse(c, "UpdateCourse", h.logger)

	c.JSON(http.StatusAccepted, course)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	course, err := h.courseService.GetCourse(id)
	if err != nil {
		c.Error(err)
		return
	}

	LoggingResponse(c, "GetCourse", h.logger)

	c.JSON(http.StatusAccepted, course)
}

func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.courseService.GetAllCourses()
	if err != nil {
		c.Error(err)
		return
	}

	LoggingResponse(c, "GetAllCourses", h.logger)

	c.JSON(http.StatusAccepted, courses)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseID, err := GetId(c, h.logger)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.courseService.DeleteCourse(courseID.String()); err != nil {
		c.Error(err)
		return
	}

	LoggingResponse(c, "DeleteCourse", h.logger)

	c.JSON(http.StatusAccepted, "Course deleted")

}
