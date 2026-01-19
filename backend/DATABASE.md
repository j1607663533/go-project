# 数据库集成说明

## ✅ 数据库已成功集成

项目已成功连接到 MySQL 数据库，所有数据现在都持久化存储在数据库中。

## 数据库配置

### 当前配置

- **数据库类型**: MySQL
- **数据库名**: projectTest
- **主机**: localhost
- **端口**: 3306
- **用户名**: root
- **密码**: 123456

### 配置位置

配置信息在 `config/config.go` 中：

```go
DB: DatabaseConfig{
    Host:     getEnv("DB_HOST", "localhost"),
    Port:     getEnv("DB_PORT", "3306"),
    User:     getEnv("DB_USER", "root"),
    Password: getEnv("DB_PASSWORD", "123456"),
    DBName:   getEnv("DB_NAME", "projectTest"),
}
```

## 数据表结构

### users 表

| 字段       | 类型            | 说明         | 约束                        |
| ---------- | --------------- | ------------ | --------------------------- |
| id         | BIGINT UNSIGNED | 主键         | PRIMARY KEY, AUTO_INCREMENT |
| username   | VARCHAR(50)     | 用户名       | UNIQUE, NOT NULL            |
| email      | VARCHAR(100)    | 邮箱         | UNIQUE, NOT NULL            |
| password   | VARCHAR(255)    | 密码（加密） | NOT NULL                    |
| nickname   | VARCHAR(50)     | 昵称         | -                           |
| avatar     | VARCHAR(500)    | 头像 URL     | -                           |
| created_at | DATETIME        | 创建时间     | -                           |
| updated_at | DATETIME        | 更新时间     | -                           |

### 索引

- `idx_users_username`: username 字段的唯一索引
- `idx_users_email`: email 字段的唯一索引

## 架构说明

### 数据流向

```
Controller → Service → Repository → Database
    ↓          ↓           ↓            ↓
  HTTP      业务逻辑    数据访问      MySQL
```

### Repository 层实现

Repository 层使用 GORM 进行数据库操作：

```go
// 查询示例
func (r *userRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    if err := r.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// 创建示例
func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}
```

## 自动迁移

项目启动时会自动创建/更新数据表：

```go
// main.go
if err := config.AutoMigrate(&models.User{}); err != nil {
    log.Fatalf("数据库迁移失败: %v", err)
}
```

### 迁移规则

- 自动创建不存在的表
- 自动添加新字段
- **不会删除**已存在的字段
- **不会修改**已存在字段的类型

## 连接池配置

为了优化性能，已配置数据库连接池：

```go
sqlDB.SetMaxIdleConns(10)   // 最大空闲连接数
sqlDB.SetMaxOpenConns(100)  // 最大打开连接数
```

## 使用示例

### 1. 创建用户

```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "dbuser001",
  "email": "dbuser001@example.com",
  "password": "password123",
  "nickname": "数据库用户"
}
```

数据会自动保存到 `users` 表。

### 2. 查询用户

```bash
GET /api/v1/users/1
```

从数据库中查询 ID 为 1 的用户。

### 3. 更新用户

```bash
PUT /api/v1/users/1
Content-Type: application/json

{
  "nickname": "新昵称"
}
```

更新数据库中的用户信息。

### 4. 删除用户

```bash
DELETE /api/v1/users/1
```

从数据库中删除用户（软删除或硬删除取决于配置）。

## 事务支持

GORM 支持事务操作，可以在 Service 层使用：

```go
func (s *userService) CreateUserWithProfile(req *CreateUserRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 创建用户
        user := &models.User{...}
        if err := tx.Create(user).Error; err != nil {
            return err
        }

        // 创建用户资料
        profile := &models.Profile{UserID: user.ID}
        if err := tx.Create(profile).Error; err != nil {
            return err
        }

        return nil
    })
}
```

## 查询优化

### 预加载关联

```go
// 预加载用户的订单
var user models.User
db.Preload("Orders").First(&user, 1)
```

### 选择字段

```go
// 只查询需要的字段
var users []models.User
db.Select("id", "username", "email").Find(&users)
```

### 分页查询

```go
// 分页
var users []models.User
db.Limit(10).Offset(20).Find(&users)
```

## 数据库维护

### 查看表结构

```sql
DESC users;
```

### 查看索引

```sql
SHOW INDEX FROM users;
```

### 查看数据

```sql
SELECT * FROM users;
```

### 备份数据库

```bash
mysqldump -u root -p projectTest > backup.sql
```

### 恢复数据库

```bash
mysql -u root -p projectTest < backup.sql
```

## 性能优化建议

### 1. 使用索引

已为 `username` 和 `email` 创建唯一索引，查询这些字段会很快。

### 2. 避免 N+1 查询

使用 `Preload` 预加载关联数据：

```go
db.Preload("Orders").Find(&users)
```

### 3. 批量操作

使用批量插入而不是循环插入：

```go
db.CreateInBatches(users, 100)
```

### 4. 使用原生 SQL

对于复杂查询，可以使用原生 SQL：

```go
db.Raw("SELECT * FROM users WHERE created_at > ?", time.Now()).Scan(&users)
```

## 环境变量配置

如果需要在不同环境使用不同的数据库，可以创建 `.env` 文件：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=projectTest
```

**注意**: `.env` 文件已被 `.gitignore` 忽略，不会提交到 Git。

## 切换数据库

### 切换到 PostgreSQL

1. 安装驱动：

```bash
go get -u gorm.io/driver/postgres
```

2. 修改 `config/database.go`：

```go
import "gorm.io/driver/postgres"

dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
    AppConfig.DB.Host,
    AppConfig.DB.User,
    AppConfig.DB.Password,
    AppConfig.DB.DBName,
    AppConfig.DB.Port,
)
DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

### 切换到 SQLite

1. 安装驱动：

```bash
go get -u gorm.io/driver/sqlite
```

2. 修改 `config/database.go`：

```go
import "gorm.io/driver/sqlite"

DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
```

## 故障排查

### 连接失败

1. 检查 MySQL 服务是否运行
2. 检查用户名和密码是否正确
3. 检查数据库是否存在
4. 检查防火墙设置

### 迁移失败

1. 检查用户是否有创建表的权限
2. 检查字段定义是否正确
3. 查看详细错误日志

### 性能问题

1. 检查是否缺少索引
2. 检查是否有 N+1 查询
3. 使用 `db.Debug()` 查看生成的 SQL
4. 检查连接池配置

## 下一步

- [ ] 实现软删除（DeletedAt 字段）
- [ ] 添加更多索引优化查询
- [ ] 实现数据库读写分离
- [ ] 添加 Redis 缓存层
- [ ] 实现数据库备份策略
- [ ] 添加数据库监控

---

**数据库集成完成时间**: 2025-12-25  
**GORM 版本**: v1.31.1  
**MySQL 驱动版本**: v1.6.0
