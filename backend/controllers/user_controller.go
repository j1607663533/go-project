package controllers

import (
	"fmt"
	"gin-backend/models"
	"gin-backend/services"
	"gin-backend/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUsers 获取用户列表
// @Summary 获取所有用户
// @Description 获取系统中所有注册用户的列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func (ctrl *UserController) GetUsers(c *gin.Context) {

	fmt.Println("获取用户列表")

	users, err := ctrl.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    users,
	})
}

// GetUser 获取单个用户
// @Summary 获取指定ID的用户
// @Description 通过用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Router /users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

// CreateUser 创建用户
// @Summary 创建新用户
// @Description 手动创建一个新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "用户信息"
// @Success 201 {object} map[string]interface{}
// @Failure 400,409,500 {object} map[string]interface{}
// @Router /users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"errors":  utils.FormatValidationErrors(err),
		})
		return
	}

	user, err := ctrl.userService.CreateUser(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已存在" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "用户创建成功",
		"data":    user,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户通过手机/邮箱进行注册，包含验证码校验
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param register body object true "注册信息"
// @Success 201 {object} map[string]interface{}
// @Failure 400,409,500 {object} map[string]interface{}
// @Router /register [post]
func (ctrl *UserController) Register(c *gin.Context) {
	var req struct {
		models.UserCreateRequest
		CaptchaID string `json:"captcha_id" binding:"required"`
		Captcha   string `json:"captcha" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"errors":  utils.FormatValidationErrors(err),
		})
		return
	}

	// 验证验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "验证码错误或已过期",
		})
		return
	}

	// 创建用户
	user, err := ctrl.userService.CreateUser(&req.UserCreateRequest)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已存在" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "注册成功",
		"data":    user,
	})
}

// UpdateUser 更新用户
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	fmt.Println("更新用户", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"errors":  utils.FormatValidationErrors(err),
		})
		return
	}

	user, err := ctrl.userService.UpdateUser(uint(id), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "用户不存在" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "邮箱已被使用" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户更新成功",
		"data":    user,
	})
}

// DeleteUser 删除用户
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	err = ctrl.userService.DeleteUser(uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "用户不存在" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户删除成功",
	})
}

// GetProfile 获取当前用户信息（需要认证）
func (ctrl *UserController) GetProfile(c *gin.Context) {
	// 从上下文中获取当前用户ID（由认证中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	user, err := ctrl.userService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码获取 JWT Token
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func (ctrl *UserController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"errors":  utils.FormatValidationErrors(err),
		})
		return
	}

	loginResp, err := ctrl.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    loginResp,
	})
}

// ExitLogin 退出登录
func (ctrl *UserController) ExitLogin(c *gin.Context) {
	// 从请求头获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供认证令牌",
		})
		return
	}

	// 验证 token 格式 (Bearer token)
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "认证令牌格式错误",
		})
		return
	}

	token := parts[1]

	// 调用服务层退出登录
	err := ctrl.userService.ExitLogin(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "退出登录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "退出登录成功",
	})
}
