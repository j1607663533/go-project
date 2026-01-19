package controllers

import (
	"gin-backend/models"
	"gin-backend/services"
	"gin-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	roleService services.RoleService
}

// NewRoleController 创建角色控制器实例
func NewRoleController(roleService services.RoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

// GetAllRoles 获取所有角色
func (ctrl *RoleController) GetAllRoles(c *gin.Context) {
	roles, err := ctrl.roleService.GetAllRoles()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取角色列表失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "获取角色列表成功", roles)
}

// GetRoleByID 根据ID获取角色
func (ctrl *RoleController) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	role, err := ctrl.roleService.GetRoleByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取角色失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "获取角色成功", role)
}

// CreateRole 创建角色
func (ctrl *RoleController) CreateRole(c *gin.Context) {
	var req models.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := ctrl.roleService.CreateRole(&req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建角色失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "创建角色成功", nil)
}

// UpdateRole 更新角色
func (ctrl *RoleController) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req models.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := ctrl.roleService.UpdateRole(uint(id), &req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新角色失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "更新角色成功", nil)
}

// DeleteRole 删除角色
func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	if err := ctrl.roleService.DeleteRole(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除角色失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "删除角色成功", nil)
}

// AssignMenus 为角色分配菜单
func (ctrl *RoleController) AssignMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		MenuIDs []uint `json:"menu_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := ctrl.roleService.AssignMenus(uint(id), req.MenuIDs); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "分配菜单失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "分配菜单成功", nil)
}
