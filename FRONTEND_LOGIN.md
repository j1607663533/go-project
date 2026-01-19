# 前端登录页面使用说明

## ✅ 已完成功能

### 1. 登录页面 (`/login`)

- ✅ 漂亮的渐变背景
- ✅ Material-UI 设计
- ✅ 用户名和密码输入
- ✅ 图形验证码显示
- ✅ 验证码刷新功能
- ✅ 密码显示/隐藏切换
- ✅ 错误提示
- ✅ 加载状态
- ✅ 自动跳转（已登录用户）

### 2. 首页 (`/`)

- ✅ 用户信息展示
- ✅ 退出登录功能
- ✅ 受保护路由（未登录自动跳转）

### 3. API 集成

- ✅ 获取验证码
- ✅ 刷新验证码
- ✅ 用户登录
- ✅ Token 存储
- ✅ 用户信息存储

## 🚀 使用方法

### 1. 启动后端服务

```bash
cd backend
go run main.go
```

后端将运行在 `http://localhost:8080`

### 2. 启动前端服务

```bash
cd d:\project-test\ReactFunTime
npm run dev
```

前端将运行在 `http://localhost:5173`

### 3. 访问登录页面

打开浏览器访问：`http://localhost:5173/login`

### 4. 登录测试

使用之前创建的测试账号：

- 用户名：`bob`
- 密码：`password123`
- 验证码：查看图片输入

## 📁 新增文件

```
src/
├── api/
│   └── auth.js          # 认证相关 API
├── pages/
│   ├── Login.jsx        # 登录页面
│   └── Home.jsx         # 首页
└── App.jsx              # 路由配置（已更新）
```

## 🎨 页面特性

### 登录页面

**设计亮点**：

- 渐变紫色背景
- 毛玻璃效果卡片
- 平滑的动画过渡
- 响应式设计

**功能特性**：

- 实时表单验证
- 验证码自动刷新
- 错误提示
- 加载状态反馈

**交互细节**：

- 点击验证码图片刷新
- 密码显示/隐藏切换
- Enter 键提交表单
- 自动聚焦用户名输入框

### 首页

**功能**：

- 显示用户信息
- 退出登录
- 自动检查登录状态

## 🔧 API 函数说明

### `getCaptcha(refreshId)`

获取或刷新验证码

```javascript
// 获取新验证码
const captcha = await getCaptcha();

// 刷新验证码
const newCaptcha = await getCaptcha(captcha.captcha_id);
```

**返回值**：

```javascript
{
  captcha_id: "abc123",
  captcha_image: "data:image/png;base64,..."
}
```

### `login(username, password, captchaId, captcha)`

用户登录

```javascript
const result = await login("bob", "password123", "abc123", "123456");
```

**返回值**：

```javascript
{
  token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  user: {
    id: 2,
    username: "bob",
    email: "bob@example.com",
    nickname: "Bob",
    created_at: "2025-12-25T17:18:00Z"
  }
}
```

**副作用**：

- 自动保存 token 到 localStorage
- 自动保存用户信息到 localStorage

### `logout()`

退出登录

```javascript
logout();
```

**副作用**：

- 清除 localStorage 中的 token
- 清除 localStorage 中的用户信息

### `isAuthenticated()`

检查是否已登录

```javascript
if (isAuthenticated()) {
  // 已登录
} else {
  // 未登录
}
```

### `getCurrentUser()`

获取当前用户信息

```javascript
const user = getCurrentUser();
console.log(user.username);
```

## 🛡️ 路由保护

### ProtectedRoute 组件

自动检查登录状态，未登录用户会被重定向到登录页面：

```jsx
<Route
  path="/"
  element={
    <ProtectedRoute>
      <Home />
    </ProtectedRoute>
  }
/>
```

## 🎯 工作流程

### 登录流程

```
1. 用户访问 /login
   ↓
2. 自动获取验证码
   ↓
3. 用户输入用户名、密码、验证码
   ↓
4. 点击登录按钮
   ↓
5. 调用后端 API
   ↓
6. 验证成功
   ↓
7. 保存 token 和用户信息
   ↓
8. 跳转到首页
```

### 访问受保护页面流程

```
1. 用户访问 /
   ↓
2. ProtectedRoute 检查登录状态
   ↓
3. 已登录？
   ├─ 是 → 显示首页
   └─ 否 → 重定向到 /login
```

## 🔍 调试技巧

### 查看 Token

```javascript
// 在浏览器控制台
console.log(localStorage.getItem("token"));
```

### 查看用户信息

```javascript
// 在浏览器控制台
console.log(JSON.parse(localStorage.getItem("user")));
```

### 清除登录状态

```javascript
// 在浏览器控制台
localStorage.clear();
```

## 🎨 自定义主题

在 `App.jsx` 中修改主题：

```javascript
const theme = createTheme({
  palette: {
    primary: {
      main: "#667eea", // 主色调
    },
    secondary: {
      main: "#764ba2", // 次色调
    },
  },
});
```

## 📱 响应式设计

页面已适配：

- ✅ 桌面端（1920px+）
- ✅ 笔记本（1366px+）
- ✅ 平板（768px+）
- ✅ 手机（375px+）

## 🚧 待扩展功能

- [ ] 注册页面
- [ ] 忘记密码
- [ ] 记住我功能
- [ ] 第三方登录（GitHub, Google）
- [ ] 双因素认证
- [ ] 用户资料编辑
- [ ] 头像上传

## 🐛 常见问题

### 1. 验证码不显示

**原因**：后端服务未启动或 CORS 问题

**解决**：

- 确保后端运行在 `http://localhost:8080`
- 检查后端 CORS 配置

### 2. 登录后立即退出

**原因**：Token 未正确保存

**解决**：

- 检查浏览器控制台错误
- 清除 localStorage 重试

### 3. 验证码错误

**原因**：验证码已过期或输入错误

**解决**：

- 点击验证码图片刷新
- 仔细输入验证码（6 位数字）

## 📚 相关文档

- Material-UI: https://mui.com/
- React Router: https://reactrouter.com/
- Fetch API: https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API

---

**创建时间**: 2025-12-25  
**技术栈**: React 19 + Material-UI 7 + React Router 7  
**状态**: ✅ 完成
