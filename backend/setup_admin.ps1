# 数据库检查和超级管理员设置脚本
# 使用方法: .\setup_admin.ps1

Write-Host "=== 数据库检查和超级管理员设置 ===" -ForegroundColor Green
Write-Host ""

$DB_HOST = "localhost"
$DB_PORT = "3306"
$DB_USER = "root"
$DB_PASSWORD = "123456"
$DB_NAME = "projecttest"

Write-Host "数据库: $DB_NAME" -ForegroundColor Cyan
Write-Host ""

# 1. 检查菜单数据
Write-Host "1. 检查菜单数据..." -ForegroundColor Yellow
$menuQuery = "SELECT COUNT(*) as count FROM menus;"
$menuCount = mysql -h $DB_HOST -P $DB_PORT -u $DB_USER "-p$DB_PASSWORD" $DB_NAME -N -e $menuQuery
Write-Host "   菜单数量: $menuCount" -ForegroundColor White

# 2. 检查角色数据
Write-Host "2. 检查角色数据..." -ForegroundColor Yellow
$roleQuery = "SELECT COUNT(*) as count FROM roles;"
$roleCount = mysql -h $DB_HOST -P $DB_PORT -u $DB_USER "-p$DB_PASSWORD" $DB_NAME -N -e $roleQuery
Write-Host "   角色数量: $roleCount" -ForegroundColor White

# 3. 检查用户数据
Write-Host "3. 检查用户数据..." -ForegroundColor Yellow
$userQuery = "SELECT id, username, email, role_id FROM users;"
Write-Host "   用户列表:" -ForegroundColor White
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER "-p$DB_PASSWORD" $DB_NAME -e $userQuery

Write-Host ""

# 4. 如果没有菜单和角色，执行初始化
if ($menuCount -eq 0 -or $roleCount -eq 0) {
    Write-Host "检测到数据不完整，正在初始化..." -ForegroundColor Yellow
    Get-Content init_data.sql | mysql -h $DB_HOST -P $DB_PORT -u $DB_USER "-p$DB_PASSWORD" $DB_NAME
    Write-Host "初始化完成!" -ForegroundColor Green
    Write-Host ""
}

# 5. 设置超级管理员
Write-Host "5. 设置超级管理员" -ForegroundColor Yellow
$username = Read-Host "请输入要设置为超级管理员的用户名"

if ($username) {
    $updateQuery = "UPDATE users SET role_id = 1 WHERE username = '$username';"
    mysql -h $DB_HOST -P $DB_PORT -u $DB_USER "-p$DB_PASSWORD" $DB_NAME -e $updateQuery
    
    # 验证
    $verifyQuery = "SELECT id, username, email, role_id FROM users WHERE username = '$username';"
    Write-Host ""
    Write-Host "设置结果:" -ForegroundColor Cyan
    mysql -h $DB_HOST -P $DB_PORT -u $DB_USER "-p$DB_PASSWORD" $DB_NAME -e $verifyQuery
    
    Write-Host ""
    Write-Host "✅ 超级管理员设置成功!" -ForegroundColor Green
    Write-Host ""
    Write-Host "下一步:" -ForegroundColor Yellow
    Write-Host "  1. 重新登录系统" -ForegroundColor White
    Write-Host "  2. 你将看到所有菜单（包括系统管理）" -ForegroundColor White
} else {
    Write-Host "未输入用户名，跳过设置" -ForegroundColor Red
}

Write-Host ""
Write-Host "=== 完成 ===" -ForegroundColor Green
Write-Host ""
Read-Host "按回车键退出"
