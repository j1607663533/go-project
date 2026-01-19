package models

import "time"

// Role 角色模型
type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:50"` // 角色名称
	Code        string    `json:"code" gorm:"uniqueIndex;not null;size:50"` // 角色编码，如 admin, user
	Description string    `json:"description" gorm:"size:200"`              // 角色描述
	IsSuper     bool      `json:"is_super" gorm:"default:false"`            // 是否为超级管理员
	Status      int       `json:"status" gorm:"default:1"`                  // 状态：1-启用，0-禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Menus       []Menu    `json:"menus" gorm:"many2many:role_menus;"` // 角色拥有的菜单
}

// Menu 菜单模型
type Menu struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ParentID  uint      `json:"parent_id" gorm:"default:0"`    // 父菜单ID，0表示顶级菜单
	Name      string    `json:"name" gorm:"not null;size:50"`  // 菜单名称
	Path      string    `json:"path" gorm:"not null;size:200"` // 路由路径
	Component string    `json:"component" gorm:"size:200"`     // 组件路径
	Icon      string    `json:"icon" gorm:"size:50"`           // 图标
	Sort      int       `json:"sort" gorm:"default:0"`         // 排序
	Type      int       `json:"type" gorm:"default:1"`         // 类型：1-菜单，2-按钮
	Status    int       `json:"status" gorm:"default:1"`       // 状态：1-启用，0-禁用
	Hidden    bool      `json:"hidden" gorm:"default:false"`   // 是否隐藏
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Children  []Menu    `json:"children" gorm:"-"` // 子菜单（不存储在数据库）
}

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50" validate:"required,min=2,max=50"`
	Code        string `json:"code" binding:"required,min=2,max=50,alphanum" validate:"required,min=2,max=50,alphanum"`
	Description string `json:"description" binding:"omitempty,max=200" validate:"omitempty,max=200"`
	MenuIDs     []uint `json:"menu_ids" binding:"omitempty" validate:"omitempty"`
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=50" validate:"omitempty,min=2,max=50"`
	Description string `json:"description" binding:"omitempty,max=200" validate:"omitempty,max=200"`
	Status      *int   `json:"status" binding:"omitempty,oneof=0 1" validate:"omitempty,oneof=0 1"`
	MenuIDs     []uint `json:"menu_ids" binding:"omitempty" validate:"omitempty"`
}

// MenuCreateRequest 创建菜单请求
type MenuCreateRequest struct {
	ParentID  uint   `json:"parent_id" binding:"omitempty" validate:"omitempty"`
	Name      string `json:"name" binding:"required,min=2,max=50" validate:"required,min=2,max=50"`
	Path      string `json:"path" binding:"required,max=200" validate:"required,max=200"`
	Component string `json:"component" binding:"omitempty,max=200" validate:"omitempty,max=200"`
	Icon      string `json:"icon" binding:"omitempty,max=50" validate:"omitempty,max=50"`
	Sort      int    `json:"sort" binding:"omitempty" validate:"omitempty"`
	Type      int    `json:"type" binding:"omitempty,oneof=1 2" validate:"omitempty,oneof=1 2"`
	Hidden    bool   `json:"hidden" binding:"omitempty" validate:"omitempty"`
}

// MenuUpdateRequest 更新菜单请求
type MenuUpdateRequest struct {
	ParentID  *uint  `json:"parent_id" binding:"omitempty" validate:"omitempty"`
	Name      string `json:"name" binding:"omitempty,min=2,max=50" validate:"omitempty,min=2,max=50"`
	Path      string `json:"path" binding:"omitempty,max=200" validate:"omitempty,max=200"`
	Component string `json:"component" binding:"omitempty,max=200" validate:"omitempty,max=200"`
	Icon      string `json:"icon" binding:"omitempty,max=50" validate:"omitempty,max=50"`
	Sort      *int   `json:"sort" binding:"omitempty" validate:"omitempty"`
	Type      *int   `json:"type" binding:"omitempty,oneof=1 2" validate:"omitempty,oneof=1 2"`
	Status    *int   `json:"status" binding:"omitempty,oneof=0 1" validate:"omitempty,oneof=0 1"`
	Hidden    *bool  `json:"hidden" binding:"omitempty" validate:"omitempty"`
}

// MenuTreeResponse 菜单树响应
type MenuTreeResponse struct {
	ID        uint               `json:"id"`
	ParentID  uint               `json:"parent_id"`
	Name      string             `json:"name"`
	Path      string             `json:"path"`
	Component string             `json:"component"`
	Icon      string             `json:"icon"`
	Sort      int                `json:"sort"`
	Type      int                `json:"type"`
	Hidden    bool               `json:"hidden"`
	Children  []MenuTreeResponse `json:"children,omitempty"`
}
