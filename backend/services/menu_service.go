package services

import (
	"gin-backend/models"
	"gin-backend/repositories"
)

// MenuService 菜单服务接口
type MenuService interface {
	CreateMenu(req *models.MenuCreateRequest) error
	UpdateMenu(id uint, req *models.MenuUpdateRequest) error
	DeleteMenu(id uint) error
	GetMenuByID(id uint) (*models.Menu, error)
	GetAllMenus() ([]models.Menu, error)
	GetMenuTree() ([]models.MenuTreeResponse, error)
	GetUserMenus(userID uint) ([]models.MenuTreeResponse, error)
}

// menuService 菜单服务实现
type menuService struct {
	menuRepo repositories.MenuRepository
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

// NewMenuService 创建菜单服务实例
func NewMenuService(menuRepo repositories.MenuRepository, userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateMenu 创建菜单
func (s *menuService) CreateMenu(req *models.MenuCreateRequest) error {
	menu := &models.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Type:      req.Type,
		Hidden:    req.Hidden,
		Status:    1,
	}

	// 创建菜单
	if err := s.menuRepo.Create(menu); err != nil {
		return err
	}

	// 自动将新菜单分配给超级管理员角色
	// 使用 GORM 的 Association 功能，直接操作主表关联
	superAdminRole, err := s.roleRepo.FindByCode("super_admin")
	if err == nil && superAdminRole != nil {
		// 使用 GORM Association Append 添加菜单到角色
		// 这会自动在 role_menus 中间表中创建关联记录
		_ = s.roleRepo.GetDB().Model(superAdminRole).Association("Menus").Append(menu)
	}

	return nil
}

// UpdateMenu 更新菜单
func (s *menuService) UpdateMenu(id uint, req *models.MenuUpdateRequest) error {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return err
	}

	if req.ParentID != nil {
		menu.ParentID = *req.ParentID
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Sort != nil {
		menu.Sort = *req.Sort
	}
	if req.Type != nil {
		menu.Type = *req.Type
	}
	if req.Status != nil {
		menu.Status = *req.Status
	}
	if req.Hidden != nil {
		menu.Hidden = *req.Hidden
	}

	return s.menuRepo.Update(menu)
}

// DeleteMenu 删除菜单
func (s *menuService) DeleteMenu(id uint) error {
	return s.menuRepo.Delete(id)
}

// GetMenuByID 根据ID获取菜单
func (s *menuService) GetMenuByID(id uint) (*models.Menu, error) {
	return s.menuRepo.FindByID(id)
}

// GetAllMenus 获取所有菜单
func (s *menuService) GetAllMenus() ([]models.Menu, error) {
	return s.menuRepo.FindAll()
}

// GetMenuTree 获取菜单树
func (s *menuService) GetMenuTree() ([]models.MenuTreeResponse, error) {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.menuRepo.BuildMenuTree(menus), nil
}

// GetUserMenus 获取用户菜单
func (s *menuService) GetUserMenus(userID uint) ([]models.MenuTreeResponse, error) {
	// 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 获取用户角色的菜单
	menus, err := s.menuRepo.FindByRoleID(user.RoleID)
	if err != nil {
		return nil, err
	}

	return s.menuRepo.BuildMenuTree(menus), nil
}
