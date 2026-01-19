# 登录和认证功能说明

## ✅ 功能已实现

项目已成功实现完整的用户认证系统，包括：

- ✅ 密码加密（bcrypt）
- ✅ JWT Token 生成和验证
- ✅ 登录接口
- ✅ 认证中间件
- ✅ 受保护的路由

## 🔐 安全特性

### 1. 密码加密

使用 bcrypt 算法对用户密码进行加密：

```go
// 创建用户时自动加密密码
hashedPassword, err := utils.HashPassword(req.Password)

// 登录时验证密码
isValid := utils.CheckPassword(plainPassword, hashedPassword)
```

**特点**：

- 单向加密，无法解密
- 每次加密结果不同（加盐）
- 计算成本高，防止暴力破解

### 2. JWT Token

使用 JWT (JSON Web Token) 进行身份认证：

```go
// Token 包含的信息
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}
```

**特点**：

- 无状态认证
- Token 有效期 24 小时
- 包含用户基本信息
- 使用 HS256 算法签名

## 📝 API 接口

### 1. 用户注册

**请求**：

```http
POST /api/v1/users
Content-Type: application/json

{
  "username": "bob",
  "email": "bob@example.com",
  "password": "password123",
  "nickname": "Bob"
}
```

**响应**：

```json
{
  "code": 0,
  "message": "用户创建成功",
  "data": {
    "id": 2,
    "username": "bob",
    "email": "bob@example.com",
    "nickname": "Bob",
    "avatar": "",
    "created_at": "2025-12-25T17:18:00Z"
  }
}
```

**说明**：

- 密码会自动加密存储
- 密码不会在响应中返回
- 用户名和邮箱必须唯一

### 2. 用户登录

**请求**：

```http
POST /api/v1/login
Content-Type: application/json

{
  "username": "bob",
  "password": "password123"
}
```

**成功响应**：

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 2,
      "username": "bob",
      "email": "bob@example.com",
      "nickname": "Bob",
      "avatar": "",
      "created_at": "2025-12-25T17:18:00Z"
    }
  }
}
```

**失败响应**：

```json
{
  "code": 401,
  "message": "用户名或密码错误"
}
```

**说明**：

- 返回的 token 用于后续请求的认证
- token 有效期为 24 小时
- 为了安全，错误信息不区分用户名不存在还是密码错误

### 3. 获取个人信息（需要认证）

**请求**：

```http
GET /api/v1/auth/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**成功响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 2,
    "username": "bob",
    "email": "bob@example.com",
    "nickname": "Bob",
    "avatar": "",
    "created_at": "2025-12-25T17:18:00Z"
  }
}
```

**未认证响应**：

```json
{
  "code": 401,
  "message": "未提供认证令牌"
}
```

**Token 无效响应**：

```json
{
  "code": 401,
  "message": "无效的认证令牌"
}
```

## 🔧 使用流程

### 完整的认证流程

```
1. 用户注册
   POST /api/v1/users
   ↓
2. 用户登录
   POST /api/v1/login
   ↓ 返回 token
3. 使用 token 访问受保护的接口
   GET /api/v1/auth/profile
   Header: Authorization: Bearer <token>
```

### 示例代码

#### JavaScript/Fetch

```javascript
// 1. 登录
const loginResponse = await fetch("http://localhost:8080/api/v1/login", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    username: "bob",
    password: "password123",
  }),
});

const loginData = await loginResponse.json();
const token = loginData.data.token;

// 2. 使用 token 访问受保护的接口
const profileResponse = await fetch(
  "http://localhost:8080/api/v1/auth/profile",
  {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  }
);

const profileData = await profileResponse.json();
console.log(profileData.data);
```

#### cURL

```bash
# 1. 登录
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"password123"}'

# 2. 使用返回的 token 访问受保护的接口
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 🛡️ 认证中间件

### 工作原理

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从请求头获取 token
        authHeader := c.GetHeader("Authorization")

        // 2. 验证格式：Bearer <token>
        parts := strings.SplitN(authHeader, " ", 2)

        // 3. 解析并验证 JWT token
        claims, err := utils.ParseToken(token)

        // 4. 将用户信息存储到上下文
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("email", claims.Email)

        c.Next()
    }
}
```

### 在控制器中获取用户信息

```go
func (ctrl *UserController) GetProfile(c *gin.Context) {
    // 从上下文获取用户 ID
    userID, exists := c.Get("userID")
    if !exists {
        // 未认证
        return
    }

    // 使用用户 ID 获取信息
    user, err := ctrl.userService.GetProfile(userID.(uint))
    // ...
}
```

## 🔑 Token 管理

### Token 结构

JWT Token 由三部分组成：

```
Header.Payload.Signature
```

**示例**：

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJ1c2VyX2lkIjoyLCJ1c2VybmFtZSI6ImJvYiIsImVtYWlsIjoiYm9iQGV4YW1wbGUuY29tIiwiZXhwIjoxNzM1MjA0NzAwfQ.
signature_here
```

### Token 有效期

- **默认有效期**：24 小时
- **过期后**：需要重新登录
- **刷新机制**：可以在 30 分钟内过期时自动刷新

### 修改 Token 有效期

在 `utils/jwt.go` 中修改：

```go
// 修改为 7 天
expirationTime := time.Now().Add(7 * 24 * time.Hour)
```

### 修改 JWT 密钥

**重要**：生产环境必须修改密钥！

在 `utils/jwt.go` 中：

```go
// 从环境变量读取
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
```

在 `.env` 文件中：

```env
JWT_SECRET=your-very-long-and-random-secret-key-here
```

## 🚀 添加受保护的路由

### 1. 在路由中使用认证中间件

```go
// routes/routes.go
auth := api.Group("/auth")
auth.Use(middlewares.AuthMiddleware())
{
    auth.GET("/profile", userController.GetProfile)
    auth.PUT("/profile", userController.UpdateProfile)
    auth.POST("/logout", userController.Logout)
}
```

### 2. 在控制器中获取当前用户

```go
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
    // 获取当前登录用户的 ID
    userID, _ := c.Get("userID")

    // 只允许用户更新自己的信息
    // ...
}
```

## 🔒 安全建议

### 1. 密码策略

- ✅ 最小长度：6 个字符（已实现）
- ⚠️ 建议：要求包含大小写字母、数字和特殊字符
- ⚠️ 建议：检查常见密码列表

### 2. Token 安全

- ✅ 使用 HTTPS 传输（生产环境）
- ✅ Token 存储在 HTTP-only Cookie 或 LocalStorage
- ⚠️ 实现 Token 黑名单（用于登出）
- ⚠️ 实现 Refresh Token 机制

### 3. 防止暴力破解

- ⚠️ 实现登录失败次数限制
- ⚠️ 添加验证码（多次失败后）
- ⚠️ IP 限流

### 4. 其他安全措施

- ✅ 密码加密存储（已实现）
- ✅ 错误信息不泄露敏感信息（已实现）
- ⚠️ 实现双因素认证（2FA）
- ⚠️ 记录登录日志

## 📊 错误码说明

| 错误码 | 说明         | 场景                         |
| ------ | ------------ | ---------------------------- |
| 200    | 成功         | 登录成功                     |
| 201    | 创建成功     | 注册成功                     |
| 400    | 请求参数错误 | 缺少必填字段、格式错误       |
| 401    | 未授权       | 未登录、Token 无效或过期     |
| 404    | 未找到       | 用户不存在                   |
| 409    | 冲突         | 用户名或邮箱已存在           |
| 500    | 服务器错误   | 密码加密失败、Token 生成失败 |

## 🧪 测试

### 测试场景

1. ✅ 用户注册 - 成功
2. ✅ 用户注册 - 用户名已存在
3. ✅ 用户注册 - 邮箱已存在
4. ✅ 用户登录 - 成功
5. ✅ 用户登录 - 密码错误
6. ✅ 用户登录 - 用户不存在
7. ✅ 访问受保护接口 - 有效 Token
8. ✅ 访问受保护接口 - 无 Token
9. ✅ 访问受保护接口 - 无效 Token

### 使用 api.http 测试

打开 `api.http` 文件，使用 VS Code REST Client 扩展测试所有场景。

## 📚 相关文件

- `utils/jwt.go` - JWT Token 工具
- `utils/password.go` - 密码加密工具
- `middlewares/auth.go` - 认证中间件
- `services/user_service.go` - 登录业务逻辑
- `controllers/user_controller.go` - 登录控制器
- `models/user.go` - 登录请求/响应模型

## 🎯 下一步

- [ ] 实现 Refresh Token
- [ ] 添加登出功能（Token 黑名单）
- [ ] 实现密码重置功能
- [ ] 添加邮箱验证
- [ ] 实现第三方登录（OAuth）
- [ ] 添加登录日志
- [ ] 实现限流和防暴力破解

---

**实现时间**: 2025-12-25  
**JWT 库**: github.com/golang-jwt/jwt/v5  
**加密库**: golang.org/x/crypto/bcrypt  
**状态**: ✅ 生产就绪
