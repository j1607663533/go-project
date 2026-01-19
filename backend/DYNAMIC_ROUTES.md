# 动态路由和角色权限系统实现文档

## 概述

本文档描述了基于角色的动态路由系统的实现，包括超级管理员角色的创建和保护机制。

## 后端实现

### 1. 数据模型

#### Role（角色模型）

- `id`: 主键
- `name`: 角色名称
- `code`: 角色编码（如 super_admin, user）
- `description`: 角色描述
- `is_super`: 是否为超级管理员（超级管理员不能删除和修改）
- `status`: 状态（1-启用，0-禁用）
- `menus`: 角色拥有的菜单（多对多关系）

#### Menu（菜单模型）

- `id`: 主键
- `parent_id`: 父菜单 ID（0 表示顶级菜单）
- `name`: 菜单名称
- `path`: 路由路径
- `component`: 组件路径
- `icon`: 图标
- `sort`: 排序
- `type`: 类型（1-菜单，2-按钮）
- `status`: 状态（1-启用，0-禁用）
- `hidden`: 是否隐藏

#### User（用户模型更新）

- 添加了 `role_id` 字段，关联到角色表
- 添加了 `role` 关联字段

### 2. 核心功能

#### 角色管理

- **创建角色**: 支持创建新角色并分配菜单权限
- **更新角色**: 更新角色信息和菜单权限
- **删除角色**: 删除角色（超级管理员角色受保护，不能删除）
- **查询角色**: 获取角色列表和详情

#### 菜单管理

- **创建菜单**: 支持创建顶级菜单和子菜单
- **更新菜单**: 更新菜单信息
- **删除菜单**: 删除菜单
- **菜单树**: 构建层级菜单树结构
- **用户菜单**: 根据用户角色获取可访问的菜单

#### 超级管理员保护

- 超级管理员角色 `is_super = true`
- 超级管理员角色不能被删除
- 超级管理员角色不能被修改
- 超级管理员角色的菜单权限不能被修改

### 3. API 接口

#### 角色接口

```
GET    /api/v1/roles           - 获取所有角色
GET    /api/v1/roles/:id       - 获取角色详情
POST   /api/v1/roles           - 创建角色
PUT    /api/v1/roles/:id       - 更新角色
DELETE /api/v1/roles/:id       - 删除角色
POST   /api/v1/roles/:id/menus - 为角色分配菜单
```

#### 菜单接口

```
GET    /api/v1/menus/tree      - 获取菜单树
GET    /api/v1/menus/user      - 获取当前用户菜单
POST   /api/v1/menus           - 创建菜单
PUT    /api/v1/menus/:id       - 更新菜单
DELETE /api/v1/menus/:id       - 删除菜单
```

#### 登录接口更新

登录接口现在返回用户菜单信息：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "...",
    "user": {
      "id": 1,
      "username": "admin",
      "role_id": 1,
      "role_name": "超级管理员"
    },
    "menus": [
      {
        "id": 1,
        "name": "首页",
        "path": "/",
        "component": "AntdHome",
        "icon": "HomeOutlined",
        "children": []
      }
    ]
  }
}
```

### 4. 默认数据

系统启动时会自动初始化以下数据：

#### 默认角色

1. **超级管理员** (super_admin)

   - 拥有所有菜单权限
   - 不能删除和修改
   - ID: 1

2. **普通用户** (user)
   - 只有首页和订单管理权限
   - ID: 2

#### 默认菜单

1. 首页 (/)
2. 订单管理 (/orders)
3. 系统管理 (/system)
   - 用户管理 (/system/users)
   - 角色管理 (/system/roles)
   - 菜单管理 (/system/menus)

## 前端实现（待完成）

### 1. 动态路由加载

- 登录后从后端获取用户菜单
- 根据菜单动态生成路由
- 更新侧边栏菜单显示

### 2. 路由守卫

- 检查用户是否有权限访问当前路由
- 未授权自动跳转到 403 页面

### 3. 菜单渲染

- 根据后端返回的菜单数据渲染侧边栏
- 支持多级菜单
- 支持图标显示

## 使用说明

### 创建新角色

```bash
POST /api/v1/roles
{
  "name": "编辑",
  "code": "editor",
  "description": "内容编辑角色",
  "menu_ids": [1, 2]
}
```

### 为角色分配菜单

```bash
POST /api/v1/roles/2/menus
{
  "menu_ids": [1, 2, 3]
}
```

### 创建新菜单

```bash
POST /api/v1/menus
{
  "parent_id": 0,
  "name": "产品管理",
  "path": "/products",
  "component": "Products",
  "icon": "ShoppingOutlined",
  "sort": 4,
  "type": 1
}
```

## 注意事项

1. 超级管理员角色（ID=1）受到特殊保护，不能删除和修改
2. 新注册用户默认分配普通用户角色（ID=2）
3. 所有菜单和角色接口都需要认证
4. 菜单的 component 字段应该对应前端的组件名称
5. 菜单的 icon 字段应该使用 Ant Design 的图标名称

## 数据库表结构

### roles 表

```sql
CREATE TABLE roles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(200),
    is_super BOOLEAN DEFAULT FALSE,
    status INT DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### menus 表

```sql
CREATE TABLE menus (
    id INT PRIMARY KEY AUTO_INCREMENT,
    parent_id INT DEFAULT 0,
    name VARCHAR(50) NOT NULL,
    path VARCHAR(200) NOT NULL,
    component VARCHAR(200),
    icon VARCHAR(50),
    sort INT DEFAULT 0,
    type INT DEFAULT 1,
    status INT DEFAULT 1,
    hidden BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### role_menus 表（关联表）

```sql
CREATE TABLE role_menus (
    role_id INT,
    menu_id INT,
    PRIMARY KEY (role_id, menu_id)
);
```

### users 表更新

```sql
ALTER TABLE users ADD COLUMN role_id INT DEFAULT 2;
```
