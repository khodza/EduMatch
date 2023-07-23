package handlers

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandlerInterface interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	ProtectedEndpoint(roles ...models.Role) gin.HandlerFunc
	RefreshToken(c *gin.Context)
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

func (h *AuthHandler) ProtectedEndpoint(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		userID, role, err := services.ValidateToken(token)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("user_role", role)

		if len(roles) > 0 {
			// Check if the user has the required role
			authorized := false
			for _, requiredRole := range roles {
				if role == requiredRole {
					authorized = true
					break
				}
			}

			if !authorized {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Retrieve the refresh token from the request body or header
	refreshToken := c.PostForm("refresh_token")
	if refreshToken == "" {
		refreshToken = c.GetHeader("Authorization")
	}

	// Validate the refresh token
	userID, userRole, err := services.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate a new JWT token
	token, err := services.GenerateToken(userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
