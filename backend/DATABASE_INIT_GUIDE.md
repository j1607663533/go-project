# 数据库初始化指南

## 方法 1: 使用 MySQL 命令行

### Windows PowerShell

```powershell
# 进入 backend 目录
cd d:\project-test\ReactFunTime\backend

# 执行 SQL 文件
mysql -u root -p reactfuntime < init_data.sql
```

### 或者使用完整路径

```powershell
mysql -u root -p -e "USE reactfuntime; SOURCE d:/project-test/ReactFunTime/backend/init_data.sql"
```

## 方法 2: 使用 MySQL Workbench

1. 打开 MySQL Workbench
2. 连接到数据库
3. 选择 `reactfuntime` 数据库
4. 打开 `init_data.sql` 文件
5. 点击执行（闪电图标）

## 方法 3: 使用 phpMyAdmin

1. 打开 phpMyAdmin
2. 选择 `reactfuntime` 数据库
3. 点击 "SQL" 标签
4. 复制 `init_data.sql` 的内容
5. 粘贴并执行

## 方法 4: 使用 PowerShell 脚本

```powershell
# 运行提供的脚本
.\run_init_sql.ps1
```

## 方法 5: 手动执行 SQL

如果上述方法都不行，可以手动复制以下 SQL 并执行：

```sql
-- 插入菜单数据
INSERT INTO menus (id, parent_id, name, path, component, icon, sort, type, status, hidden, created_at, updated_at) VALUES
(1, 0, '首页', '/', 'AntdHome', 'HomeOutlined', 1, 1, 1, 0, NOW(), NOW()),
(2, 0, '订单管理', '/orders', 'Orders', 'ShoppingCartOutlined', 2, 1, 1, 0, NOW(), NOW()),
(3, 0, '系统管理', '/system', '', 'SettingOutlined', 3, 1, 1, 0, NOW(), NOW()),
(4, 3, '用户管理', '/system/users', 'SystemUsers', 'UserOutlined', 1, 1, 1, 0, NOW(), NOW()),
(5, 3, '角色管理', '/system/roles', 'SystemRoles', 'TeamOutlined', 2, 1, 1, 0, NOW(), NOW()),
(6, 3, '菜单管理', '/system/menus', 'SystemMenus', 'MenuOutlined', 3, 1, 1, 0, NOW(), NOW());

-- 插入角色数据
INSERT INTO roles (id, name, code, description, is_super, status, created_at, updated_at) VALUES
(1, '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, 1, NOW(), NOW()),
(2, '普通用户', 'user', '普通用户角色', 0, 1, NOW(), NOW());

-- 为超级管理员分配所有菜单
INSERT INTO role_menus (role_id, menu_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6);

-- 为普通用户分配首页和订单管理
INSERT INTO role_menus (role_id, menu_id) VALUES
(2, 1), (2, 2);
```

## 验证数据

执行完成后，运行以下 SQL 验证：

```sql
-- 查看菜单
SELECT id, parent_id, name, path FROM menus ORDER BY parent_id, sort;

-- 查看角色
SELECT id, name, code, is_super FROM roles;

-- 查看角色菜单关联
SELECT r.name as role_name, m.name as menu_name
FROM role_menus rm
JOIN roles r ON rm.role_id = r.id
JOIN menus m ON rm.menu_id = m.id
ORDER BY r.id, m.id;
```

## 创建超级管理员账号

1. 先在系统中注册一个账号
2. 然后执行以下 SQL：

```sql
-- 将用户设置为超级管理员
UPDATE users SET role_id = 1 WHERE username = 'your_username';

-- 验证
SELECT id, username, email, role_id FROM users WHERE username = 'your_username';
```

## 常见问题

### 1. 提示 "Duplicate entry" 错误

这是正常的，说明数据已经存在。可以忽略。

### 2. 提示 "Table doesn't exist" 错误

需要先运行后端程序，让 GORM 自动创建表：

```bash
cd backend
go run main.go
```

### 3. 无法连接数据库

检查 `backend/.env` 文件中的数据库配置是否正确。

### 4. 权限不足

确保 MySQL 用户有足够的权限：

```sql
GRANT ALL PRIVILEGES ON reactfuntime.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
```
