package services

import (
	"errors"
	"gin-backend/models"
	"gin-backend/repositories"
)

// RoleService 角色服务接口
type RoleService interface {
	CreateRole(req *models.RoleCreateRequest) error
	UpdateRole(id uint, req *models.RoleUpdateRequest) error
	DeleteRole(id uint) error
	GetRoleByID(id uint) (*models.Role, error)
	GetAllRoles() ([]models.Role, error)
	AssignMenus(roleID uint, menuIDs []uint) error
}

// roleService 角色服务实现
type roleService struct {
	roleRepo repositories.RoleRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo repositories.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

// CreateRole 创建角色
func (s *roleService) CreateRole(req *models.RoleCreateRequest) error {
	role := &models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsSuper:     false, // 新创建的角色不能是超级管理员
		Status:      1,
	}

	// 创建角色
	if err := s.roleRepo.Create(role); err != nil {
		return err
	}

	// 分配菜单
	if len(req.MenuIDs) > 0 {
		return s.roleRepo.AssignMenus(role.ID, req.MenuIDs)
	}

	return nil
}

// UpdateRole 更新角色
func (s *roleService) UpdateRole(id uint, req *models.RoleUpdateRequest) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 超级管理员角色不能修改
	if role.IsSuper {
		return errors.New("超级管理员角色不能修改")
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}

	// 更新角色
	if err := s.roleRepo.Update(role); err != nil {
		return err
	}

	// 更新菜单关联
	if req.MenuIDs != nil {
		return s.roleRepo.AssignMenus(id, req.MenuIDs)
	}

	return nil
}

// DeleteRole 删除角色
func (s *roleService) DeleteRole(id uint) error {
	return s.roleRepo.Delete(id)
}

// GetRoleByID 根据ID获取角色
func (s *roleService) GetRoleByID(id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(id)
}

// GetAllRoles 获取所有角色
func (s *roleService) GetAllRoles() ([]models.Role, error) {
	return s.roleRepo.FindAll()
}

// AssignMenus 为角色分配菜单
func (s *roleService) AssignMenus(roleID uint, menuIDs []uint) error {
	return s.roleRepo.AssignMenus(roleID, menuIDs)
}
