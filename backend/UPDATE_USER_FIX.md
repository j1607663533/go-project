# UpdateUser 接收 body 参数问题 - 已修复

## 🐛 问题描述

前端调用 `updateUser` API 时传递了 `role_id` 参数，但后端没有接收和处理该参数。

## 🔍 问题原因

`UserUpdateRequest` 结构体中缺少 `RoleID` 字段，导致：

1. 前端传递的 `role_id` 参数被忽略
2. 用户角色无法通过更新接口修改

## ✅ 解决方案

### 1. 更新 `UserUpdateRequest` 结构体

**文件**: `backend/models/user.go`

**修改前**:

```go
type UserUpdateRequest struct {
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty,url,max=500"`
}
```

**修改后**:

```go
type UserUpdateRequest struct {
	RoleID   uint   `json:"role_id" binding:"omitempty,min=1"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty,url,max=500"`
}
```

### 2. 更新 `UpdateUser` 服务方法

**文件**: `backend/services/user_service.go`

**添加的代码**:

```go
// 更新字段
if req.RoleID > 0 {
	user.RoleID = req.RoleID
}
if req.Nickname != "" {
	user.Nickname = req.Nickname
}
// ...
```

## 📝 修改说明

### 字段验证规则

```go
RoleID uint `json:"role_id" binding:"omitempty,min=1"`
```

- `json:"role_id"`: JSON 字段名为 `role_id`
- `binding:"omitempty,min=1"`:
  - `omitempty`: 字段可选
  - `min=1`: 如果提供，值必须 >= 1

### 更新逻辑

```go
if req.RoleID > 0 {
	user.RoleID = req.RoleID
}
```

- 只有当 `RoleID > 0` 时才更新
- 避免将角色设置为 0（无效值）

## 🎯 使用示例

### 前端调用

```javascript
import { updateUser } from "../api/auth";

// 更新用户角色
await updateUser(userId, {
  role_id: 1, // 设置为超级管理员
  nickname: "管理员",
  email: "admin@example.com",
});
```

### API 请求

```http
PUT /api/v1/users/1
Content-Type: application/json
Authorization: Bearer <token>

{
  "role_id": 1,
  "nickname": "管理员",
  "email": "admin@example.com"
}
```

### API 响应

```json
{
  "code": 0,
  "message": "用户更新成功",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "管理员",
    "role_id": 1,
    "role_name": "超级管理员",
    "created_at": "2026-01-06T10:00:00Z",
    "updated_at": "2026-01-06T16:00:00Z"
  }
}
```

## ✅ 验证步骤

1. **重启后端服务**

   ```bash
   cd backend
   go run main.go
   ```

2. **在前端测试**

   - 进入"系统管理 > 用户管理"
   - 点击"编辑"某个用户
   - 修改角色
   - 保存

3. **检查结果**
   - 用户角色应该成功更新
   - 刷新页面，角色显示正确

## 🔧 相关文件

- `backend/models/user.go` - 添加了 `RoleID` 字段
- `backend/services/user_service.go` - 添加了 `RoleID` 更新逻辑
- `backend/controllers/user_controller.go` - 已有正确的参数绑定
- `src/api/auth.js` - 前端 API 调用
- `src/pages/SystemUsers.jsx` - 用户管理页面

## 📋 支持的更新字段

现在 `UpdateUser` 接口支持更新以下字段：

- ✅ `role_id` - 用户角色
- ✅ `email` - 邮箱（会检查是否重复）
- ✅ `nickname` - 昵称
- ✅ `avatar` - 头像 URL

## 🎉 总结

问题已修复！现在可以通过用户管理页面正常更新用户角色了。

**修改内容**:

1. ✅ 在 `UserUpdateRequest` 中添加 `RoleID` 字段
2. ✅ 在 `UpdateUser` 服务方法中添加角色更新逻辑
3. ✅ 支持通过 API 更新用户角色

**下一步**:

- 重启后端服务
- 测试用户角色更新功能
