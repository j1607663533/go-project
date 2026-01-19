package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-this-in-production") // 生产环境应该从环境变量读取

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(userID uint, username, email string) (string, error) {
	// 设置过期时间为 24 小时
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的 token")
}

// RefreshToken 刷新 token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 如果 token 在 30 分钟内过期，则刷新
	if time.Until(claims.ExpiresAt.Time) < 30*time.Minute {
		return GenerateToken(claims.UserID, claims.Username, claims.Email)
	}

	return tokenString, nil
}

// ClearToken 清除 token（将 token 加入黑名单）
func ClearToken(tokenString string) error {
	// 解析 token 获取过期时间
	claims, err := ParseToken(tokenString)
	if err != nil {
		return err
	}

	// 计算 token 剩余有效时间
	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime <= 0 {
		// token 已过期，无需加入黑名单
		return nil
	}

	// 将 token 加入黑名单，过期时间与 token 剩余有效时间一致
	blacklistKey := "blacklist:token:" + tokenString
	return CacheSet(blacklistKey, true, remainingTime)
}

// IsTokenBlacklisted 检查 token 是否在黑名单中
func IsTokenBlacklisted(tokenString string) bool {
	blacklistKey := "blacklist:token:" + tokenString
	var isBlacklisted bool
	err := CacheGet(blacklistKey, &isBlacklisted)
	return err == nil && isBlacklisted
}

// SetUserToken 设置用户的当前 token（单点登录）
// 如果用户已有 token，会将旧 token 加入黑名单
func SetUserToken(userID uint, newToken string) error {
	userTokenKey := fmt.Sprintf("user:token:%d", userID)

	// 获取用户的旧 token
	var oldToken string
	err := CacheGet(userTokenKey, &oldToken)
	// 打印日志
	fmt.Println("oldToken", oldToken)
	fmt.Println("newToken", newToken)
	fmt.Println("err", err)
	if err == nil && oldToken != "" && oldToken != newToken {
		// 将旧 token 加入黑名单
		ClearToken(oldToken)
	}

	// 存储新 token，过期时间 24 小时（与 token 有效期一致）
	return CacheSet(userTokenKey, newToken, 24*time.Hour)
}

// GetUserToken 获取用户的当前 token
func GetUserToken(userID uint) (string, error) {
	userTokenKey := fmt.Sprintf("user:token:%d", userID)
	var token string
	err := CacheGet(userTokenKey, &token)
	return token, err
}

// IsUserCurrentToken 检查 token 是否是用户的当前有效 token
func IsUserCurrentToken(userID uint, tokenString string) bool {
	currentToken, err := GetUserToken(userID)
	if err != nil {
		return false
	}
	return currentToken == tokenString
}
