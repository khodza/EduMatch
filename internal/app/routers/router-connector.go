package routers

import (
	"edumatch/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func ConnectRoutersToHandlers(router *gin.Engine, handlersMap map[string]interface{}) {
	for route, handler := range handlersMap {
		routeGroup := router.Group("/" + route)
		switch route {
		case "users":
			userHandler := handler.(handlers.UserHandlerInterface)
			SetupUserRouter(routeGroup, userHandler)
		case "eduCenters":
			eduCenterHandler := handler.(handlers.EduCenterHandlerInterface)
			SetupEduCenter(routeGroup, eduCenterHandler)
		case "auth":
			authHandler := handler.(handlers.AuthHandlerInterface)
			SetupAuthRouter(routeGroup, authHandler)
		}
	}
}
