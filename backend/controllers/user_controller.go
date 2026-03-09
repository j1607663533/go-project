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
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// UserController 用户控制器
type UserController struct {
	userService   services.UserService
	wechatService services.WechatService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService services.UserService, wechatService services.WechatService) *UserController {
	return &UserController{
		userService:   userService,
		wechatService: wechatService,
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
// @Security Bearer
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
// @Summary 更新用户信息
// @Description 修改指定用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param user body models.UserUpdateRequest true "更新信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,409 {object} map[string]interface{}
// @Router /users/{id} [put]
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
// @Summary 删除指定用户
// @Description 通过用户ID删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /users/{id} [delete]
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

// GetProfile 获取当前用户信息
// @Summary 获取个人资料
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401,404 {object} map[string]interface{}
// @Router /users/profile [get]
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

// GetWeChatQRCode 获取微信登录二维码
func (ctrl *UserController) GetWeChatQRCode(c *gin.Context) {
	session, qrURL, err := ctrl.wechatService.GetQRCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取二维码失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"qr_url":    qrURL,
			"scene_id":  session.SceneID,
			"expire_at": session.ExpireAt.UnixMilli(),
		},
	})
}

// CheckWeChatStatus 检查微信登录状态（轮询接口）
func (ctrl *UserController) CheckWeChatStatus(c *gin.Context) {
	sceneID := c.Query("scene_id")
	if sceneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "scene_id 不能为空"})
		return
	}

	session, err := ctrl.wechatService.CheckStatus(sceneID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 如果状态是成功，则返回用户 Token
	if session.Status == services.StatusSuccess {
		loginResp, err := ctrl.userService.LoginByUserID(session.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "授权登录失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"status": "SUCCESS",
				"token":  loginResp.Token,
				"user":   loginResp.User,
				"menus":  loginResp.Menus,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"status": session.Status,
		},
	})
}

// MockWeChatScan 模拟微信扫码成功（仅供开发测试使用）
func (ctrl *UserController) MockWeChatScan(c *gin.Context) {
	sceneID := c.Query("scene_id")
	userIDStr := c.DefaultQuery("user_id", "1") // 默认为管理员用户
	userID, _ := strconv.Atoi(userIDStr)

	err := ctrl.wechatService.MockScan(sceneID, uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "模拟扫码授权成功，请查看前端状态",
	})
}

// WechatCallback 处理微信服务端推送的消息
func (ctrl *UserController) WechatCallback(c *gin.Context) {
	// 获取微信 SDK 的 Server 对象
	server := ctrl.wechatService.GetServer(c.Request, c.Writer)
	if server == nil {
		c.String(http.StatusOK, "success") // 或者返回错误，视情况而定
		return
	}

	// 设置消息处理回调
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		return ctrl.wechatService.HandleCallback(*msg)
	})

	// 处理请求并发送响应
	err := server.Serve()
	if err != nil {
		fmt.Printf("微信回调处理失败: %v\n", err)
		return
	}
}
