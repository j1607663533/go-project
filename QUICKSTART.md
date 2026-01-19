# 快速启动指南

## 系统概述

这是一个基于 React + Ant Design + Go + Gin 的后台管理系统，实现了完整的 RBAC（基于角色的访问控制）权限管理。

## 主要功能

- ✅ 用户认证（登录/注册/退出）
- ✅ 动态菜单（根据角色权限显示）
- ✅ 角色管理（创建、编辑、删除、分配菜单）
- ✅ 菜单管理（创建、编辑、删除、树形结构）
- ✅ 用户管理（查看、编辑、分配角色）
- ✅ 订单管理（CRUD 操作）
- ✅ 超级管理员保护机制

## 快速启动

### 1. 环境准备

**后端要求**:

- Go 1.25+
- MySQL 5.7+
- Redis（可选，用于缓存）

**前端要求**:

- Node.js 16+
- npm 或 yarn

### 2. 数据库配置

1. 创建数据库：

```sql
CREATE DATABASE reactfuntime CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 配置后端环境变量（`backend/.env`）：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=你的密码
DB_NAME=reactfuntime

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

### 3. 启动后端

```bash
cd backend
go run main.go
```

后端会自动：

- 连接数据库
- 迁移表结构
- 初始化角色和菜单数据
- 启动在 http://localhost:8080

**初始化数据**:

- 超级管理员角色（ID=1）
- 普通用户角色（ID=2）
- 6 个默认菜单（首页、订单管理、系统管理及其子菜单）

### 4. 启动前端

```bash
npm install
npm run dev
```

前端启动在 http://localhost:5173

### 5. 注册和登录

1. 访问 http://localhost:5173
2. 点击"注册"创建新账号
3. 注册成功后会自动分配"普通用户"角色
4. 登录后可以看到首页和订单管理菜单

### 6. 创建超级管理员账号

**方法 1: 直接修改数据库**

```sql
-- 将某个用户设置为超级管理员
UPDATE users SET role_id = 1 WHERE username = 'admin';
```

**方法 2: 使用现有超级管理员账号**

- 如果已有超级管理员账号，登录后在用户管理中为其他用户分配角色

## 使用流程

### 作为超级管理员

1. **登录系统**

   - 使用超级管理员账号登录
   - 可以看到所有菜单（包括系统管理）

2. **管理菜单**

   - 进入"系统管理 > 菜单管理"
   - 创建、编辑、删除菜单
   - 设置菜单图标、排序、隐藏等

3. **管理角色**

   - 进入"系统管理 > 角色管理"
   - 创建新角色（如：编辑、审核员等）
   - 为角色分配菜单权限

4. **管理用户**
   - 进入"系统管理 > 用户管理"
   - 查看所有用户
   - 为用户分配角色

### 作为普通用户

1. **注册账号**

   - 填写用户名、邮箱、密码
   - 完成验证码验证

2. **登录系统**

   - 默认只能看到"首页"和"订单管理"

3. **使用功能**
   - 查看订单列表
   - 创建、编辑、删除订单

## 权限说明

### 角色类型

1. **超级管理员** (super_admin)

   - 拥有所有权限
   - 不能被删除
   - 不能被修改
   - 菜单权限不能被修改

2. **普通用户** (user)

   - 默认角色
   - 只有首页和订单管理权限
   - 可以被修改

3. **自定义角色**
   - 可以自由创建
   - 可以分配任意菜单权限
   - 可以编辑和删除

### 菜单类型

1. **菜单** (type=1)

   - 显示在侧边栏
   - 可以点击跳转

2. **按钮** (type=2)
   - 不显示在侧边栏
   - 用于按钮级权限控制

## 目录结构

```
ReactFunTime/
├── backend/                 # 后端代码
│   ├── config/             # 配置文件
│   ├── controllers/        # 控制器
│   ├── models/             # 数据模型
│   ├── repositories/       # 数据访问层
│   ├── routes/             # 路由配置
│   ├── services/           # 业务逻辑层
│   ├── utils/              # 工具函数
│   ├── middlewares/        # 中间件
│   ├── main.go             # 入口文件
│   ├── init_data.sql       # 数据初始化脚本
│   └── API_TEST.md         # API 测试文档
├── src/                    # 前端代码
│   ├── api/                # API 接口
│   ├── components/         # 组件
│   ├── pages/              # 页面
│   │   ├── Login.jsx       # 登录页
│   │   ├── Register.jsx    # 注册页
│   │   ├── AntdHome.jsx    # 首页
│   │   ├── Orders.jsx      # 订单管理
│   │   ├── SystemUsers.jsx # 用户管理
│   │   ├── SystemRoles.jsx # 角色管理
│   │   └── SystemMenus.jsx # 菜单管理
│   ├── utils/              # 工具函数
│   └── main.jsx            # 入口文件
└── docs/                   # 文档
    ├── DYNAMIC_ROUTES.md           # 动态路由文档
    ├── DYNAMIC_ROUTES_SUMMARY.md   # 功能总结
    └── SYSTEM_PAGES_GUIDE.md       # 系统页面指南
```

## 常见问题

### 1. 登录后看不到菜单？

**原因**: 用户没有分配角色或角色没有菜单权限

**解决**:

1. 检查用户的 `role_id` 是否正确
2. 检查角色是否有菜单权限
3. 使用超级管理员账号为用户分配角色

### 2. 菜单数据为空？

**原因**: 数据库初始化失败

**解决**:

1. 检查后端启动日志
2. 手动执行 `backend/init_data.sql` 脚本
3. 重启后端服务

### 3. 无法访问系统管理页面？

**原因**: 当前用户角色没有系统管理权限

**解决**:

1. 使用超级管理员账号登录
2. 或者在数据库中将用户的 `role_id` 改为 1

### 4. Redis 连接失败？

**原因**: Redis 服务未启动

**解决**:

1. 启动 Redis 服务
2. 或者忽略（系统会使用内存存储作为降级方案）

## 技术栈

### 后端

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT
- **验证**: validator

### 前端

- **框架**: React 19
- **UI**: Ant Design
- **路由**: React Router v7
- **HTTP**: Axios
- **构建**: Vite

## 下一步

1. **添加更多功能模块**

   - 创建新的菜单
   - 开发对应的前端页面
   - 分配给相应角色

2. **完善权限控制**

   - 添加按钮级权限
   - 添加数据权限
   - 添加操作日志

3. **优化用户体验**

   - 添加加载动画
   - 优化错误提示
   - 添加操作确认

4. **部署上线**
   - 配置生产环境
   - 优化性能
   - 配置 HTTPS

## 相关文档

- [动态路由实现文档](./backend/DYNAMIC_ROUTES.md)
- [系统页面使用指南](./SYSTEM_PAGES_GUIDE.md)
- [API 测试文档](./backend/API_TEST.md)
- [功能总结](./DYNAMIC_ROUTES_SUMMARY.md)

## 联系支持

如有问题，请查看相关文档或提交 Issue。
