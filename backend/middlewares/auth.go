package middlewares

import (
	"gin-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 token

		authHeader := c.GetHeader("Authorization")
		token := ""

		if authHeader != "" {
			// 验证 token 格式
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 如果 Header 中没有，尝试从查询参数获取
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 检查 token 是否在黑名单中
		if utils.IsTokenBlacklisted(token) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌已失效",
			})
			c.Abort()
			return
		}

		// 解析并验证 JWT token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 单点登录检查：验证 token 是否是用户的当前有效 token
		if !utils.IsUserCurrentToken(claims.UserID, token) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "您的账号已在其他设备登录",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()

	}
}
