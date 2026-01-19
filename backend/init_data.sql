-- 数据库初始化脚本（支持重复执行）
-- 如果数据已存在，会自动跳过

-- 清空现有数据（可选，如果需要重新初始化请取消注释）
-- DELETE FROM role_menus;
-- DELETE FROM menus;
-- DELETE FROM roles;

-- 插入菜单数据（如果已存在则忽略）
INSERT IGNORE INTO menus (id, parent_id, name, path, component, icon, sort, type, status, hidden, created_at, updated_at) VALUES
(1, 0, '首页', '/', 'AntdHome', 'HomeOutlined', 1, 1, 1, 0, NOW(), NOW()),
(2, 0, '订单管理', '/orders', 'Orders', 'ShoppingCartOutlined', 2, 1, 1, 0, NOW(), NOW()),
(3, 0, '系统管理', '/system', '', 'SettingOutlined', 3, 1, 1, 0, NOW(), NOW()),
(4, 3, '用户管理', '/system/users', 'SystemUsers', 'UserOutlined', 1, 1, 1, 0, NOW(), NOW()),
(5, 3, '角色管理', '/system/roles', 'SystemRoles', 'TeamOutlined', 2, 1, 1, 0, NOW(), NOW()),
(6, 3, '菜单管理', '/system/menus', 'SystemMenus', 'MenuOutlined', 3, 1, 1, 0, NOW(), NOW());

-- 插入角色数据（如果已存在则忽略）
INSERT IGNORE INTO roles (id, name, code, description, is_super, status, created_at, updated_at) VALUES
(1, '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, 1, NOW(), NOW()),
(2, '普通用户', 'user', '普通用户角色', 0, 1, NOW(), NOW());

-- 为超级管理员分配所有菜单（如果已存在则忽略）
INSERT IGNORE INTO role_menus (role_id, menu_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6);

-- 为普通用户分配首页和订单管理（如果已存在则忽略）
INSERT IGNORE INTO role_menus (role_id, menu_id) VALUES
(2, 1), (2, 2);

-- 更新现有用户的角色（如果 role_id 为 NULL 或 0，则设置为普通用户）
UPDATE users SET role_id = 2 WHERE role_id IS NULL OR role_id = 0;

-- 查询验证
SELECT '=== 菜单列表 ===' as info;
SELECT id, parent_id, name, path, component, icon FROM menus ORDER BY parent_id, sort;

SELECT '=== 角色列表 ===' as info;
SELECT id, name, code, is_super FROM roles;

SELECT '=== 角色菜单关联 ===' as info;
SELECT r.name as role_name, m.name as menu_name 
FROM role_menus rm 
JOIN roles r ON rm.role_id = r.id 
JOIN menus m ON rm.menu_id = m.id 
ORDER BY r.id, m.id;

SELECT '=== 用户列表 ===' as info;
SELECT id, username, email, role_id FROM users;
