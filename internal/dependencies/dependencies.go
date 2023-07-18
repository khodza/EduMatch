package dependencies

import (
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/handlers"
	"edumatch/internal/app/repositories"
	"edumatch/internal/app/services"
	"edumatch/internal/app/validators"
	database "edumatch/pkg/db"
	"edumatch/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

func InitDependencies() (*custom_errors.GlobalErrorHandler, map[string]interface{}, *zap.Logger, error) {
	// INITIALIZE LOGGER
	logger, err := logger.CreateLogger()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error on initializing logger")
	}
	// Get the DB instance
	db := database.GetDB()
	if db == nil {
		return nil, nil, nil, fmt.Errorf("error on getting db instance")
	}

	// INITIALIZE REPOSITORIES
	userRepository := repositories.NewUserRepository(db)
	eduCenterRepository := repositories.NewEduCenterRepository(db)

	//INITIALIZE VALIDATORS
	userValidator := validators.NewUserValidator()
	eduCenterValidator := validators.NewEduCenterValidator()

	// INITIALIZE SERVICES
	userService := services.NewUserService(userRepository, userValidator)
	authService := services.NewAuthService(userService)
	eduCenterService := services.NewEduCenterService(eduCenterRepository, eduCenterValidator)

	// INITIALIZE HANDLERS
	userHandler := handlers.NewUserHandler(userService, logger)
	eduCenterHandler := handlers.NewEduCenterHandler(eduCenterService, logger)
	authHandler := handlers.NewAuthHandler(authService, logger)

	//INITIALIZE Global Error Handler
	globalErrorHandler := custom_errors.NewGlobalErrorHandler(logger)

	handlersMap := map[string]interface{}{
		"users":      userHandler,
		"eduCenters": eduCenterHandler,
		"auth":       authHandler,
	}

	return globalErrorHandler, handlersMap, logger, nil
}
