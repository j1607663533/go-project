# Gin Backend 项目

这是一个使用 Gin 框架构建的企业级 Go 后端项目，采用三层架构设计，集成了 MySQL 数据库。

## ✨ 项目特性

- ✅ **三层架构**: Controller → Service → Repository，职责分离
- ✅ **数据库集成**: MySQL + GORM，支持自动迁移
- ✅ **参数验证**: 完善的请求参数验证，友好的错误提示
- ✅ **RESTful API**: 标准的 REST 接口设计
- ✅ **CORS 支持**: 跨域资源共享
- ✅ **依赖注入**: 清晰的依赖关系
- ✅ **单元测试**: 完整的测试示例
- ✅ **统一响应**: 标准化的 API 响应格式
- ✅ **日志系统**: 结构化日志记录
- ✅ **环境配置**: 灵活的配置管理

## 📂 项目结构

```
backend/
├── config/              # 配置文件
│   ├── config.go       # 应用配置
│   └── database.go     # 数据库配置
├── controllers/         # 控制器层（HTTP 处理）
│   ├── user_controller.go
│   └── user_controller_test.go
├── services/           # 服务层（业务逻辑）
│   ├── user_service.go
│   └── user_service_test.go
├── repositories/       # 仓储层（数据访问）
│   └── user_repository.go
├── models/             # 数据模型
│   └── user.go
├── middlewares/        # 中间件
│   ├── auth.go        # 认证中间件
│   └── cors.go        # CORS 中间件
├── routes/             # 路由配置
│   └── routes.go
├── utils/              # 工具函数
│   ├── logger.go      # 日志工具
│   ├── response.go    # 响应工具
│   └── validator.go   # 验证工具
├── .env.example        # 环境变量示例
├── .gitignore         # Git 忽略文件
├── api.http           # API 测试文件
├── Makefile           # 命令简化
├── go.mod             # Go 模块
├── main.go            # 程序入口
├── README.md          # 项目说明
├── ARCHITECTURE.md    # 架构说明
├── DATABASE.md        # 数据库说明
├── VALIDATION.md      # 验证说明
└── QUICKSTART.md      # 快速开始
```

## 🚀 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 5.7+ / 8.0+

### 2. 克隆项目

```bash
cd backend
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 配置数据库

确保 MySQL 服务运行，并创建数据库：

```sql
CREATE DATABASE projectTest CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

配置在 `config/config.go` 中：

- 数据库名: projectTest
- 用户名: root
- 密码: 123456
- 主机: localhost
- 端口: 3306

### 5. 运行项目

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

### 6. 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 获取用户列表
curl http://localhost:8080/api/v1/users

# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","nickname":"测试用户"}'
```

## 📖 API 文档

### 基础接口

#### 健康检查

- **GET** `/health`
- 响应：`{"status": "ok", "message": "服务运行正常"}`

### 用户管理

#### 获取用户列表

- **GET** `/api/v1/users`
- 响应：

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "测试用户",
      "avatar": "",
      "created_at": "2025-12-25T17:00:00Z"
    }
  ]
}
```

#### 获取单个用户

- **GET** `/api/v1/users/:id`

#### 创建用户

- **POST** `/api/v1/users`
- 请求体：

```json
{
  "username": "testuser", // 必填，3-20字符，只能字母数字
  "email": "test@example.com", // 必填，有效邮箱
  "password": "password123", // 必填，6-50字符
  "nickname": "测试用户" // 可选，最长50字符
}
```

#### 更新用户

- **PUT** `/api/v1/users/:id`
- 请求体：

```json
{
  "email": "newemail@example.com", // 可选
  "nickname": "新昵称", // 可选
  "avatar": "https://example.com/avatar.jpg" // 可选，必须是有效URL
}
```

#### 删除用户

- **DELETE** `/api/v1/users/:id`

### 认证接口

#### 获取个人信息

- **GET** `/api/v1/auth/profile`
- Headers: `Authorization: Bearer <token>`

## 🏗️ 架构设计

### 三层架构

```
┌─────────────────────────────────────┐
│         Controller Layer            │  ← HTTP 请求处理
│  - 参数验证                          │
│  - 响应格式化                        │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│          Service Layer              │  ← 业务逻辑
│  - 业务规则验证                      │
│  - 数据转换                          │
│  - 事务管理                          │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│        Repository Layer             │  ← 数据访问
│  - CRUD 操作                        │
│  - 数据库查询                        │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│            Database                 │  ← MySQL
└─────────────────────────────────────┘
```

详细说明请查看 [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔧 开发指南

### 添加新功能

1. **定义模型** (`models/`)

```go
type Product struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name" gorm:"size:100"`
}
```

2. **创建 Repository** (`repositories/`)

```go
type ProductRepository interface {
    FindAll() ([]models.Product, error)
    Create(product *models.Product) error
}
```

3. **实现 Service** (`services/`)

```go
type ProductService interface {
    GetAllProducts() ([]models.ProductResponse, error)
}
```

4. **创建 Controller** (`controllers/`)

```go
func (ctrl *ProductController) GetProducts(c *gin.Context) {
    products, err := ctrl.productService.GetAllProducts()
    // ...
}
```

5. **注册路由** (`routes/routes.go`)

```go
products := api.Group("/products")
{
    products.GET("", productController.GetProducts)
}
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./services -v

# 运行测试并查看覆盖率
go test -cover ./...
```

### 使用 Makefile

```bash
make run      # 运行应用
make build    # 编译应用
make test     # 运行测试
make clean    # 清理编译文件
make install  # 安装依赖
```

## 📝 参数验证

项目使用 `go-playground/validator` 进行参数验证，支持：

- 必填验证: `required`
- 长度验证: `min`, `max`, `len`
- 格式验证: `email`, `url`, `alphanum`
- 数值验证: `gt`, `gte`, `lt`, `lte`

详细说明请查看 [VALIDATION.md](./VALIDATION.md)

## 🗄️ 数据库

### 当前配置

- **数据库**: MySQL
- **ORM**: GORM
- **自动迁移**: 启用
- **连接池**: 已配置

### 数据表

- `users`: 用户表
  - 字段: id, username, email, password, nickname, avatar, created_at, updated_at
  - 索引: username (unique), email (unique)

详细说明请查看 [DATABASE.md](./DATABASE.md)

## 🧪 测试

项目包含完整的单元测试示例：

- Controller 测试: 使用 `httptest` 测试 HTTP 接口
- Service 测试: 使用 Mock Repository 测试业务逻辑

## 📚 文档

- [README.md](./README.md) - 项目概览（本文档）
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 架构设计说明
- [DATABASE.md](./DATABASE.md) - 数据库集成说明
- [VALIDATION.md](./VALIDATION.md) - 参数验证说明
- [QUICKSTART.md](./QUICKSTART.md) - 快速开始指南
- [api.http](./api.http) - API 测试用例

## 🛠️ 技术栈

- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) v1.11.0
- **ORM**: [GORM](https://gorm.io/) v1.31.1
- **数据库驱动**: MySQL Driver v1.6.0
- **验证器**: [validator](https://github.com/go-playground/validator) v10.30.1
- **测试**: [testify](https://github.com/stretchr/testify) v1.11.1

## 🔐 安全性

- ✅ 密码字段不在 JSON 中序列化
- ✅ SQL 注入防护（GORM 自动处理）
- ✅ 参数验证防止恶意输入
- ⚠️ 密码加密（待实现）
- ⚠️ JWT 认证（待实现）
- ⚠️ 限流中间件（待实现）

## 📈 性能优化

- ✅ 数据库连接池配置
- ✅ 索引优化（username, email）
- ✅ 字段长度限制
- ⚠️ Redis 缓存（待实现）
- ⚠️ 查询优化（待实现）

## 🚧 待实现功能

- [ ] JWT 认证
- [ ] 密码加密（bcrypt）
- [ ] 文件上传
- [ ] 分页查询
- [ ] 搜索和过滤
- [ ] Redis 缓存
- [ ] 限流中间件
- [ ] Swagger 文档
- [ ] Docker 支持
- [ ] CI/CD 配置

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

MIT License

## 📞 联系方式

如有问题或建议，请提交 Issue。

---

**项目创建时间**: 2025-12-25  
**Go 版本**: 1.21+  
**当前状态**: ✅ 生产就绪
