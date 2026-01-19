package routes

import (
	"gin-backend/controllers"
	"gin-backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupMenuRoutes 设置菜单路由
func SetupMenuRoutes(api *gin.RouterGroup, menuController *controllers.MenuController) {
	menus := api.Group("/menus")
	{
		// 公开路由（需要认证）
		menus.Use(middlewares.AuthMiddleware())
		menus.GET("/user", menuController.GetUserMenus) // 获取当前用户菜单
		menus.GET("/tree", menuController.GetMenuTree)  // 获取菜单树

		// 管理路由（需要管理员权限）
		menus.POST("", menuController.CreateMenu)       // 创建菜单
		menus.PUT("/:id", menuController.UpdateMenu)    // 更新菜单
		menus.DELETE("/:id", menuController.DeleteMenu) // 删除菜单
	}
}

// SetupRoleRoutes 设置角色路由
func SetupRoleRoutes(api *gin.RouterGroup, roleController *controllers.RoleController) {
	roles := api.Group("/roles")
	{
		// 所有路由都需要认证
		roles.Use(middlewares.AuthMiddleware())

		roles.GET("", roleController.GetAllRoles)            // 获取所有角色
		roles.GET("/:id", roleController.GetRoleByID)        // 根据ID获取角色
		roles.POST("", roleController.CreateRole)            // 创建角色
		roles.PUT("/:id", roleController.UpdateRole)         // 更新角色
		roles.DELETE("/:id", roleController.DeleteRole)      // 删除角色
		roles.POST("/:id/menus", roleController.AssignMenus) // 为角色分配菜单
	}
}
