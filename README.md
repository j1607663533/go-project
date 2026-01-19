# React Admin - 后台管理系统

一个功能完整的后台管理系统，基于 React + Ant Design + Go + Gin 构建，实现了完整的 RBAC 权限管理。

## ✨ 特性

- 🎨 **现代化 UI**: 基于 Ant Design 的美观界面
- 🔐 **完整的权限系统**: RBAC 角色权限控制
- 🚀 **动态路由**: 根据用户权限动态加载菜单和路由
- 👥 **用户管理**: 用户的增删改查和角色分配
- 🎭 **角色管理**: 角色的增删改查和菜单权限分配
- 📋 **菜单管理**: 菜单的增删改查，支持树形结构
- 🛡️ **超级管理员**: 内置超级管理员，受保护不可删除
- 📦 **订单管理**: 完整的 CRUD 示例
- 🔑 **JWT 认证**: 基于 JWT 的用户认证
- 🎯 **验证码**: 图形验证码防护

## 🖥️ 技术栈

### 前端

- React 19
- Ant Design 5
- React Router v7
- Axios
- Vite

### 后端

- Go 1.25
- Gin
- GORM
- MySQL
- Redis
- JWT

## 📸 截图

### 登录页面

- 美观的登录界面
- 图形验证码
- 表单验证

### 主页面

- 响应式布局
- 动态侧边栏菜单
- 数据统计卡片
- 图表展示

### 系统管理

- 菜单管理（树形结构）
- 角色管理（权限分配）
- 用户管理（角色分配）

## 🚀 快速开始

### 环境要求

- Node.js 16+
- Go 1.25+
- MySQL 5.7+
- Redis（可选）

### 安装步骤

1. **克隆项目**

```bash
git clone <repository-url>
cd ReactFunTime
```

2. **配置数据库**

```sql
CREATE DATABASE reactfuntime CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

3. **配置后端**

```bash
cd backend
cp .env.example .env
# 编辑 .env 文件，配置数据库连接
```

4. **启动后端**

```bash
cd backend
go run main.go
```

5. **安装前端依赖**

```bash
npm install
```

6. **启动前端**

```bash
npm run dev
```

7. **访问系统**

- 前端: http://localhost:5173
- 后端: http://localhost:8080

详细步骤请查看 [快速启动指南](./QUICKSTART.md)

## 📚 文档

- [快速启动指南](./QUICKSTART.md) - 详细的安装和使用说明
- [动态路由文档](./backend/DYNAMIC_ROUTES.md) - 动态路由实现细节
- [系统页面指南](./SYSTEM_PAGES_GUIDE.md) - 系统管理页面使用说明
- [API 测试文档](./backend/API_TEST.md) - API 接口测试
- [功能总结](./DYNAMIC_ROUTES_SUMMARY.md) - 已实现功能总结

## 🎯 核心功能

### 1. 用户认证

- ✅ 用户注册（带验证码）
- ✅ 用户登录（JWT）
- ✅ 退出登录
- ✅ 单点登录（SSO）

### 2. 权限管理

- ✅ 基于角色的访问控制（RBAC）
- ✅ 动态菜单加载
- ✅ 菜单权限控制
- ✅ 超级管理员保护

### 3. 系统管理

- ✅ 用户管理（查看、编辑、分配角色）
- ✅ 角色管理（增删改查、分配菜单）
- ✅ 菜单管理（增删改查、树形结构）

### 4. 业务功能

- ✅ 订单管理（完整 CRUD）
- ✅ 数据分页
- ✅ 数据筛选

## 🔐 默认账号

系统启动后会自动创建：

### 角色

- **超级管理员** (ID=1, code=super_admin)
  - 拥有所有权限
  - 不可删除、不可修改
- **普通用户** (ID=2, code=user)
  - 只有首页和订单管理权限

### 创建超级管理员账号

注册账号后，在数据库中执行：

```sql
UPDATE users SET role_id = 1 WHERE username = 'your_username';
```

## 📁 项目结构

```
ReactFunTime/
├── backend/                 # 后端代码
│   ├── config/             # 配置
│   ├── controllers/        # 控制器
│   ├── models/             # 模型
│   ├── repositories/       # 数据访问
│   ├── routes/             # 路由
│   ├── services/           # 业务逻辑
│   ├── utils/              # 工具
│   └── middlewares/        # 中间件
├── src/                    # 前端代码
│   ├── api/                # API 接口
│   ├── components/         # 组件
│   ├── pages/              # 页面
│   └── utils/              # 工具
└── docs/                   # 文档
```

## 🛠️ 开发

### 后端开发

```bash
cd backend
go run main.go
```

### 前端开发

```bash
npm run dev
```

### 构建

```bash
# 前端构建
npm run build

# 后端构建
cd backend
go build -o app main.go
```

## 📝 API 文档

### 认证接口

- `POST /api/v1/register` - 用户注册
- `POST /api/v1/login` - 用户登录
- `POST /api/v1/logout` - 退出登录
- `GET /api/v1/captcha` - 获取验证码

### 用户接口

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户

### 角色接口

- `GET /api/v1/roles` - 获取角色列表
- `POST /api/v1/roles` - 创建角色
- `PUT /api/v1/roles/:id` - 更新角色
- `DELETE /api/v1/roles/:id` - 删除角色
- `POST /api/v1/roles/:id/menus` - 分配菜单

### 菜单接口

- `GET /api/v1/menus/tree` - 获取菜单树
- `GET /api/v1/menus/user` - 获取用户菜单
- `POST /api/v1/menus` - 创建菜单
- `PUT /api/v1/menus/:id` - 更新菜单
- `DELETE /api/v1/menus/:id` - 删除菜单

### 订单接口

- `GET /api/v1/orders` - 获取订单列表
- `POST /api/v1/orders` - 创建订单
- `PUT /api/v1/orders/:id` - 更新订单
- `DELETE /api/v1/orders/:id` - 删除订单

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🙏 致谢

- [React](https://react.dev/)
- [Ant Design](https://ant.design/)
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)

---

**注意**: 这是一个演示项目，生产环境使用前请进行安全加固和性能优化。
