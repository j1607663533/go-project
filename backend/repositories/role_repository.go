package repositories

import (
	"errors"
	"gin-backend/models"

	"gorm.io/gorm"
)

// RoleRepository 角色仓库接口
type RoleRepository interface {
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uint) error
	FindByID(id uint) (*models.Role, error)
	FindAll() ([]models.Role, error)
	FindByCode(code string) (*models.Role, error)
	AssignMenus(roleID uint, menuIDs []uint) error
	GetDB() *gorm.DB
}

// roleRepository 角色仓库实现
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库实例
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create 创建角色
func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *roleRepository) Delete(id uint) error {
	// 检查是否为超级管理员
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return err
	}
	if role.IsSuper {
		return errors.New("超级管理员角色不能删除")
	}
	return r.db.Delete(&models.Role{}, id).Error
}

// FindByID 根据ID查找角色
func (r *roleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Menus").First(&role, id).Error
	return &role, err
}

// FindAll 查找所有角色
func (r *roleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("Menus").Find(&roles).Error
	return roles, err
}

// FindByCode 根据编码查找角色
func (r *roleRepository) FindByCode(code string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Menus").Where("code = ?", code).First(&role).Error
	return &role, err
}

// AssignMenus 为角色分配菜单
func (r *roleRepository) AssignMenus(roleID uint, menuIDs []uint) error {
	// 检查是否为超级管理员
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}
	if role.IsSuper {
		return errors.New("超级管理员角色的菜单不能修改")
	}

	// 先清除现有关联
	if err := r.db.Model(&role).Association("Menus").Clear(); err != nil {
		return err
	}

	// 添加新的关联
	if len(menuIDs) > 0 {
		var menus []models.Menu
		if err := r.db.Find(&menus, menuIDs).Error; err != nil {
			return err
		}
		return r.db.Model(&role).Association("Menus").Append(menus)
	}

	return nil
}

// GetDB 获取数据库实例
func (r *roleRepository) GetDB() *gorm.DB {
	return r.db
}
