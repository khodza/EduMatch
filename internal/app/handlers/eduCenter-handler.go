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
	GetAllEduCenters(c *gin.Context)
	CreateEduCenter(c *gin.Context)
	GetEduCenter(c *gin.Context)
	UpdateEduCenter(c *gin.Context)
	DeleteEduCenter(c *gin.Context)
	GiveRating(c *gin.Context)
	GetEduCenterByLocation(c *gin.Context)
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

// Get EduCenters ...
// @Summary Get EduCenters
// @Description This API for getting all EduCenters
// @Tags EduCenter
// @Accept json
// @Produce json
// @Success 200 {object} models.AllEduCenters
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/educenters [GET]
func (h *EduCenterHandler) GetAllEduCenters(c *gin.Context) {
	eduCenters, err := h.eduCenterService.GetAllEduCenters()

	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "GetAllEduCenters", h.logger)

	c.JSON(http.StatusOK, eduCenters)
}

// CreateEduCenter ...
// @Summary Create Edu Center
// @Description This API for creating Edu Center
// @Security BearerAuth
// @Tags EduCenter
// @Accept json
// @Produce json
// @Param body body models.EduCenter true "CourseBody"
// @Success 200 {object} models.EduCenter
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/educenters [POST]
func (h *EduCenterHandler) CreateEduCenter(c *gin.Context) {
	var eduCenter models.CreateEduCenterDto
	if err := HandleFormDataBinding(c, &eduCenter, h.logger); err != nil {
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

// Get EduCenter ...
// @Summary Get Edu Center
// @Description This API for getting EduCenter
// @Tags EduCenter
// @Accept json
// @Produce json
// @Param id path string true "EduCenter_ID"
// @Success 200 {object} models.EduCenter
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/educenters [GET]
func (h *EduCenterHandler) GetEduCenter(c *gin.Context) {
	eduCenterID, err := GetId(c, h.logger)
	if err != nil {
		c.Error(err)
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

//		Update EduCenter ...
//	 @Summary Update EduCenter
//
// @Description This API for updating eduCenter
//
//	@Security BearerAuth
//	@Tags EduCenter
//	@Accept json
//	@Produse json
//	@Param body body models.EduCenter true "EduCenter"
//	@Success 200 {object} models.EduCenter
//	@Failure 400 {object} models.CustomError
//	@Failure 500 {object} models.CustomError
//	@Router /api/educenters [PATCH]
func (h *EduCenterHandler) UpdateEduCenter(c *gin.Context) {
	var eduCenter models.UpdateEduCenterDto
	if err := HandleFormDataBinding(c, &eduCenter, h.logger); err != nil {
		c.Error(err)
		return
	}
	//get edu Center ID and attaching it
	eduCenterID, err := GetId(c, h.logger)
	if err != nil {
		c.Error(err)
		return
	}
	eduCenter.ID = eduCenterID

	updatedEduCenter, err := h.eduCenterService.UpdateEduCenter(eduCenter)
	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "UpdateEduCenter", h.logger)

	c.JSON(http.StatusOK, updatedEduCenter)
}

// Delete EduCenter ....
// @Summary Delete EduCenter
// @Description This API for deleting EduCenter
// @Security BearerAuth
// @Tags EduCenter
// @Accept json
// @Produce json
// @Param id path string true "EduCenter_ID"
// @Success 200 {object} models.Empty
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/educenters [DELETE]
func (h *EduCenterHandler) DeleteEduCenter(c *gin.Context) {
	eduCenterID, err := GetId(c, h.logger)
	if err != nil {
		c.Error(err)
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

func (h *EduCenterHandler) GiveRating(c *gin.Context) {
	var rating models.EduCenterRating
	if err := HandleJSONBinding(c, &rating, h.logger); err != nil {
		c.Error(err)
		return
	}

	//attaching owner
	userID := c.MustGet("user_id").(uuid.UUID)
	rating.OwnerID = userID

	if err := h.eduCenterService.GiveRating(rating); err != nil {
		c.Error(err)
		return
	}
	//logging
	LoggingResponse(c, "GiveRating", h.logger)
	c.JSON(http.StatusOK, gin.H{"message": "Rating accepted"})
}

// GetEduCenterByLocation(location models.NearEduCenterDto) (models.AllNearEduCenters, error)
// Get EduCenter By Location
// @Summary Get EduCenter By Location
// @Description This API for getting educenters by location
// @Tags EduCenter
// @Accept json
// @Produce json
// @Param body body models.NearEduCenterDto true "EduCenter_Location"
// @Success 200 {object} models.AllNearEduCenters
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/educenters/location [POST]
func (h *EduCenterHandler) GetEduCenterByLocation(c *gin.Context) {
	var location models.NearEduCenterDto

	err := HandleJSONBinding(c, &location, h.logger)
	if err != nil {
		c.Error(err)
		return
	}
	eduCenters, err := h.eduCenterService.GetEduCenterByLocation(location)
	if err != nil {
		c.Error(err)
		return
	}

	LoggingResponse(c, "GetEduCenterByLocation", h.logger)

	c.JSON(http.StatusAccepted, eduCenters)
}
