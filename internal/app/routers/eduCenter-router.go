package routers

import (
	"edumatch/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupEduCenter(router *gin.RouterGroup, eduCenterHandler handlers.EduCenterHandlerInterface) {
	eduCenterGroup := router.Group("")
	eduCenterGroup.GET("/", eduCenterHandler.GetEduCenters)
	eduCenterGroup.POST("/", eduCenterHandler.CreateEduCenter)
	eduCenterGroup.GET("/:id", eduCenterHandler.GetEduCenter)
	eduCenterGroup.PATCH("/:id", eduCenterHandler.UpdateEduCenter)
	eduCenterGroup.DELETE("/:id", eduCenterHandler.DeleteEduCenter)
}
