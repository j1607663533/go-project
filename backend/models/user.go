package models

import "time"

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password  string    `json:"-" gorm:"not null;size:255"` // - 表示不在 JSON 中序列化
	Nickname  string    `json:"nickname" gorm:"size:50"`
	Avatar    string    `json:"avatar" gorm:"size:500"`
	RoleID    uint      `json:"role_id" gorm:"default:2"` // 角色ID，默认为普通用户
	Role      *Role     `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanum" validate:"required,min=3,max=20,alphanum"`
	Email    string `json:"email" binding:"required,email,max=100" validate:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6,max=50" validate:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"omitempty,max=50" validate:"omitempty,max=50"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	RoleID   uint   `json:"role_id" binding:"omitempty,min=1"`
	Email    string `json:"email" binding:"omitempty,email,max=100" validate:"omitempty,email,max=100"`
	Nickname string `json:"nickname" binding:"omitempty,max=50" validate:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty,url,max=500" validate:"omitempty,url,max=500"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username  string `json:"username" binding:"required" validate:"required"`
	Password  string `json:"password" binding:"required" validate:"required"`
	CaptchaID string `json:"captcha_id" binding:"required" validate:"required"`
	Captcha   string `json:"captcha" binding:"required,len=6" validate:"required,len=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string             `json:"token"`
	User  UserResponse       `json:"user"`
	Menus []MenuTreeResponse `json:"menus"` // 用户可访问的菜单
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaID  string `json:"captcha_id"`
	CaptchaURL string `json:"captcha_url"`
}

// UserResponse 用户响应（不包含敏感信息）
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	RoleID    uint      `json:"role_id"`
	RoleName  string    `json:"role_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

//

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
	resp := UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
	}
	if u.Role != nil {
		resp.RoleName = u.Role.Name
	}
	return resp
}
