package routes

import (
	"gin-backend/controllers"
	"gin-backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户相关路由
func SetupUserRoutes(api *gin.RouterGroup, userController *controllers.UserController) {
	users := api.Group("/users")
	{
		// 需要认证的用户操作
		users.Use(middlewares.AuthMiddleware())
		{
			users.GET("", userController.GetUsers)          // 获取用户列表
			users.GET("/:id", userController.GetUser)       // 获取单个用户
			users.PUT("/:id", userController.UpdateUser)    // 更新用户
			users.DELETE("/:id", userController.DeleteUser) // 删除用户
		}

		// 创建用户（可能不需要认证，根据业务需求调整）
		users.POST("", userController.CreateUser)
	}
}
