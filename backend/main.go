package main

import (
	"gin-backend/config"
	"gin-backend/models"
	"gin-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer config.CloseDB()

	// 初始化 Redis 连接
	if err := config.InitRedis(); err != nil {
		log.Printf("Redis 连接失败: %v (将使用内存存储作为降级方案)", err)
	} else {
		defer config.CloseRedis()
	}

	// 自动迁移数据库表
	log.Println("开始数据库迁移...")
	if err := config.AutoMigrate(&models.User{}, &models.Payment{}, &models.Order{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")

	// 初始化角色和菜单
	log.Println("开始初始化角色和菜单...")
	if err := config.InitRolesAndMenus(config.DB); err != nil {
		log.Fatalf("角色和菜单初始化失败: %v", err)
	}

	// 设置 Gin 模式
	if config.AppConfig.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由，传入数据库连接
	r := routes.SetupRouter(config.DB)

	// 启动服务器
	addr := config.AppConfig.Host + ":" + config.AppConfig.Port
	log.Printf("服务器启动在: %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
