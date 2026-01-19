package routes

import (
	"gin-backend/controllers"
	"gin-backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 设置认证相关路由
func SetupAuthRoutes(api *gin.RouterGroup, userController *controllers.UserController, captchaController *controllers.CaptchaController) {
	// 验证码接口（不需要认证）
	api.GET("/captcha", captchaController.GetCaptcha)            // 获取/刷新验证码
	api.POST("/captcha/verify", captchaController.VerifyCaptcha) // 验证验证码（测试用）

	// 认证接口（不需要认证）
	api.POST("/register", userController.Register) // 用户注册
	api.POST("/login", userController.Login)       // 用户登录
	api.POST("/logout", userController.ExitLogin)  // 用户退出登录

	// 需要认证的路由
	auth := api.Group("/auth")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.GET("/profile", userController.GetProfile) // 获取当前用户信息
	}
}
