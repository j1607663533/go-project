package repositories

import (
	"gin-backend/models"

	"gorm.io/gorm"
)

// MenuRepository 菜单仓库接口
type MenuRepository interface {
	Create(menu *models.Menu) error
	Update(menu *models.Menu) error
	Delete(id uint) error
	FindByID(id uint) (*models.Menu, error)
	FindAll() ([]models.Menu, error)
	FindByRoleID(roleID uint) ([]models.Menu, error)
	BuildMenuTree(menus []models.Menu) []models.MenuTreeResponse
}

// menuRepository 菜单仓库实现
type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓库实例
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

// Create 创建菜单
func (r *menuRepository) Create(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

// Update 更新菜单
func (r *menuRepository) Update(menu *models.Menu) error {
	return r.db.Save(menu).Error
}

// Delete 删除菜单
func (r *menuRepository) Delete(id uint) error {
	return r.db.Delete(&models.Menu{}, id).Error
}

// FindByID 根据ID查找菜单
func (r *menuRepository) FindByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.First(&menu, id).Error
	return &menu, err
}

// FindAll 查找所有菜单
func (r *menuRepository) FindAll() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Where("status = ?", 1).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

// FindByRoleID 根据角色ID查找菜单
func (r *menuRepository) FindByRoleID(roleID uint) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Table("menus").
		Joins("INNER JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND menus.status = ? AND menus.type = ?", roleID, 1, 1).
		Order("menus.sort ASC, menus.id ASC").
		Find(&menus).Error
	return menus, err
}

// BuildMenuTree 构建菜单树
func (r *menuRepository) BuildMenuTree(menus []models.Menu) []models.MenuTreeResponse {
	menuMap := make(map[uint]*models.MenuTreeResponse)
	var roots []*models.MenuTreeResponse

	// 第一遍：创建所有菜单节点
	for _, menu := range menus {
		menuMap[menu.ID] = &models.MenuTreeResponse{
			ID:        menu.ID,
			ParentID:  menu.ParentID,
			Name:      menu.Name,
			Path:      menu.Path,
			Component: menu.Component,
			Icon:      menu.Icon,
			Sort:      menu.Sort,
			Type:      menu.Type,
			Hidden:    menu.Hidden,
			Children:  []models.MenuTreeResponse{},
		}
	}

	// 第二遍：构建树形结构
	for _, menu := range menus {
		if menu.ParentID == 0 {
			// 顶级菜单
			roots = append(roots, menuMap[menu.ID])
		} else {
			// 子菜单
			if parent, ok := menuMap[menu.ParentID]; ok {
				parent.Children = append(parent.Children, *menuMap[menu.ID])
			}
		}
	}

	// 转换为值类型返回
	result := make([]models.MenuTreeResponse, len(roots))
	for i, root := range roots {
		result[i] = *root
	}

	return result
}
