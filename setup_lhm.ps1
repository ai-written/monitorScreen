# LibreHardwareMonitor Setup Script
# Downloads portable LHM and pre-configures web server
param(
    [string]$OutDir = "lhm"
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  LibreHardwareMonitor Setup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$repo = "LibreHardwareMonitor/LibreHardwareMonitor"

Write-Host "[1/4] Fetching latest release info..."
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest" -TimeoutSec 15
    $version = $release.tag_name
    Write-Host "  Latest version: $version"
} catch {
    Write-Host "[!] Cannot reach GitHub API. Try direct download:" -ForegroundColor Red
    Write-Host "    https://github.com/LibreHardwareMonitor/LibreHardwareMonitor/releases/latest" -ForegroundColor Yellow
    Write-Host "    Download the portable .zip, extract to ./lhm/, then re-run this script." -ForegroundColor Yellow
    exit 1
}

Write-Host "[2/4] Finding portable .zip asset..."
$asset = $release.assets | Where-Object { $_.name -like "*portable*" -and $_.name -like "*.zip" } | Select-Object -First 1
if (-not $asset) {
    $asset = $release.assets | Where-Object { $_.name -like "*.zip" } | Select-Object -First 1
}
if (-not $asset) {
    Write-Host "[!] No .zip asset found in release." -ForegroundColor Red
    Write-Host "    Please download manually from:" -ForegroundColor Yellow
    Write-Host "    https://github.com/LibreHardwareMonitor/LibreHardwareMonitor/releases/latest" -ForegroundColor Yellow
    exit 1
}
Write-Host "  Asset: $($asset.name)"

$zipPath = "$env:TEMP\lhm-portable.zip"
Write-Host "[3/4] Downloading $($asset.name)..."
Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $zipPath -TimeoutSec 120
Write-Host "  Downloaded to $zipPath"

Write-Host "[4/4] Extracting to $OutDir..."
if (Test-Path $OutDir) {
    Remove-Item -Recurse -Force $OutDir
}
Expand-Archive -Path $zipPath -DestinationPath $OutDir -Force
Remove-Item $zipPath

$configPath = Join-Path $OutDir "LibreHardwareMonitor.config.json"
$config = @{
    RemoteWebServerEnabled = $true
    RemoteWebServerPort    = 8085
    AutoUpdateEnabled      = $false
    MinimizeToTray         = $true
}

$config | ConvertTo-Json -Depth 2 | Set-Content -Path $configPath -Encoding UTF8

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  LHM setup complete!" -ForegroundColor Green
Write-Host "  Location: $((Get-Item $OutDir).FullName)" -ForegroundColor White
Write-Host "  Web server: http://localhost:8085/data.json" -ForegroundColor White
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Next: Set in config.yaml:" -ForegroundColor Yellow
Write-Host "  lhm:" -ForegroundColor White
Write-Host "    enabled: 'true'" -ForegroundColor White

if ($PSScriptRoot) {
    $targetDir = Split-Path $PSScriptRoot -Parent
    Write-Host ""
    Write-Host "Run this to copy LHM to project root:" -ForegroundColor Cyan
    Write-Host "  Copy-Item -Recurse -Force '$((Get-Item $OutDir).FullName)' '$targetDir\lhm'" -ForegroundColor White
}
