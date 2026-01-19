# 用户注册接口说明

## ✅ 功能已实现

后端已成功添加用户注册接口，包含验证码验证功能。

## 📝 API 接口

### 用户注册

**请求**：

```http
POST /api/v1/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "testuser@example.com",
  "password": "password123",
  "nickname": "测试用户",
  "captcha_id": "abc123",
  "captcha": "123456"
}
```

**成功响应** (201 Created)：

```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "id": 3,
    "username": "testuser",
    "email": "testuser@example.com",
    "nickname": "测试用户",
    "avatar": "",
    "created_at": "2025-12-25T18:00:00Z"
  }
}
```

**失败响应**：

1. **参数错误** (400 Bad Request)：

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "Username",
      "message": "Username 为必填字段"
    }
  ]
}
```

2. **验证码错误** (400 Bad Request)：

```json
{
  "code": 400,
  "message": "验证码错误或已过期"
}
```

3. **用户名已存在** (409 Conflict)：

```json
{
  "code": 409,
  "message": "用户名已存在"
}
```

4. **邮箱已存在** (409 Conflict)：

```json
{
  "code": 409,
  "message": "邮箱已存在"
}
```

## 🔧 参数说明

| 参数       | 类型   | 必填 | 说明      | 验证规则                        |
| ---------- | ------ | ---- | --------- | ------------------------------- |
| username   | string | 是   | 用户名    | 3-20 个字符，只能包含字母和数字 |
| email      | string | 是   | 邮箱      | 有效的邮箱格式，最长 100 个字符 |
| password   | string | 是   | 密码      | 6-50 个字符                     |
| nickname   | string | 否   | 昵称      | 最长 50 个字符                  |
| captcha_id | string | 是   | 验证码 ID | 从获取验证码接口返回            |
| captcha    | string | 是   | 验证码    | 6 位数字                        |

## 🚀 使用流程

### 完整的注册流程

```
1. 获取验证码
   GET /api/v1/captcha
   ↓ 返回 captcha_id 和 captcha_image

2. 用户查看验证码图片并输入

3. 提交注册请求
   POST /api/v1/register
   {
     username, email, password, nickname,
     captcha_id, captcha
   }
   ↓
4. 服务器验证
   - 验证码是否正确
   - 参数格式是否正确
   - 用户名是否已存在
   - 邮箱是否已存在
   ↓
5. 创建用户
   - 密码自动加密（bcrypt）
   - 保存到数据库
   ↓
6. 返回用户信息
```

## 💻 前端集成示例

### JavaScript/Fetch

```javascript
// 1. 获取验证码
const captchaResponse = await fetch("http://localhost:8080/api/v1/captcha");
const captchaData = await captchaResponse.json();
const { captcha_id, captcha_image } = captchaData.data;

// 2. 显示验证码图片
document.getElementById("captcha-img").src = captcha_image;

// 3. 用户填写表单并提交
const registerData = {
  username: "testuser",
  email: "testuser@example.com",
  password: "password123",
  nickname: "测试用户",
  captcha_id: captcha_id,
  captcha: "123456", // 用户输入的验证码
};

const response = await fetch("http://localhost:8080/api/v1/register", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify(registerData),
});

const result = await response.json();

if (result.code === 0) {
  console.log("注册成功！", result.data);
  // 可以直接跳转到登录页面
  window.location.href = "/login";
} else {
  console.error("注册失败：", result.message);
  // 刷新验证码
  refreshCaptcha();
}
```

### React 示例

```jsx
import { useState } from "react";

function Register() {
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
    nickname: "",
    captcha: "",
  });
  const [captchaData, setCaptchaData] = useState(null);
  const [error, setError] = useState("");

  // 获取验证码
  const loadCaptcha = async () => {
    const response = await fetch("http://localhost:8080/api/v1/captcha");
    const data = await response.json();
    setCaptchaData(data.data);
  };

  // 注册
  const handleRegister = async (e) => {
    e.preventDefault();

    try {
      const response = await fetch("http://localhost:8080/api/v1/register", {
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
        alert("注册成功！");
        // 跳转到登录页面
        window.location.href = "/login";
      } else {
        setError(result.message);
        loadCaptcha(); // 刷新验证码
      }
    } catch (err) {
      setError("注册失败：" + err.message);
    }
  };

  return (
    <form onSubmit={handleRegister}>
      <input
        type="text"
        placeholder="用户名"
        value={formData.username}
        onChange={(e) => setFormData({ ...formData, username: e.target.value })}
        required
      />

      <input
        type="email"
        placeholder="邮箱"
        value={formData.email}
        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
        required
      />

      <input
        type="password"
        placeholder="密码"
        value={formData.password}
        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
        required
      />

      <input
        type="text"
        placeholder="昵称（可选）"
        value={formData.nickname}
        onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
      />

      {captchaData && (
        <img
          src={captchaData.captcha_image}
          alt="验证码"
          onClick={loadCaptcha}
        />
      )}

      <input
        type="text"
        placeholder="验证码"
        value={formData.captcha}
        onChange={(e) => setFormData({ ...formData, captcha: e.target.value })}
        maxLength={6}
        required
      />

      {error && <div className="error">{error}</div>}

      <button type="submit">注册</button>
    </form>
  );
}
```

## 🔒 安全特性

1. **验证码验证**

   - 防止机器人注册
   - 验证码一次性使用
   - 10 分钟自动过期

2. **密码加密**

   - 使用 bcrypt 加密
   - 不可逆加密
   - 每次加密结果不同

3. **唯一性检查**

   - 用户名唯一
   - 邮箱唯一

4. **参数验证**
   - 用户名：3-20 个字符，只能字母数字
   - 邮箱：有效格式
   - 密码：6-50 个字符

## 🧪 测试

### 使用 cURL 测试

```bash
# 1. 获取验证码
curl http://localhost:8080/api/v1/captcha

# 2. 注册用户（替换 captcha_id 和 captcha）
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "testuser@example.com",
    "password": "password123",
    "nickname": "测试用户",
    "captcha_id": "YOUR_CAPTCHA_ID",
    "captcha": "123456"
  }'
```

### 使用 api.http 测试

打开 `backend/api.http` 文件，找到注册测试用例：

1. 先执行"获取验证码"请求
2. 复制返回的 `captcha_id`
3. 查看验证码图片，输入验证码
4. 执行"用户注册"请求

## 📊 与登录接口的区别

| 特性       | 注册接口           | 登录接口        |
| ---------- | ------------------ | --------------- |
| 路径       | `/api/v1/register` | `/api/v1/login` |
| 需要验证码 | ✅ 是              | ✅ 是           |
| 需要邮箱   | ✅ 是              | ❌ 否           |
| 需要昵称   | ⭕ 可选            | ❌ 否           |
| 返回 Token | ❌ 否              | ✅ 是           |
| 密码加密   | ✅ 自动            | ✅ 验证         |

## 🎯 下一步

注册成功后，用户可以：

1. 使用注册的用户名和密码登录
2. 登录后获取 JWT Token
3. 使用 Token 访问受保护的接口

## 📚 相关文档

- 登录接口：`AUTH.md`
- 验证码功能：`CAPTCHA.md`
- 参数验证：`VALIDATION.md`

---

**实现时间**: 2025-12-25  
**状态**: ✅ 完成
