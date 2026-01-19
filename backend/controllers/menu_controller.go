package controllers

import (
	"gin-backend/models"
	"gin-backend/services"
	"gin-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuController 菜单控制器
type MenuController struct {
	menuService services.MenuService
}

// NewMenuController 创建菜单控制器实例
func NewMenuController(menuService services.MenuService) *MenuController {
	return &MenuController{menuService: menuService}
}

// GetMenuTree 获取菜单树
func (ctrl *MenuController) GetMenuTree(c *gin.Context) {
	menus, err := ctrl.menuService.GetMenuTree()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取菜单树失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "获取菜单树成功", menus)
}

// GetUserMenus 获取当前用户菜单
func (ctrl *MenuController) GetUserMenus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	menus, err := ctrl.menuService.GetUserMenus(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用户菜单失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "获取用户菜单成功", menus)
}

// CreateMenu 创建菜单
func (ctrl *MenuController) CreateMenu(c *gin.Context) {
	var req models.MenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := ctrl.menuService.CreateMenu(&req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建菜单失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "创建菜单成功", nil)
}

// UpdateMenu 更新菜单
func (ctrl *MenuController) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	var req models.MenuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := ctrl.menuService.UpdateMenu(uint(id), &req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新菜单失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "更新菜单成功", nil)
}

// DeleteMenu 删除菜单
func (ctrl *MenuController) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	if err := ctrl.menuService.DeleteMenu(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除菜单失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(c, "删除菜单成功", nil)
}
