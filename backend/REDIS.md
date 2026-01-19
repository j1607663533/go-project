# Redis 缓存集成说明

## ✅ 功能已实现

项目已成功集成 Redis 缓存，用于提升性能和支持分布式部署。

## 🎯 使用场景

### 1. 验证码存储 ⭐⭐⭐⭐⭐

**优势**：

- 支持分布式部署
- 自动过期管理
- 一次性使用验证

**实现**：

```go
// 生成验证码并存储到 Redis
id, answer, err := utils.GenerateCaptchaRedis()

// 验证验证码（验证后自动删除）
isValid := utils.VerifyCaptchaRedis(id, answer)
```

**Redis 键格式**：

- 键：`captcha:{id}`
- 值：验证码答案（6 位数字）
- 过期时间：10 分钟

### 2. 用户信息缓存 ⭐⭐⭐⭐

**优势**：

- 减少数据库查询
- 提升响应速度
- 降低数据库负载

**实现**：

```go
// 自动使用缓存
user, err := userService.GetUserByID(id)
// 首次查询：数据库 → 缓存
// 后续查询：缓存 → 直接返回
```

**Redis 键格式**：

- 键：`user:{id}`
- 值：JSON 格式的用户信息
- 过期时间：30 分钟

### 3. 登录失败次数限制（待实现）

**用途**：

- 防止暴力破解
- IP 限流

### 4. JWT Token 黑名单（待实现）

**用途**：

- 实现登出功能
- Token 撤销

## 📦 安装和配置

### 1. 安装 Redis

**Windows**：

```bash
# 使用 Chocolatey
choco install redis-64

# 或下载 MSI 安装包
# https://github.com/microsoftarchive/redis/releases
```

**Linux**：

```bash
sudo apt-get install redis-server
```

**macOS**：

```bash
brew install redis
```

**Docker**：

```bash
docker run -d -p 6379:6379 --name redis redis:latest
```

### 2. 启动 Redis

```bash
# Windows
redis-server

# Linux/macOS
redis-server

# Docker
docker start redis
```

### 3. 配置环境变量

在 `.env` 文件中配置（如果没有 Redis，会自动降级到内存存储）：

```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

## 🔧 缓存工具函数

### 基础操作

```go
// 设置缓存
utils.CacheSet("key", value, 10*time.Minute)

// 获取缓存
var result MyStruct
err := utils.CacheGet("key", &result)

// 设置字符串缓存
utils.CacheSetString("key", "value", 10*time.Minute)

// 获取字符串缓存
value, err := utils.CacheGetString("key")

// 删除缓存
utils.CacheDel("key1", "key2")

// 检查是否存在
exists, err := utils.CacheExists("key")
```

### 高级操作

```go
// 自增
count, err := utils.CacheIncr("counter")

// 指定步长自增
count, err := utils.CacheIncrBy("counter", 5)

// 设置过期时间
utils.CacheExpire("key", 10*time.Minute)

// 获取剩余过期时间
ttl, err := utils.CacheTTL("key")

// 获取匹配的键
keys, err := utils.CacheKeys("user:*")
```

### 哈希操作

```go
// 设置哈希字段
utils.CacheHSet("user:1", "name", "Alice")

// 获取哈希字段
name, err := utils.CacheHGet("user:1", "name")

// 获取所有哈希字段
fields, err := utils.CacheHGetAll("user:1")

// 删除哈希字段
utils.CacheHDel("user:1", "name", "email")
```

## 🚀 性能优化

### 缓存命中率

```go
// 查看缓存统计
INFO stats
```

### 内存使用

```go
// 查看内存使用
INFO memory
```

### 连接池配置

在 `config/redis.go` 中配置：

```go
redis.NewClient(&redis.Options{
    PoolSize:     10,  // 连接池大小
    MinIdleConns: 5,   // 最小空闲连接数
    DialTimeout:  5 * time.Second,
    ReadTimeout:  3 * time.Second,
    WriteTimeout: 3 * time.Second,
})
```

## 📊 缓存策略

### 1. 缓存更新策略

**Cache-Aside（旁路缓存）** - 当前使用

```
读取：
1. 先查缓存
2. 缓存未命中 → 查数据库
3. 将结果写入缓存

更新：
1. 更新数据库
2. 删除缓存
```

**示例**：

```go
// 更新用户时清除缓存
func (s *userService) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
    // 更新数据库
    user, err := s.userRepo.Update(user)

    // 删除缓存
    cacheKey := fmt.Sprintf("user:%d", id)
    utils.CacheDel(cacheKey)

    return &response, nil
}
```

### 2. 缓存过期策略

| 数据类型     | 过期时间       | 说明                |
| ------------ | -------------- | ------------------- |
| 验证码       | 10 分钟        | 安全考虑            |
| 用户信息     | 30 分钟        | 平衡性能和实时性    |
| Token 黑名单 | Token 过期时间 | 与 JWT 过期时间一致 |
| 登录失败次数 | 15 分钟        | 防止长期锁定        |

### 3. 缓存预热

```go
// 系统启动时预加载热点数据
func WarmUpCache() {
    // 加载活跃用户
    activeUsers := getActiveUsers()
    for _, user := range activeUsers {
        cacheKey := fmt.Sprintf("user:%d", user.ID)
        utils.CacheSet(cacheKey, user, 30*time.Minute)
    }
}
```

## 🔒 降级方案

### 自动降级

如果 Redis 连接失败，系统会自动降级到内存存储：

```go
// main.go
if err := config.InitRedis(); err != nil {
    log.Printf("Redis 连接失败: %v (将使用内存存储作为降级方案)", err)
}
```

### 验证码降级

```go
// utils/captcha.go
func VerifyCaptcha(id, answer string) bool {
    // 优先使用 Redis
    if VerifyCaptchaRedis(id, answer) {
        return true
    }
    // 降级到内存存储
    return captcha.VerifyString(id, answer)
}
```

## 🧪 测试

### 测试 Redis 连接

```bash
# 连接 Redis
redis-cli

# 测试连接
PING
# 应返回: PONG

# 查看所有键
KEYS *

# 查看验证码
KEYS captcha:*

# 查看用户缓存
KEYS user:*

# 获取键的值
GET captcha:abc123

# 查看键的过期时间
TTL captcha:abc123
```

### 测试缓存功能

```bash
# 1. 获取验证码（会存储到 Redis）
curl http://localhost:8080/api/v1/captcha

# 2. 在 Redis 中查看
redis-cli
> KEYS captcha:*
> GET captcha:{返回的id}

# 3. 获取用户（首次查询数据库）
curl http://localhost:8080/api/v1/users/1

# 4. 再次获取（从缓存读取）
curl http://localhost:8080/api/v1/users/1

# 5. 在 Redis 中查看
redis-cli
> GET user:1
```

## 📈 监控和维护

### Redis 监控命令

```bash
# 查看 Redis 信息
INFO

# 查看内存使用
INFO memory

# 查看统计信息
INFO stats

# 实时监控命令
MONITOR

# 查看慢查询
SLOWLOG GET 10
```

### 清理缓存

```bash
# 删除特定模式的键
redis-cli KEYS "user:*" | xargs redis-cli DEL

# 清空当前数据库（慎用）
redis-cli FLUSHDB

# 清空所有数据库（慎用）
redis-cli FLUSHALL
```

## 🎯 最佳实践

### 1. 键命名规范

```
{业务}:{类型}:{ID}
例如：
- captcha:abc123
- user:1
- token:blacklist:xyz789
- login:fail:192.168.1.1
```

### 2. 设置合理的过期时间

```go
// 避免永久缓存
utils.CacheSet("key", value, 0) // ❌ 错误

// 设置合理的过期时间
utils.CacheSet("key", value, 30*time.Minute) // ✅ 正确
```

### 3. 处理缓存穿透

```go
// 缓存空值防止穿透
if user == nil {
    utils.CacheSet(cacheKey, "null", 5*time.Minute)
    return nil, errors.New("用户不存在")
}
```

### 4. 使用连接池

```go
// 已在 config/redis.go 中配置
PoolSize:     10,
MinIdleConns: 5,
```

## 📚 相关文件

- `config/redis.go` - Redis 连接配置
- `utils/cache.go` - 缓存工具函数
- `utils/captcha.go` - 验证码 Redis 存储
- `services/user_service.go` - 用户信息缓存

## 🔄 下一步优化

- [ ] 实现登录失败次数限制
- [ ] 实现 JWT Token 黑名单
- [ ] 添加缓存预热功能
- [ ] 实现缓存统计和监控
- [ ] 添加分布式锁
- [ ] 实现 Redis 集群支持

---

**实现时间**: 2025-12-25  
**Redis 客户端**: github.com/redis/go-redis/v9  
**状态**: ✅ 生产就绪（支持降级）
