# 参数验证说明文档

## 概述

本项目使用 `go-playground/validator` 进行参数验证，提供了友好的中文错误提示。

## 验证规则

### 用户创建 (UserCreateRequest)

| 字段     | 类型   | 验证规则                          | 说明                                  |
| -------- | ------ | --------------------------------- | ------------------------------------- |
| username | string | required, min=3, max=20, alphanum | 必填，3-20 个字符，只能包含字母和数字 |
| email    | string | required, email, max=100          | 必填，必须是有效邮箱，最长 100 字符   |
| password | string | required, min=6, max=50           | 必填，6-50 个字符                     |
| nickname | string | omitempty, max=50                 | 可选，最长 50 字符                    |

### 用户更新 (UserUpdateRequest)

| 字段     | 类型   | 验证规则                  | 说明                                |
| -------- | ------ | ------------------------- | ----------------------------------- |
| email    | string | omitempty, email, max=100 | 可选，必须是有效邮箱，最长 100 字符 |
| nickname | string | omitempty, max=50         | 可选，最长 50 字符                  |
| avatar   | string | omitempty, url, max=500   | 可选，必须是有效 URL，最长 500 字符 |

## 支持的验证标签

### 基础验证

- `required`: 必填项
- `omitempty`: 可选项（为空时跳过其他验证）

### 字符串验证

- `min=n`: 最小长度
- `max=n`: 最大长度
- `len=n`: 固定长度
- `email`: 邮箱格式
- `url`: URL 格式
- `uri`: URI 格式
- `alphanum`: 只能包含字母和数字
- `alpha`: 只能包含字母
- `numeric`: 只能包含数字

### 数字验证

- `gt=n`: 大于 n
- `gte=n`: 大于等于 n
- `lt=n`: 小于 n
- `lte=n`: 小于等于 n

### 枚举验证

- `oneof=value1 value2`: 必须是指定值之一

## 错误响应格式

### 验证错误响应

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "username",
      "message": "Username 最小长度为 3"
    },
    {
      "field": "email",
      "message": "Email 必须是有效的邮箱地址"
    }
  ]
}
```

### 业务错误响应

```json
{
  "code": 409,
  "message": "用户名已存在"
}
```

## 测试示例

### 1. 缺少必填字段

**请求**:

```bash
POST /api/v1/users
Content-Type: application/json

{
  "email": "test@example.com"
}
```

**响应**:

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "username",
      "message": "Username 是必填项"
    },
    {
      "field": "password",
      "message": "Password 是必填项"
    }
  ]
}
```

### 2. 用户名太短

**请求**:

```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "ab",
  "email": "test@example.com",
  "password": "123456"
}
```

**响应**:

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "username",
      "message": "Username 最小长度为 3"
    }
  ]
}
```

### 3. 用户名包含特殊字符

**请求**:

```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "test@user",
  "email": "test@example.com",
  "password": "123456"
}
```

**响应**:

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "username",
      "message": "Username 只能包含字母和数字"
    }
  ]
}
```

### 4. 邮箱格式错误

**请求**:

```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "testuser",
  "email": "invalid-email",
  "password": "123456"
}
```

**响应**:

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "email",
      "message": "Email 必须是有效的邮箱地址"
    }
  ]
}
```

### 5. 密码太短

**请求**:

```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123"
}
```

**响应**:

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "password",
      "message": "Password 最小长度为 6"
    }
  ]
}
```

### 6. 头像 URL 格式错误

**请求**:

```bash
PUT /api/v1/users/1
Content-Type: application/json

{
  "avatar": "not-a-url"
}
```

**响应**:

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "avatar",
      "message": "Avatar 必须是有效的 URL"
    }
  ]
}
```

### 7. 正确的请求

**请求**:

```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123456",
  "nickname": "测试用户"
}
```

**响应**:

```json
{
  "code": 0,
  "message": "用户创建成功",
  "data": {
    "id": 3,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "",
    "created_at": "2025-12-25T16:40:00Z"
  }
}
```

## 自定义验证

如果需要添加自定义验证规则，可以在 `utils/validator.go` 中注册：

```go
func init() {
    validate = validator.New()

    // 注册自定义验证
    validate.RegisterValidation("custom_rule", func(fl validator.FieldLevel) bool {
        // 自定义验证逻辑
        return true
    })
}
```

然后在模型中使用：

```go
type MyModel struct {
    Field string `json:"field" binding:"custom_rule"`
}
```

## 最佳实践

1. **前端也要验证**: 虽然后端有验证，前端也应该进行基础验证以提升用户体验
2. **明确的错误信息**: 确保错误信息清晰，帮助用户快速定位问题
3. **合理的验证规则**: 不要过于严格，也不要过于宽松
4. **安全性考虑**: 对于敏感字段（如密码），不要在错误信息中暴露具体值
5. **国际化**: 如果需要支持多语言，可以扩展 `getErrorMessage` 函数

## 工具函数

### ValidateStruct

直接验证结构体：

```go
err := utils.ValidateStruct(&req)
if err != nil {
    errors := utils.FormatValidationErrors(err)
    // 处理错误
}
```

### FormatValidationErrors

格式化验证错误为友好的格式：

```go
errors := utils.FormatValidationErrors(err)
// 返回 []ValidationError
```

### GetValidationErrorMessage

获取第一条错误消息：

```go
message := utils.GetValidationErrorMessage(err)
// 返回字符串
```
