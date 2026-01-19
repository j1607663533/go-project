# 测试后端 API 接口

## 1. 测试健康检查

GET http://localhost:8080/health

### 2. 获取验证码

GET http://localhost:8080/api/v1/captcha

### 3. 用户登录（需要先注册用户）

POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
"username": "test",
"password": "123456",
"captcha_id": "替换为实际的验证码 ID",
"captcha": "替换为实际的验证码"
}

### 4. 获取菜单树（需要登录后的 token）

GET http://localhost:8080/api/v1/menus/tree
Authorization: Bearer 你的 token

### 5. 获取当前用户菜单

GET http://localhost:8080/api/v1/menus/user
Authorization: Bearer 你的 token

### 6. 获取所有角色

GET http://localhost:8080/api/v1/roles
Authorization: Bearer 你的 token

### 7. 创建菜单

POST http://localhost:8080/api/v1/menus
Content-Type: application/json
Authorization: Bearer 你的 token

{
"parent_id": 0,
"name": "测试菜单",
"path": "/test",
"component": "Test",
"icon": "HomeOutlined",
"sort": 10,
"type": 1,
"hidden": false
}

### 8. 创建角色

POST http://localhost:8080/api/v1/roles
Content-Type: application/json
Authorization: Bearer 你的 token

{
"name": "测试角色",
"code": "test_role",
"description": "这是一个测试角色",
"menu_ids": [1, 2]
}

### 9. 为角色分配菜单

POST http://localhost:8080/api/v1/roles/2/menus
Content-Type: application/json
Authorization: Bearer 你的 token

{
"menu_ids": [1, 2, 3]
}

### 10. 获取所有用户

GET http://localhost:8080/api/v1/users
Authorization: Bearer 你的 token
