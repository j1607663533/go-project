# 图形验证码功能说明

## ✅ 功能已实现

项目已成功集成图形验证码功能，用于登录时的安全验证。

## 🎯 功能特点

- ✅ **Base64 图片返回** - 一次请求获取验证码 ID 和图片
- ✅ **自动刷新** - 支持通过参数刷新验证码
- ✅ **自动过期** - 验证码 10 分钟后自动过期
- ✅ **内存存储** - 使用内存存储，无需数据库
- ✅ **登录集成** - 登录时必须提供验证码

## 📝 API 接口

### 1. 获取验证码

**请求**：

```http
GET /api/v1/captcha
```

**响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "captcha_id": "abc123xyz",
    "captcha_image": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
  }
}
```

**说明**：

- `captcha_id`: 验证码唯一标识，登录时需要提供
- `captcha_image`: Base64 编码的 PNG 图片，可直接在前端显示

**前端使用示例**：

```html
<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..." />
```

### 2. 刷新验证码

**请求**：

```http
GET /api/v1/captcha?refresh=abc123xyz
```

**响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "captcha_id": "abc123xyz",
    "captcha_image": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
  }
}
```

**说明**：

- 使用相同的 `captcha_id`，但生成新的验证码图片
- 如果原验证码已过期，会自动生成新的 ID

### 3. 验证验证码（测试用）

**请求**：

```http
POST /api/v1/captcha/verify
Content-Type: application/json

{
  "captcha_id": "abc123xyz",
  "captcha": "123456"
}
```

**响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true
  }
}
```

### 4. 登录（需要验证码）

**请求**：

```http
POST /api/v1/login
Content-Type: application/json

{
  "username": "bob",
  "password": "password123",
  "captcha_id": "abc123xyz",
  "captcha": "123456"
}
```

**成功响应**：

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 2,
      "username": "bob",
      "email": "bob@example.com"
    }
  }
}
```

**验证码错误响应**：

```json
{
  "code": 401,
  "message": "验证码错误或已过期"
}
```

## 🔧 使用流程

### 完整的登录流程

```
1. 获取验证码
   GET /api/v1/captcha
   ↓ 返回 captcha_id 和 captcha_image

2. 用户查看验证码图片并输入

3. 提交登录请求
   POST /api/v1/login
   {
     username, password,
     captcha_id, captcha
   }
   ↓
4. 服务器验证验证码和密码
   ↓
5. 返回 JWT token
```

## 💻 前端集成示例

### React 示例

```jsx
import React, { useState } from "react";

function LoginForm() {
  const [captchaData, setCaptchaData] = useState(null);
  const [formData, setFormData] = useState({
    username: "",
    password: "",
    captcha: "",
  });

  // 获取验证码
  const getCaptcha = async () => {
    const response = await fetch("http://localhost:8080/api/v1/captcha");
    const data = await response.json();
    setCaptchaData(data.data);
  };

  // 刷新验证码
  const refreshCaptcha = async () => {
    if (captchaData) {
      const response = await fetch(
        `http://localhost:8080/api/v1/captcha?refresh=${captchaData.captcha_id}`
      );
      const data = await response.json();
      setCaptchaData(data.data);
    }
  };

  // 登录
  const handleLogin = async (e) => {
    e.preventDefault();

    const response = await fetch("http://localhost:8080/api/v1/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        ...formData,
        captcha_id: captchaData.captcha_id,
      }),
    });

    const result = await response.json();
    if (result.code === 0) {
      // 登录成功
      localStorage.setItem("token", result.data.token);
    } else {
      // 登录失败，刷新验证码
      refreshCaptcha();
    }
  };

  // 组件加载时获取验证码
  React.useEffect(() => {
    getCaptcha();
  }, []);

  return (
    <form onSubmit={handleLogin}>
      <input
        type="text"
        placeholder="用户名"
        value={formData.username}
        onChange={(e) => setFormData({ ...formData, username: e.target.value })}
      />

      <input
        type="password"
        placeholder="密码"
        value={formData.password}
        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
      />

      <div>
        {captchaData && (
          <img
            src={captchaData.captcha_image}
            alt="验证码"
            onClick={refreshCaptcha}
            style={{ cursor: "pointer" }}
          />
        )}
      </div>

      <input
        type="text"
        placeholder="验证码"
        value={formData.captcha}
        onChange={(e) => setFormData({ ...formData, captcha: e.target.value })}
        maxLength={6}
      />

      <button type="submit">登录</button>
    </form>
  );
}
```

### Vue 示例

```vue
<template>
  <form @submit.prevent="handleLogin">
    <input v-model="formData.username" placeholder="用户名" />
    <input v-model="formData.password" type="password" placeholder="密码" />

    <div v-if="captchaData">
      <img
        :src="captchaData.captcha_image"
        alt="验证码"
        @click="refreshCaptcha"
        style="cursor: pointer"
      />
    </div>

    <input v-model="formData.captcha" placeholder="验证码" maxlength="6" />

    <button type="submit">登录</button>
  </form>
</template>

<script>
export default {
  data() {
    return {
      captchaData: null,
      formData: {
        username: "",
        password: "",
        captcha: "",
      },
    };
  },

  mounted() {
    this.getCaptcha();
  },

  methods: {
    async getCaptcha() {
      const response = await fetch("http://localhost:8080/api/v1/captcha");
      const data = await response.json();
      this.captchaData = data.data;
    },

    async refreshCaptcha() {
      if (this.captchaData) {
        const response = await fetch(
          `http://localhost:8080/api/v1/captcha?refresh=${this.captchaData.captcha_id}`
        );
        const data = await response.json();
        this.captchaData = data.data;
      }
    },

    async handleLogin() {
      const response = await fetch("http://localhost:8080/api/v1/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          ...this.formData,
          captcha_id: this.captchaData.captcha_id,
        }),
      });

      const result = await response.json();
      if (result.code === 0) {
        localStorage.setItem("token", result.data.token);
      } else {
        this.refreshCaptcha();
      }
    },
  },
};
</script>
```

## ⚙️ 配置说明

### 验证码参数

在 `utils/captcha.go` 中可以配置：

```go
// 验证码长度（默认 6 位）
captcha.NewLen(6)

// 验证码过期时间（默认 10 分钟）
const CaptchaExpiration = 10 * time.Minute

// 验证码图片尺寸
captcha.StdWidth  // 240px
captcha.StdHeight // 80px
```

### 修改验证码长度

```go
// utils/captcha.go
func GenerateCaptcha() string {
    return captcha.NewLen(4) // 改为 4 位
}
```

### 修改过期时间

验证码库默认过期时间为 10 分钟，存储在内存中。

## 🔒 安全特性

1. **一次性使用** - 验证码验证后即失效
2. **自动过期** - 10 分钟后自动失效
3. **随机生成** - 每次生成的验证码都是随机的
4. **内存存储** - 不持久化，重启服务器后清空

## 📊 错误处理

| 场景           | 错误码 | 错误信息             |
| -------------- | ------ | -------------------- |
| 验证码错误     | 401    | 验证码错误或已过期   |
| 验证码过期     | 401    | 验证码错误或已过期   |
| 缺少验证码     | 400    | 请求参数错误         |
| 验证码长度错误 | 400    | Captcha 长度必须为 6 |

## 🎨 自定义验证码样式

如果需要自定义验证码样式，可以使用其他验证码库，例如：

- `github.com/mojocn/base64Captcha` - 支持更多样式
- `github.com/steambap/captcha` - 支持数学题验证码

## 🚀 性能优化

### 内存使用

验证码存储在内存中，每个验证码约占用几 KB 空间。默认配置下：

- 最多存储 1000 个验证码
- 自动清理过期验证码

### 并发处理

验证码库使用线程安全的存储，支持高并发访问。

## 🧪 测试

### 测试验证码生成

```bash
curl http://localhost:8080/api/v1/captcha
```

### 测试验证码刷新

```bash
curl "http://localhost:8080/api/v1/captcha?refresh=YOUR_CAPTCHA_ID"
```

### 测试验证码验证

```bash
curl -X POST http://localhost:8080/api/v1/captcha/verify \
  -H "Content-Type: application/json" \
  -d '{"captcha_id":"YOUR_ID","captcha":"123456"}'
```

### 测试登录

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username":"bob",
    "password":"password123",
    "captcha_id":"YOUR_ID",
    "captcha":"123456"
  }'
```

## 📚 相关文件

- `utils/captcha.go` - 验证码工具
- `controllers/captcha_controller.go` - 验证码控制器
- `models/user.go` - 登录请求模型（包含验证码字段）
- `services/user_service.go` - 登录业务逻辑（验证验证码）

## 🎯 下一步优化

- [ ] 添加验证码点击次数限制
- [ ] 支持语音验证码
- [ ] 支持滑块验证码
- [ ] 添加验证码难度配置
- [ ] 支持自定义验证码字符集

---

**实现时间**: 2025-12-25  
**验证码库**: github.com/dchest/captcha v1.1.0  
**状态**: ✅ 生产就绪
