package handlers

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandlerInterface interface {
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	// CreateUser(c *gin.Context)
}
type UserHandler struct {
	userService services.UserServiceInterface
	logger      *zap.Logger
}

func NewUserHandler(userService services.UserServiceInterface, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// Create User ...
// @Summary Create User
// @Description This API for creating new user
// @Security BearerAuth
// @Tags user
// @Accept json
// @Produce json
// @Param body body models.User true "User_body"
// @Success 200 {object} models.User
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/users/ [POST]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()

	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "GetUsers", h.logger)

	c.JSON(http.StatusOK, users)
}

// Get User ...
// @Summary Get User
// @Description This API for getting userby ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User_ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/users/{id} [GET]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := GetId(c, h.logger)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := h.userService.GetUser(userID)

	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "GetUser", h.logger)

	c.JSON(http.StatusOK, user)
}

// Update User ...
// @Summary Update User
// @Description This API for updating user
// @Security BearerAuth
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User_ID"
// @Param body body models.User true "User_body"
// @Success 200 {object} models.User
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/users/{id} [PATCH]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.UpdateUserDto
	if err := HandleFormDataBinding(c, &user, h.logger); err != nil {
		c.Error(err)
		return
	}

	userID, err := GetId(c, h.logger)
	if err != nil {
		c.Error(err)
		return
	}
	user.ID = userID

	updatedUser, err := h.userService.UpdateUser(user)

	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "UpdateUser", h.logger)

	c.JSON(http.StatusOK, updatedUser)
}

// Delete User ..
// @Summary Delete User
// @Description This API for deleting user
// @Security BearerAuth
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User_ID"
// @Success 200 {object} models.Empty
// @Failure 400 {object} models.CustomError
// @Failure 500 {object} models.CustomError
// @Router /api/users/{id} [DELETE]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := GetId(c, h.logger)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "DeleteUser", h.logger)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

//NOT USED HANDLERS

// func (h *UserHandler) CreateUser(c *gin.Context) {
// 	var user models.User
// 	if err := HandleJSONBinding(c, &user, h.logger); err != nil {
// 		return
// 	}

// 	createdUser, err := h.userService.CreateUser(user)

// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	//logging
// 	LoggingResponse(c, "CreateUser", h.logger)

// 	c.JSON(http.StatusCreated, createdUser)
// }
