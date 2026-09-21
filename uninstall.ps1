# Tole Uninstaller for Windows
# Removes tole.exe and removes %LOCALAPPDATA%\Tole from User PATH

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "       Tole - Windows Uninstaller       " -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host ""

$installDir = Join-Path $env:LOCALAPPDATA "Tole"

# Remove from User PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath) {
    $paths = $currentPath.Split(";", [System.StringSplitOptions]::RemoveEmptyEntries) | Where-Object { $_ -ne $installDir }
    $newPath = $paths -join ";"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "[OK] Removed $installDir from your User PATH." -ForegroundColor Green
}

# Remove directory
if (Test-Path $installDir) {
    Remove-Item -Path $installDir -Recurse -Force
    Write-Host "[OK] Deleted $installDir" -ForegroundColor Green
}

Write-Host ""
Write-Host "Tole has been successfully uninstalled from your system." -ForegroundColor Green
Write-Host ""
