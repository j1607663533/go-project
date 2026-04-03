package routes

import (
	"gin-backend/controllers"
	"gin-backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupLotteryRoutes 设置抽奖相关路由
func SetupLotteryRoutes(r *gin.RouterGroup, controller *controllers.LotteryController) {
	lotteryGroup := r.Group("/lottery")
	// 需要登录的接口
	lotteryGroup.Use(middlewares.AuthMiddleware())
	{
		lotteryGroup.GET("/info", controller.GetInfo)
		lotteryGroup.POST("/draw", controller.Draw)
		lotteryGroup.GET("/records/my", controller.GetRecords)
		lotteryGroup.GET("/records/public", controller.GetPublicRecords)

		// 管理后台接口
		lotteryGroup.GET("/admin/activities", controllers.AdminGetActivities)
		lotteryGroup.POST("/admin/activities", controllers.AdminSaveConfig)
		lotteryGroup.PUT("/admin/activities/:id/status", controllers.AdminToggleStatus)
		lotteryGroup.DELETE("/admin/activities/:id", controllers.AdminDeleteActivity)
		
		// 抽奖统计流水
		lotteryGroup.GET("/admin/records", controllers.AdminGetRecords)
	}
}
