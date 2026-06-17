# LibreHardwareMonitor offline setup
# Place LibreHardwareMonitor.zip in the same directory as this script, then run.
param(
    [string]$ZipFile = "LibreHardwareMonitor.zip",
    [string]$OutDir = "lhm"
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $ZipFile)) {
    Write-Host "[!] $ZipFile not found." -ForegroundColor Red
    Write-Host "    Download from:" -ForegroundColor Yellow
    Write-Host "    https://github.com/LibreHardwareMonitor/LibreHardwareMonitor/releases/latest" -ForegroundColor White
    Write-Host "    (Choose the portable .zip, place it here, then re-run.)" -ForegroundColor Yellow
    exit 1
}

Write-Host "[1/2] Extracting to $OutDir..."
if (Test-Path $OutDir) {
    Remove-Item -Recurse -Force $OutDir
}
Expand-Archive -Path $ZipFile -DestinationPath $OutDir -Force

Write-Host "[2/2] Creating config with web server enabled..."
$config = @{
    RemoteWebServerEnabled = $true
    RemoteWebServerPort    = 8085
    MinimizeToTray         = $true
    AutoUpdateEnabled      = $false
}
$configPath = Join-Path $OutDir "LibreHardwareMonitor.config.json"
$config | ConvertTo-Json -Depth 2 | Set-Content -Path $configPath -Encoding UTF8

Write-Host ""
Write-Host "LHM setup complete!" -ForegroundColor Green
Write-Host "Location: $((Get-Item $OutDir).FullName)" -ForegroundColor White
Write-Host ""
Write-Host "Next step: edit config.yaml, set lhm.enabled to 'true'" -ForegroundColor Yellow
