# 执行 SQL 初始化脚本
# 使用方法: .\run_init_sql.ps1

Write-Host "=== 执行数据库初始化脚本 ===" -ForegroundColor Green

# 数据库配置（请根据实际情况修改）
$DB_HOST = "localhost"
$DB_PORT = "3306"
$DB_USER = "root"
$DB_NAME = "reactfuntime"
$SQL_FILE = "init_data.sql"

Write-Host "数据库: $DB_NAME" -ForegroundColor Cyan
Write-Host "SQL文件: $SQL_FILE" -ForegroundColor Cyan
Write-Host ""

# 提示输入密码
$DB_PASSWORD = Read-Host "请输入MySQL密码" -AsSecureString
$BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($DB_PASSWORD)
$PlainPassword = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)

Write-Host ""
Write-Host "正在执行SQL脚本..." -ForegroundColor Yellow

# 执行SQL文件
try {
    # 方法1: 使用 mysql 命令
    $mysqlCmd = "mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$PlainPassword $DB_NAME"
    Get-Content $SQL_FILE | & mysql -h $DB_HOST -P $DB_PORT -u $DB_USER "-p$PlainPassword" $DB_NAME
    
    Write-Host ""
    Write-Host "=== SQL脚本执行成功! ===" -ForegroundColor Green
    Write-Host ""
    Write-Host "已初始化:" -ForegroundColor Cyan
    Write-Host "  - 6个菜单" -ForegroundColor White
    Write-Host "  - 2个角色（超级管理员、普通用户）" -ForegroundColor White
    Write-Host "  - 角色-菜单关联" -ForegroundColor White
    Write-Host ""
    Write-Host "下一步:" -ForegroundColor Yellow
    Write-Host "  1. 注册一个新账号" -ForegroundColor White
    Write-Host "  2. 执行 SQL: UPDATE users SET role_id = 1 WHERE username = 'your_username';" -ForegroundColor White
    Write-Host "  3. 使用超级管理员账号登录" -ForegroundColor White
}
catch {
    Write-Host ""
    Write-Host "=== 执行失败 ===" -ForegroundColor Red
    Write-Host "错误信息: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "请检查:" -ForegroundColor Yellow
    Write-Host "  1. MySQL 服务是否已启动" -ForegroundColor White
    Write-Host "  2. 数据库 '$DB_NAME' 是否存在" -ForegroundColor White
    Write-Host "  3. 用户名和密码是否正确" -ForegroundColor White
    Write-Host "  4. mysql 命令是否在 PATH 中" -ForegroundColor White
}

Write-Host ""
Read-Host "按回车键退出"
