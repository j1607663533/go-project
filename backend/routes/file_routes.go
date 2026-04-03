package routes

import (
	"gin-backend/controllers"
	"gin-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupFileRoutes(api *gin.RouterGroup, ctrl *controllers.FileController) {
	files := api.Group("/files")
	files.Use(middlewares.AuthMiddleware()) // 需要登录
	{
		files.POST("/upload", ctrl.Upload)
		files.GET("", ctrl.List)
		files.DELETE("/:id", ctrl.Delete)
		files.GET("/download/:id", ctrl.Download)
	}
}
