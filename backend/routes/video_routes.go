package routes

import (
	"gin-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupVideoRoutes(api *gin.RouterGroup, ctrl *controllers.VideoController) {
	video := api.Group("/video")
	{
		video.POST("/remove-watermark", ctrl.RemoveWatermark)
	}
}
