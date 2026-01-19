package routes

import (
	"gin-backend/controllers"
	"gin-backend/middlewares"
	"gin-backend/repositories"
	"gin-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 禁用自动重定向，防止 301 问题
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// 使用 CORS 中间件
	r.Use(middlewares.CORS())

	// 使用请求日志中间件
	r.Use(middlewares.RequestLogger())

	// 初始化依赖注入
	// Repository 层 - 注入数据库连接
	userRepo := repositories.NewUserRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	menuRepo := repositories.NewMenuRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	// Service 层 - 注入 Repository
	userService := services.NewUserService(userRepo, menuRepo)
	orderService := services.NewOrderService(orderRepo)
	menuService := services.NewMenuService(menuRepo, userRepo, roleRepo)
	roleService := services.NewRoleService(roleRepo)

	// Controller 层 - 注入 Service
	userController := controllers.NewUserController(userService)
	orderController := controllers.NewOrderController(orderService)
	menuController := controllers.NewMenuController(menuService)
	roleController := controllers.NewRoleController(roleService)
	captchaController := controllers.NewCaptchaController()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// API 路由组
	api := r.Group("/api/v1")

	// 设置各模块路由
	SetupAuthRoutes(api, userController, captchaController) // 认证路由
	SetupUserRoutes(api, userController)                    // 用户路由
	SetupOrderRoutes(api, orderController)                  // 订单路由
	SetupMenuRoutes(api, menuController)                    // 菜单路由
	SetupRoleRoutes(api, roleController)                    // 角色路由

	return r
}
