package handlers

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EduCenterHandlerInterface interface {
	GetEduCenters(c *gin.Context)
	CreateEduCenter(c *gin.Context)
	GetEduCenter(c *gin.Context)
	UpdateEduCenter(c *gin.Context)
	DeleteEduCenter(c *gin.Context)
}
type EduCenterHandler struct {
	eduCenterService services.EduCenterServiceInterface
	logger           *zap.Logger
}

func NewEduCenterHandler(eduCenterService services.EduCenterServiceInterface, logger *zap.Logger) EduCenterHandlerInterface {
	return &EduCenterHandler{
		eduCenterService: eduCenterService,
		logger:           logger,
	}
}

func (h *EduCenterHandler) GetEduCenters(c *gin.Context) {
	eduCenters, err := h.eduCenterService.GetEduCenters()

	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "GetEduCenters", h.logger)

	c.JSON(http.StatusOK, eduCenters)
}

func (h *EduCenterHandler) CreateEduCenter(c *gin.Context) {
	var eduCenter models.EduCenter
	if err := HandleJSONBinding(c, &eduCenter, h.logger); err != nil {
		c.Error(err)
		return
	}

	//attaching owner
	userID := c.MustGet("user_id").(uuid.UUID)
	eduCenter.OwnerID = userID

	createdEduCenter, err := h.eduCenterService.CreateEduCenter(eduCenter)

	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "CreateEduCenter", h.logger)

	c.JSON(http.StatusCreated, createdEduCenter)
}

func (h *EduCenterHandler) GetEduCenter(c *gin.Context) {
	eduCenterID, err := GetId(c, h.logger)
	if err != nil {
		return
	}

	eduCenter, err := h.eduCenterService.GetEduCenter(eduCenterID)

	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "GetEduCenter", h.logger)

	c.JSON(http.StatusOK, eduCenter)
}

func (h *EduCenterHandler) UpdateEduCenter(c *gin.Context) {
	var eduCenter models.EduCenter
	if err := HandleJSONBinding(c, &eduCenter, h.logger); err != nil {
		return
	}

	updatedEduCenter, err := h.eduCenterService.UpdateEduCenter(eduCenter)
	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "UpdateEduCenter", h.logger)

	c.JSON(http.StatusOK, updatedEduCenter)
}

func (h *EduCenterHandler) DeleteEduCenter(c *gin.Context) {
	eduCenterID, err := GetId(c, h.logger)
	if err != nil {
		return
	}
	if err := h.eduCenterService.DeleteEduCenter(eduCenterID); err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "DeleteEduCenter", h.logger)

	c.JSON(http.StatusOK, gin.H{"message": "Education Center deleted successfully"})
}
