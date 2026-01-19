package routes

import (
	"gin-backend/controllers"
	"gin-backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes 设置订单相关路由
func SetupOrderRoutes(api *gin.RouterGroup, orderController *controllers.OrderController) {
	orders := api.Group("/orders")
	// 所有订单操作都需要认证
	orders.Use(middlewares.AuthMiddleware())
	{
		orders.GET("/", orderController.GetOrderListWithPage) // 分页列表 (匹配 /orders/ )
		orders.GET("", orderController.GetOrderListWithPage)  // 分页列表 (匹配 /orders )
		orders.GET("/all", orderController.GetOrderList)      // 全部列表
		orders.GET("/:id", orderController.GetOrderById)      // 详情
		orders.POST("/", orderController.CreateOrder)         // 创建 (带斜杠)
		orders.POST("", orderController.CreateOrder)          // 创建 (不带斜杠)
		orders.PUT("/:id", orderController.UpdateOrder)       // 更新
		orders.DELETE("/:id", orderController.DeleteOrder)    // 删除
	}

}
