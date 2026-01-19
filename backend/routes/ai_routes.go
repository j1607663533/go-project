package routes

import (
	"gin-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAIRoutes(r *gin.RouterGroup) {
	aiController := controllers.NewAIController()

	aiGroup := r.Group("/ai")
	{
		aiGroup.POST("/chat", aiController.Chat)
	}
}
