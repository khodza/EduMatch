package handlers

import (
	custom_errors "edumatch/internal/app/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetId(c *gin.Context, logger *zap.Logger) (uuid.UUID, error) {
	id := c.Param("id")
	ID, err := uuid.Parse(id)

	if err != nil {
		//logging
		logger.Error(custom_errors.ErrInvalidID.Error(),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", http.StatusBadRequest),
			zap.Error(err))
		return uuid.Nil, custom_errors.ErrInvalidID
	}
	return ID, nil
}

func HandleJSONBinding(c *gin.Context, target interface{}, logger *zap.Logger) error {
	if err := c.ShouldBindJSON(&target); err != nil {
		//logging
		logger.Error(custom_errors.ErrHandleBinding.Error(),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", http.StatusBadRequest),
			zap.Error(err))
		return custom_errors.ErrHandleBinding
	}
	return nil
}

func HandleFormDataBinding(c *gin.Context, target interface{}, logger *zap.Logger) error {
	if err := c.ShouldBind(target); err != nil {
		//logging
		logger.Error(custom_errors.ErrHandleBinding.Error(),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", http.StatusBadRequest),
			zap.Error(err))
		return custom_errors.ErrHandleBinding
	}
	return nil
}

func LoggingResponse(c *gin.Context, info string, logger *zap.Logger) {
	logger.Info(info,
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", http.StatusOK),
	)
}
