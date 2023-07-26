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
	UpdateCourse(c *gin.Context)
	GetCourse(c *gin.Context)
	GetAllCourses(c *gin.Context)
	DeleteCourse(c *gin.Context)
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

// CreateCourse ...
// @Summary CreateCourse
// @Description This API for creating course
// @Security BearerAuth
// @Tags Course
// @Accept json
// @Produce json
// @Param body body models.Course true "CourseBody"
// @Success 200 {object} models.Course
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/courses [POST]
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

// UpdateCourse ...
// @Summary UpdateCourse
// @Description This API for updating Course
// @Security BearerAuth
// @Tags Course
// @Accept json
// @Produce json
// @Param body body models.Course true "CourseBody"
// @Success 200 {object} models.Course
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/courses [PUT]
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

// GetCourse ...
// @Summary GetCourse
// @Description This API for getting Course
// @Tags Course
// @Accept json
// @Produce json
// @Param id path string true "Course_id"
// @Success 200 {object} models.Course
// @Failure 400 {object} models.CustomError
// @Failuer 500 {object} models.CustomError
// @Router /api/courses/{id} [GET]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID, err := GetId(c, h.logger)
	if err != nil {
		return
	}
	course, err := h.courseService.GetCourse(courseID)
	if err != nil {
		c.Error(err)
		return
	}

	LoggingResponse(c, "GetCourse", h.logger)

	c.JSON(http.StatusAccepted, course)
}

// GetAllCourses ...
// @Summary GetAllCourses
// @Description This API for getting all courses
// @Tags Course
// @Accept json
// @Produce json
// @Success 200 {object} models.AllCourses
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/courses [GET]
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.courseService.GetAllCourses()
	if err != nil {
		c.Error(err)
		return
	}

	LoggingResponse(c, "GetAllCourses", h.logger)

	c.JSON(http.StatusAccepted, courses)
}

// DeleteCourse ...
// @Summary DeleteCourse
// @Description This API for deleting Course
// @Security BearerAuth
// @Tags Course
// @Accept json
// @Produce json
// @Param id path string true "Course_id"
// @Success 200 {object} models.Empty
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/courses/{id} [DELETE]
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
