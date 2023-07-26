package handlers

import (
	custom_errors "edumatch/internal/app/errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetId(c *gin.Context, logger *zap.Logger) (uuid.UUID, error) {
	id := c.Param("id")
	ID, err := uuid.Parse(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Id provided"})

		//logging
		logger.Error("Invalid Id provided",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", http.StatusBadRequest),
			zap.Error(err))
		return uuid.Nil, err
	}
	return ID, nil
}

func HandleJSONBinding(c *gin.Context, target interface{}) error {
	if err := c.ShouldBindJSON(&target); err != nil {
		err = custom_errors.ErrHandleBinding
		return err
	}
	return nil
}

func HandleFormDataBinding(c *gin.Context, target interface{}) error {
	if err := c.ShouldBind(target); err != nil {
		fmt.Println(err)
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
