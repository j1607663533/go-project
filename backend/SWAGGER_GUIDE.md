# Swagger 接口文档集成与使用指南

本文档详细介绍了如何在 Gin 项目中引入、配置并使用 Swagger 生成自动化接口文档。

---

## 1. 安装依赖

首先，需要安装 `swag` 命令行工具以及 Gin 的 Swagger 中间件。

```bash
# 安装 swag 命令行工具（用于生成 docs 文件）
go install github.com/swaggo/swag/cmd/swag@latest

# 在项目中引入 Gin-Swagger 依赖
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
```

---

## 2. 第一步：配置全局 API 信息

在 `main.go` 的 `main` 函数上方添加全局注释。这些注释定义了文档的标题、版本、基础路径以及安全认证（JWT）的方案。

```go
// @title Gin Backend API
// @version 1.0
// @description This is a sample Gin backend service.
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() { ... }
```

- **@securityDefinitions.apikey**: 定义认证方案名称为 `Bearer`。
- **@name Authorization**: 指定 Token 存放在 HTTP Header 的 `Authorization` 字段。

---

## 3. 第二步：注册 Swagger 路由

在路由初始化文件（如 `routes/routes.go`）中导入生成的 `docs` 包，并挂载 Swagger 的访问路径。

```go
import (
    _ "gin-backend/docs" // 必须导入生成的 docs 包，否则会报错
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // 注册 Swagger 访问路径
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // ... 其他路由设置
    return r
}
```

---

## 4. 第三步：为 Controller 添加接口注释

在具体的 Controller 方法上方添加注释，描述接口的行为。

### 示例：受保护的接口（需要 Token）

```go
// @Summary 创建新用户
// @Description 只有持有有效 Token 的管理员可以手动创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body models.UserCreateRequest true "用户信息"
// @Success 201 {object} map[string]interface{}
// @Router /users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) { ... }
```

- **@Summary**: 接口简要说明。
- **@Tags**: 分组标签，UI 会根据标签对接口进行分类。
- **@Security Bearer**: 标记该接口需要 `Bearer` 认证（UI 上会显示一把锁）。
- **@Success**: 定义成功返回的数据结构。

---

## 5. 第四步：生成与更新文档

每当你修改了代码中的 `@` 注释后，都需要重新运行 `swag init` 来更新 `docs` 目录下的文件。

```bash
# 在项目根目录下执行
# 如果已安装 swag 到 PATH
swag init

# 或者使用 go run 方式（推荐，无需手动配置环境变量）
go run github.com/swaggo/swag/cmd/swag init
```

---

## 6. 如何使用

1.  **访问地址**: 启动服务后，访问 `http://localhost:8080/swagger/index.html`。
2.  **认证测试**:
    - 点击右上方的绿色 **Authorize** 按钮。
    - 在弹出框中输入 `Bearer YOUR_TOKEN`（注意 `Bearer` 和 Token 之间有一个空格）。
    - 点击 **Authorize** 后，所有带锁图标的接口在测试时都会自动带上请求头。
3.  **调试接口**: 展开具体接口，点击 **Try it out**，输入参数后点击 **Execute** 即可看到真实请求结果。

---

## 7. 常见问题 (FAQ)

- **访问 404?**: 确认是否在 `routes.go` 中注册了路由，且端口号与 `main.go` 中的 `@host` 一致。
- **文档不更新?**: 每次修改注释后必须执行 `swag init`。
- **Token 无效?**: 检查设置中 `@name` 是否为 `Authorization`，且输入时是否带了 `Bearer ` 前缀。
