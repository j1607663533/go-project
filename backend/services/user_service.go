package services

import (
	"errors"
	"fmt"
	"gin-backend/models"
	"gin-backend/repositories"
	"gin-backend/utils"
	"time"
)

// UserService 用户业务逻辑接口
type UserService interface {
	GetAllUsers() ([]models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error)
	UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
	DeleteUser(id uint) error
	GetProfile(userID uint) (*models.UserResponse, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	ExitLogin(token string) error
	// 批量操作（演示 make 和 Channel）
	GetUsersByIDs(ids []uint) ([]models.UserResponse, error)
	BatchCreateUsers(requests []*models.UserCreateRequest) ([]models.UserResponse, []error)
}

// userService 用户业务逻辑实现
type userService struct {
	userRepo repositories.UserRepository
	menuRepo repositories.MenuRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository, menuRepo repositories.MenuRepository) UserService {
	return &userService{
		userRepo: userRepo,
		menuRepo: menuRepo,
	}
}

// GetAllUsers 获取所有用户
func (s *userService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var response []models.UserResponse
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	return response, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*models.UserResponse, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("user:%d", id)
	var cachedUser models.UserResponse
	err := utils.CacheGet(cacheKey, &cachedUser)
	if err == nil {
		// 缓存命中
		return &cachedUser, nil
	}

	// 缓存未命中，从数据库查询
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()

	// 存入缓存，过期时间 30 分钟
	utils.CacheSet(cacheKey, response, 30*time.Minute)

	return &response, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error) {
	// 业务逻辑：检查用户名是否已存在
	if existingUser, _ := s.userRepo.FindByUsername(req.Username); existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 业务逻辑：检查邮箱是否已存在
	if existingUser, _ := s.userRepo.FindByEmail(req.Email); existingUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	// 业务逻辑：对密码进行加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户实体
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Nickname:  req.Nickname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 调用仓储层保存数据
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	// 业务逻辑：检查用户是否存在
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 业务逻辑：如果更新邮箱，检查邮箱是否已被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		if existingUser, _ := s.userRepo.FindByEmail(req.Email); existingUser != nil && existingUser.ID != id {
			return nil, errors.New("邮箱已被使用")
		}
		user.Email = req.Email
	}

	// 更新字段
	if req.RoleID > 0 {
		user.RoleID = req.RoleID
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	user.UpdatedAt = time.Now()

	// 调用仓储层更新数据
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	// 业务逻辑：检查用户是否存在
	if _, err := s.userRepo.FindByID(id); err != nil {
		return err
	}

	// 业务逻辑：这里可以添加其他检查，比如是否有关联数据等

	// 调用仓储层删除数据
	return s.userRepo.Delete(id)
}

// GetProfile 获取用户个人信息
func (s *userService) GetProfile(userID uint) (*models.UserResponse, error) {
	// 业务逻辑：获取当前用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// Login 用户登录
func (s *userService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// 业务逻辑：验证验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		return nil, errors.New("验证码错误或已过期")
	}

	// 业务逻辑：查找用户并预加载角色信息
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil || user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 业务逻辑：验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 业务逻辑：生成 JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, errors.New("生成 token 失败")
	}

	// 单点登录：设置用户的当前 token，旧 token 会被自动加入黑名单
	if err := utils.SetUserToken(user.ID, token); err != nil {
		return nil, errors.New("设置用户 token 失败")
	}

	// 获取用户菜单
	var menus []models.MenuTreeResponse
	if user.RoleID > 0 {
		userMenus, err := s.menuRepo.FindByRoleID(user.RoleID)
		if err == nil {
			menus = s.menuRepo.BuildMenuTree(userMenus)
		}
	}

	// 返回登录响应
	return &models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
		Menus: menus,
	}, nil

}

// ExitLogin 退出登录
func (s *userService) ExitLogin(token string) error {
	// 业务逻辑：验证 token
	if _, err := utils.ParseToken(token); err != nil {
		return errors.New("token 无效")
	}

	// 清除 token 缓存
	utils.ClearToken(token)

	return nil
}

// GetUsersByIDs 批量获取用户（演示 make 创建切片和 Channel 并发）
func (s *userService) GetUsersByIDs(ids []uint) ([]models.UserResponse, error) {
	// 使用 make 创建结果切片
	results := make([]models.UserResponse, 0, len(ids))

	// 使用 make 创建 Channel
	type userResult struct {
		user *models.UserResponse
		err  error
	}
	resultChan := make(chan userResult, len(ids))

	// 并发获取用户信息
	for _, id := range ids {
		go func(userID uint) {
			user, err := s.GetUserByID(userID)
			resultChan <- userResult{user: user, err: err}
		}(id)
	}

	// 收集结果
	for i := 0; i < len(ids); i++ {
		result := <-resultChan
		if result.err == nil && result.user != nil {
			results = append(results, *result.user)
		}
	}

	return results, nil
}

// BatchCreateUsers 批量创建用户（演示 make 和 Channel）
func (s *userService) BatchCreateUsers(requests []*models.UserCreateRequest) ([]models.UserResponse, []error) {
	// 使用 make 创建结果切片
	successUsers := make([]models.UserResponse, 0)
	errors := make([]error, 0)

	// 使用 make 创建 Channel
	type createResult struct {
		user  *models.UserResponse
		err   error
		index int
	}
	resultChan := make(chan createResult, len(requests))

	// 并发创建用户
	for i, req := range requests {
		go func(index int, request *models.UserCreateRequest) {
			user, err := s.CreateUser(request)
			resultChan <- createResult{user: user, err: err, index: index}
		}(i, req)
	}

	// 收集结果
	for i := 0; i < len(requests); i++ {
		result := <-resultChan
		if result.err != nil {
			errors = append(errors, fmt.Errorf("索引 %d: %v", result.index, result.err))
		} else if result.user != nil {
			successUsers = append(successUsers, *result.user)
		}
	}

	return successUsers, errors
}
