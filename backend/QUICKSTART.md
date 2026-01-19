# Gin 后端项目快速开始指南

## 📦 项目已创建完成！

你的 Gin 后端项目已经成功创建，包含以下功能：

### ✅ 已实现的功能

1. **基础架构**

   - ✅ Gin 框架集成
   - ✅ 项目结构规范（MVC 模式）
   - ✅ 环境变量配置
   - ✅ 统一响应格式

2. **中间件**

   - ✅ CORS 跨域支持
   - ✅ 认证中间件（Bearer Token）

3. **API 接口**

   - ✅ 健康检查接口
   - ✅ 用户 CRUD 接口
   - ✅ RESTful API 设计

4. **开发工具**

   - ✅ Makefile 命令简化
   - ✅ API 测试文件（api.http）
   - ✅ 单元测试示例
   - ✅ 日志工具
   - ✅ Git 配置

5. **文档**
   - ✅ README 文档
   - ✅ API 文档
   - ✅ 代码注释

## 🚀 当前状态

服务器正在运行：

- **地址**: http://localhost:8080
- **状态**: ✅ 运行中

## 📝 快速测试

### 1. 测试健康检查

```bash
curl http://localhost:8080/health
```

### 2. 获取用户列表

```bash
curl http://localhost:8080/api/v1/users
```

### 3. 创建新用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "password123",
    "nickname": "新用户"
  }'
```

## 📂 项目结构

```
backend/
├── config/              # 配置文件
│   ├── config.go       # 应用配置
│   └── database.go     # 数据库配置（待使用）
├── controllers/         # 控制器
│   ├── user_controller.go
│   └── user_controller_test.go
├── middlewares/         # 中间件
│   ├── auth.go         # 认证中间件
│   └── cors.go         # CORS 中间件
├── models/              # 数据模型
│   └── user.go         # 用户模型
├── routes/              # 路由配置
│   └── routes.go       # 路由定义
├── utils/               # 工具函数
│   ├── logger.go       # 日志工具
│   └── response.go     # 统一响应
├── .env.example         # 环境变量示例
├── .gitignore          # Git 忽略文件
├── api.http            # API 测试文件
├── Makefile            # 命令简化
├── README.md           # 项目文档
├── go.mod              # Go 模块
└── main.go             # 程序入口
```

## 🔧 下一步建议

### 1. 集成数据库

```bash
# 安装 GORM 和 MySQL 驱动
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

然后在 `main.go` 中添加：

```go
import "gin-backend/config"

func main() {
    config.LoadConfig()

    // 初始化数据库
    if err := config.InitDB(); err != nil {
        log.Fatal(err)
    }
    defer config.CloseDB()

    // ... 其他代码
}
```

### 2. 实现 JWT 认证

```bash
go get -u github.com/golang-jwt/jwt/v5
```

### 3. 添加 Swagger 文档

```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 4. 添加热重载（开发模式）

```bash
go install github.com/cosmtrek/air@latest
air init
air
```

### 5. 添加验证器增强

```bash
go get -u github.com/go-playground/validator/v10
```

## 🎯 可用的 API 端点

| 方法   | 路径                 | 描述         | 认证 |
| ------ | -------------------- | ------------ | ---- |
| GET    | /health              | 健康检查     | ❌   |
| GET    | /api/v1/users        | 获取用户列表 | ❌   |
| GET    | /api/v1/users/:id    | 获取单个用户 | ❌   |
| POST   | /api/v1/users        | 创建用户     | ❌   |
| PUT    | /api/v1/users/:id    | 更新用户     | ❌   |
| DELETE | /api/v1/users/:id    | 删除用户     | ❌   |
| GET    | /api/v1/auth/profile | 获取个人信息 | ✅   |

## 💡 使用技巧

1. **使用 api.http 文件测试**

   - 安装 VS Code 的 REST Client 扩展
   - 打开 `api.http` 文件
   - 点击 "Send Request" 测试 API

2. **查看日志**

   - 服务器日志会显示所有请求信息
   - 使用 `utils.LogInfo()` 等函数记录自定义日志

3. **环境配置**
   - 复制 `.env.example` 为 `.env`
   - 修改配置以适应你的环境

## 🐛 常见问题

**Q: 如何修改端口？**
A: 在 `.env` 文件中设置 `APP_PORT=你的端口`

**Q: 如何启用生产模式？**
A: 在 `.env` 文件中设置 `APP_MODE=release`

**Q: 数据在哪里存储？**
A: 当前使用内存存储（重启会丢失），建议集成数据库

## 📚 学习资源

- [Gin 官方文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Go 标准库](https://pkg.go.dev/std)

---

**项目创建时间**: 2025-12-25
**Go 版本**: 请运行 `go version` 查看
**Gin 版本**: v1.10.0

祝你开发愉快！🎉
