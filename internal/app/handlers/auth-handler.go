package handlers

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandlerInterface interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	ProtectedEndpoint(c *gin.Context)
}

type AuthHandler struct {
	authService services.AuthServiceInterface
	logger      *zap.Logger
}

func NewAuthHandler(authService services.AuthServiceInterface, logger *zap.Logger) AuthHandlerInterface {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var user models.RegUser
	if err := HandleJSONBinding(c, &user, h.logger); err != nil {
		c.Error(err)
		return
	}
	// register user
	tokens, err := h.authService.RegisterUser(user)
	// ...
	if err != nil {
		c.Error(err)
		return
	}

	//logging
	LoggingResponse(c, "SignUp", h.logger)

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loggingUser models.LoggingUser
	if err := HandleJSONBinding(c, &loggingUser, h.logger); err != nil {
		c.Error(err)
		return
	}
	// Generate a new JWT token
	tokens, err := h.authService.Login(loggingUser)
	if err != nil {
		c.Error(err)
		return
	}

	// Logging
	LoggingResponse(c, "Login", h.logger)

	c.JSON(http.StatusOK, tokens)
}

// ProtectedEndpoint is an example of a protected endpoint that requires authentication
func (h *AuthHandler) ProtectedEndpoint(c *gin.Context) {
	fmt.Println("PROTECTED")
	// Retrieve the token from the request header or query parameter
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}

	// Validate the token
	userID, err := services.ValidateToken(token)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Set("user_id", userID)
	c.Next()
}

// // RefreshToken handles the refresh token request and generates a new JWT token
// func (h *AuthHandler) RefreshToken(c *gin.Context) {
// 	// Retrieve the refresh token from the request body or header
// 	refreshToken := c.PostForm("refresh_token")
// 	if refreshToken == "" {
// 		refreshToken = c.GetHeader("Authorization")
// 	}

// 	// Validate the refresh token
// 	userID, err := services.ValidateRefreshToken(refreshToken)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
// 		return
// 	}

// 	// Generate a new JWT token
// 	token, err := services.GenerateToken(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }
