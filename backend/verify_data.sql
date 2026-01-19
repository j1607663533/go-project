-- 验证数据库数据
-- 在 MySQL 中执行此脚本来检查数据

USE projecttest;

-- 检查菜单
SELECT '=== 菜单数据 ===' as '';
SELECT COUNT(*) as menu_count FROM menus;
SELECT id, parent_id, name, path FROM menus ORDER BY parent_id, sort;

-- 检查角色
SELECT '=== 角色数据 ===' as '';
SELECT COUNT(*) as role_count FROM roles;
SELECT id, name, code, is_super FROM roles;

-- 检查角色菜单关联
SELECT '=== 角色菜单关联 ===' as '';
SELECT COUNT(*) as role_menu_count FROM role_menus;
SELECT r.name as role_name, m.name as menu_name 
FROM role_menus rm 
JOIN roles r ON rm.role_id = r.id 
JOIN menus m ON rm.menu_id = m.id 
ORDER BY r.id, m.id;

-- 检查用户
SELECT '=== 用户数据 ===' as '';
SELECT COUNT(*) as user_count FROM users;
SELECT id, username, email, role_id FROM users;
