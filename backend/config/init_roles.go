package config

import (
	"fmt"
	"gin-backend/models"

	"gorm.io/gorm"
)

// InitRolesAndMenus 初始化角色和菜单数据
func InitRolesAndMenus(db *gorm.DB) error {
	// 自动迁移表结构
	if err := db.AutoMigrate(&models.Role{}, &models.Menu{}); err != nil {
		return fmt.Errorf("迁移角色和菜单表失败: %v", err)
	}

	// 检查是否已经初始化
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count > 0 {
		fmt.Println("角色和菜单已初始化，跳过")
		return nil
	}

	// 创建默认菜单
	menus := []models.Menu{
		{
			Name:      "首页",
			Path:      "/",
			Component: "AntdHome",
			Icon:      "HomeOutlined",
			Sort:      1,
			Type:      1,
			Status:    1,
		},
		{
			Name:      "订单管理",
			Path:      "/orders",
			Component: "Orders",
			Icon:      "ShoppingCartOutlined",
			Sort:      2,
			Type:      1,
			Status:    1,
		},
		{
			Name:      "系统管理",
			Path:      "/system",
			Component: "",
			Icon:      "SettingOutlined",
			Sort:      3,
			Type:      1,
			Status:    1,
		},
	}

	for i := range menus {
		if err := db.Create(&menus[i]).Error; err != nil {
			return fmt.Errorf("创建菜单失败: %v", err)
		}
	}

	// 创建系统管理子菜单
	systemMenuID := menus[2].ID
	subMenus := []models.Menu{
		{
			ParentID:  systemMenuID,
			Name:      "用户管理",
			Path:      "/system/users",
			Component: "SystemUsers",
			Icon:      "UserOutlined",
			Sort:      1,
			Type:      1,
			Status:    1,
		},
		{
			ParentID:  systemMenuID,
			Name:      "角色管理",
			Path:      "/system/roles",
			Component: "SystemRoles",
			Icon:      "TeamOutlined",
			Sort:      2,
			Type:      1,
			Status:    1,
		},
		{
			ParentID:  systemMenuID,
			Name:      "菜单管理",
			Path:      "/system/menus",
			Component: "SystemMenus",
			Icon:      "MenuOutlined",
			Sort:      3,
			Type:      1,
			Status:    1,
		},
	}

	for i := range subMenus {
		if err := db.Create(&subMenus[i]).Error; err != nil {
			return fmt.Errorf("创建子菜单失败: %v", err)
		}
	}

	// 收集所有菜单ID
	var allMenus []models.Menu
	db.Find(&allMenus)

	// 创建超级管理员角色
	superAdminRole := models.Role{
		Name:        "超级管理员",
		Code:        "super_admin",
		Description: "系统超级管理员，拥有所有权限",
		IsSuper:     true,
		Status:      1,
		Menus:       allMenus, // 分配所有菜单
	}

	if err := db.Create(&superAdminRole).Error; err != nil {
		return fmt.Errorf("创建超级管理员角色失败: %v", err)
	}

	// 创建普通用户角色（只有首页和订单管理）
	userRole := models.Role{
		Name:        "普通用户",
		Code:        "user",
		Description: "普通用户角色",
		IsSuper:     false,
		Status:      1,
		Menus:       menus[:2], // 只分配首页和订单管理
	}

	if err := db.Create(&userRole).Error; err != nil {
		return fmt.Errorf("创建普通用户角色失败: %v", err)
	}

	fmt.Println("角色和菜单初始化成功")
	fmt.Printf("超级管理员角色ID: %d\n", superAdminRole.ID)
	fmt.Printf("普通用户角色ID: %d\n", userRole.ID)

	return nil
}
