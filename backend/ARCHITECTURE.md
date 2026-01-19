# 项目架构说明

## 三层架构设计

本项目采用经典的三层架构模式，将业务逻辑和数据访问分离，提高代码的可维护性和可测试性。

```
┌─────────────────────────────────────────────────────────┐
│                     HTTP Request                         │
└─────────────────────┬───────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────┐
│                 Controller Layer                         │
│  职责：处理 HTTP 请求和响应                                │
│  - 参数验证和绑定                                          │
│  - HTTP 状态码处理                                        │
│  - 响应格式化                                             │
└─────────────────────┬───────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────┐
│                  Service Layer                           │
│  职责：业务逻辑处理                                        │
│  - 业务规则验证                                           │
│  - 数据转换和组装                                         │
│  - 事务管理                                              │
│  - 调用多个 Repository                                   │
└─────────────────────┬───────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────┐
│                Repository Layer                          │
│  职责：数据访问                                           │
│  - CRUD 操作                                            │
│  - 数据库查询                                            │
│  - 数据持久化                                            │
└─────────────────────┬───────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────┐
│                     Database                             │
└─────────────────────────────────────────────────────────┘
```

## 目录结构

```
backend/
├── controllers/         # 控制器层
│   ├── user_controller.go
│   └── user_controller_test.go
├── services/           # 服务层（业务逻辑）
│   └── user_service.go
├── repositories/       # 仓储层（数据访问）
│   └── user_repository.go
├── models/             # 数据模型
│   └── user.go
├── middlewares/        # 中间件
│   ├── auth.go
│   └── cors.go
├── routes/             # 路由配置
│   └── routes.go
├── config/             # 配置
│   ├── config.go
│   └── database.go
└── utils/              # 工具函数
    ├── logger.go
    └── response.go
```

## 各层职责详解

### 1. Controller Layer（控制器层）

**文件**: `controllers/user_controller.go`

**职责**:

- 接收和解析 HTTP 请求
- 参数验证和绑定
- 调用 Service 层处理业务
- 格式化响应数据
- 设置 HTTP 状态码

**示例**:

```go
func (ctrl *UserController) GetUser(c *gin.Context) {
    // 1. 解析参数
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    // 2. 调用 Service
    user, err := ctrl.userService.GetUserByID(uint(id))

    // 3. 返回响应
    c.JSON(http.StatusOK, gin.H{
        "code": 0,
        "data": user,
    })
}
```

**特点**:

- ✅ 不包含业务逻辑
- ✅ 不直接访问数据库
- ✅ 只负责 HTTP 层面的处理

### 2. Service Layer（服务层）

**文件**: `services/user_service.go`

**职责**:

- 实现核心业务逻辑
- 业务规则验证（如：用户名唯一性检查）
- 数据转换和组装
- 调用一个或多个 Repository
- 事务管理（如果需要）

**示例**:

```go
func (s *userService) CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error) {
    // 1. 业务验证：检查用户名是否存在
    if existingUser, _ := s.userRepo.FindByUsername(req.Username); existingUser != nil {
        return nil, errors.New("用户名已存在")
    }

    // 2. 业务验证：检查邮箱是否存在
    if existingUser, _ := s.userRepo.FindByEmail(req.Email); existingUser != nil {
        return nil, errors.New("邮箱已存在")
    }

    // 3. 业务逻辑：密码加密（示例）
    // hashedPassword := hashPassword(req.Password)

    // 4. 创建实体
    user := &models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    // 5. 调用 Repository 保存
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    // 6. 返回响应
    response := user.ToResponse()
    return &response, nil
}
```

**特点**:

- ✅ 包含所有业务逻辑
- ✅ 不关心 HTTP 细节
- ✅ 通过 Repository 访问数据
- ✅ 可以调用多个 Repository

### 3. Repository Layer（仓储层）

**文件**: `repositories/user_repository.go`

**职责**:

- 封装数据访问逻辑
- CRUD 操作
- 数据库查询
- 数据持久化

**示例**:

```go
func (r *userRepository) FindByID(id uint) (*models.User, error) {
    // 直接的数据库操作
    for _, user := range users {
        if user.ID == id {
            return &user, nil
        }
    }
    return nil, errors.New("用户不存在")
}

func (r *userRepository) Create(user *models.User) error {
    // 直接的数据库操作
    user.ID = nextID
    nextID++
    users = append(users, *user)
    return nil
}
```

**特点**:

- ✅ 只负责数据访问
- ✅ 不包含业务逻辑
- ✅ 使用接口定义，便于测试和替换实现
- ✅ 可以轻松切换数据源（内存、MySQL、PostgreSQL 等）

## 依赖注入流程

在 `routes/routes.go` 中实现依赖注入：

```go
// 1. 创建 Repository
userRepo := repositories.NewUserRepository()

// 2. 创建 Service，注入 Repository
userService := services.NewUserService(userRepo)

// 3. 创建 Controller，注入 Service
userController := controllers.NewUserController(userService)

// 4. 注册路由
users.GET("", userController.GetUsers)
```

## 数据流向

### 创建用户的完整流程：

```
1. HTTP POST /api/v1/users
   ↓
2. Controller.CreateUser()
   - 解析 JSON 请求体
   - 验证参数格式
   ↓
3. Service.CreateUser()
   - 检查用户名是否存在 (调用 Repository.FindByUsername)
   - 检查邮箱是否存在 (调用 Repository.FindByEmail)
   - 密码加密等业务逻辑
   - 创建用户实体
   ↓
4. Repository.Create()
   - 保存到数据库
   ↓
5. Service 返回结果给 Controller
   ↓
6. Controller 格式化响应
   ↓
7. HTTP Response (JSON)
```

## 优势

### 1. **关注点分离**

- 每一层只关注自己的职责
- 代码更清晰，易于理解

### 2. **可测试性**

- 每一层都可以独立测试
- 使用接口便于 Mock

```go
// 可以轻松创建 Mock Repository 进行测试
type mockUserRepository struct {}
func (m *mockUserRepository) FindByID(id uint) (*models.User, error) {
    return &models.User{ID: id}, nil
}
```

### 3. **可维护性**

- 修改数据库实现不影响业务逻辑
- 修改业务逻辑不影响 HTTP 处理
- 易于添加新功能

### 4. **可扩展性**

- 轻松切换数据源（内存 → MySQL → PostgreSQL）
- 可以添加缓存层
- 可以添加消息队列等

### 5. **复用性**

- Service 可以被多个 Controller 使用
- Repository 可以被多个 Service 使用

## 实际应用示例

### 切换到真实数据库

只需修改 `repositories/user_repository.go`：

```go
type userRepository struct {
    db *gorm.DB  // 添加数据库连接
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    if err := r.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

**Service 和 Controller 完全不需要修改！**

### 添加缓存

在 Service 层添加缓存逻辑：

```go
func (s *userService) GetUserByID(id uint) (*models.UserResponse, error) {
    // 1. 先查缓存
    if cachedUser := s.cache.Get(id); cachedUser != nil {
        return cachedUser, nil
    }

    // 2. 缓存未命中，查数据库
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    s.cache.Set(id, user)

    return user.ToResponse(), nil
}
```

## 最佳实践

1. **Controller 应该很薄**：只做参数解析和响应格式化
2. **Service 包含业务逻辑**：所有业务规则都在这里
3. **Repository 只做数据访问**：不要在这里写业务逻辑
4. **使用接口**：便于测试和替换实现
5. **依赖注入**：通过构造函数注入依赖
6. **错误处理**：每一层都应该妥善处理错误

## 总结

这种架构模式的核心思想是：

> **每一层只做一件事，并且做好这件事**

- Controller：处理 HTTP
- Service：处理业务
- Repository：处理数据

这样的设计让代码更加清晰、可测试、可维护！
