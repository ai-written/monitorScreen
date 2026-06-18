@echo off
echo Building LHM sensor bridge (HTTP server mode)...
set CSC=%SystemRoot%\Microsoft.NET\Framework64\v4.0.30319\csc.exe
if not exist "%CSC%" set CSC=%SystemRoot%\Microsoft.NET\Framework\v4.0.30319\csc.exe
if not exist "%CSC%" (
    echo [!] C# compiler not found. Install .NET Framework 4.7.2 SDK.
    exit /b 1
)

if not exist "lhm\LibreHardwareMonitorLib.dll" (
    echo [!] lhm\LibreHardwareMonitorLib.dll not found.
    echo     Extract LHM to lhm\ directory first.
    exit /b 1
)

"%CSC%" /target:exe /out:lhm\sensor_bridge.exe /reference:lhm\LibreHardwareMonitorLib.dll lhm_bridge.cs
if %ERRORLEVEL% neq 0 (
    echo [!] Compilation failed.
    exit /b 1
)

echo Done: sensor_bridge.exe (HTTP server on port 16533)
exit /b 0
