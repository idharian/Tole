# Tole Installer for Windows
# Copies tole.exe to %LOCALAPPDATA%\Tole and adds it to User PATH

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "       Tole - Windows Installer         " -ForegroundColor Cyan
Write-Host "       Dev by agushariyanto             " -ForegroundColor Gray
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$sourceExe = Join-Path $PSScriptRoot "tole.exe"
if (-not (Test-Path $sourceExe)) {
    Write-Host "[*] Building tole.exe from source..." -ForegroundColor Yellow
    $cmdDir = Join-Path $PSScriptRoot "cmd\tole"
    & go build -o $sourceExe $cmdDir
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[-] Failed to build tole.exe. Please ensure Go is installed." -ForegroundColor Red
        exit 1
    }
}

$installDir = Join-Path $env:LOCALAPPDATA "Tole"
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

$targetExe = Join-Path $installDir "tole.exe"
Copy-Item -Path $sourceExe -Destination $targetExe -Force
Write-Host "[OK] Copied tole.exe to $installDir" -ForegroundColor Green

# Add to User PATH if not present
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($null -eq $currentPath) { $currentPath = "" }
$paths = $currentPath.Split(";", [System.StringSplitOptions]::RemoveEmptyEntries)

if ($paths -notcontains $installDir) {
    if ($currentPath -eq "") {
        $newPath = $installDir
    } else {
        $newPath = $currentPath + ";" + $installDir
    }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = $env:Path + ";" + $installDir
    Write-Host "[OK] Added $installDir to your User PATH." -ForegroundColor Green
} else {
    Write-Host "[OK] $installDir is already in your PATH." -ForegroundColor Gray
}

Write-Host ""
Write-Host "Installation successful!" -ForegroundColor Green
Write-Host "You can now run 'tole' from any terminal (PowerShell or Command Prompt)." -ForegroundColor Cyan
Write-Host "Try typing: tole" -ForegroundColor Yellow
Write-Host ""
