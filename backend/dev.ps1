$AIR_PATH = "C:\Users\jiacan\go\bin\air.exe"
if (Test-Path $AIR_PATH) {
    & $AIR_PATH
} else {
    Write-Host "Air is not found at $AIR_PATH. Trying to run 'air' from PATH..."
    air
}
