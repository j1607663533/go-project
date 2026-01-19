# 前端路由配置完成

## ✅ 路由已配置

路由配置已经在 `src/main.jsx` 中完成：

### 路由列表

1. **`/login`** - 登录页面

   - 未登录用户的入口
   - 已登录用户访问会自动跳转到首页

2. **`/`** - 首页（受保护）

   - 需要登录才能访问
   - 未登录用户会自动跳转到登录页面

3. **`/*`** - 其他路径
   - 自动重定向到首页

## 🚀 如何访问

### 1. 确保服务运行

**后端**：

```bash
cd backend
go run main.go
# 运行在 http://localhost:8080
```

**前端**：

```bash
cd d:\project-test\ReactFunTime
npm run dev
# 运行在 http://localhost:5173
```

### 2. 访问页面

- 登录页面：http://localhost:5173/login
- 首页：http://localhost:5173/

### 3. 测试流程

1. 打开浏览器访问 http://localhost:5173
2. 如果未登录，会自动跳转到 http://localhost:5173/login
3. 输入用户名和密码登录
4. 登录成功后自动跳转到首页

## 🔧 路由配置说明

### main.jsx 中的路由

```javascript
const router = createBrowserRouter([
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/",
    element: (
      <ProtectedRoute>
        <Home />
      </ProtectedRoute>
    ),
  },
  {
    path: "*",
    element: <Navigate to="/" replace />,
  },
]);
```

### ProtectedRoute 组件

```javascript
const ProtectedRoute = ({ children }) => {
  return isAuthenticated() ? children : <Navigate to="/login" replace />;
};
```

**作用**：

- 检查用户是否已登录
- 已登录：显示页面内容
- 未登录：重定向到登录页面

## 📝 测试账号

使用后端创建的测试账号：

- **用户名**：`bob`
- **密码**：`password123`
- **验证码**：查看图片输入（6 位数字）

## 🎨 页面特性

### 登录页面特性

- ✨ 渐变紫色背景
- 🎯 Material-UI 设计
- 🖼️ 图形验证码
- 🔄 点击刷新验证码
- 👁️ 密码显示/隐藏
- ⚠️ 错误提示
- ⏳ 加载状态

### 首页特性

- 👤 用户信息展示
- 📊 信息卡片布局
- 🚪 退出登录按钮
- 🛡️ 自动路由保护

## 🐛 故障排查

### 问题 1：页面空白

**可能原因**：

- 前端服务未启动
- 浏览器缓存问题

**解决方法**：

```bash
# 重启前端服务
Ctrl+C
npm run dev

# 清除浏览器缓存
Ctrl+Shift+R (硬刷新)
```

### 问题 2：登录后立即跳回登录页

**可能原因**：

- Token 未正确保存
- LocalStorage 被禁用

**解决方法**：

```javascript
// 在浏览器控制台检查
console.log(localStorage.getItem("token"));

// 清除并重试
localStorage.clear();
```

### 问题 3：验证码不显示

**可能原因**：

- 后端服务未启动
- CORS 问题

**解决方法**：

```bash
# 确保后端运行
cd backend
go run main.go

# 检查后端日志
```

### 问题 4：路由不工作

**可能原因**：

- React Router 配置错误
- 依赖未安装

**解决方法**：

```bash
# 重新安装依赖
npm install

# 重启服务
npm run dev
```

## 📦 依赖检查

确保以下依赖已安装：

```json
{
  "dependencies": {
    "react": "^19.2.3",
    "react-dom": "^19.2.3",
    "react-router-dom": "^7.11.0",
    "@mui/material": "^7.3.6",
    "@mui/icons-material": "^7.3.6",
    "@emotion/react": "^11.14.0",
    "@emotion/styled": "^11.14.1"
  }
}
```

如果缺少依赖，运行：

```bash
npm install
```

## 🎯 下一步

现在你可以：

1. ✅ 访问登录页面
2. ✅ 使用测试账号登录
3. ✅ 查看用户信息
4. ✅ 退出登录

所有功能都已就绪！🎉

---

**配置完成时间**: 2025-12-25  
**状态**: ✅ 就绪
