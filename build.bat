@echo off
echo ========================================
echo   Monitor Screen - Wails Build Script
echo ========================================
echo.

set GOPROXY=https://goproxy.cn,direct
echo [GOPROXY=%GOPROXY%]

REM Set app icon
if not exist "favicon.png" (
    echo [Icon] favicon.png not found, icon may not display correctly
)

where wails >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [!] Wails CLI not found. Installing...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if %ERRORLEVEL% neq 0 (
        echo [!] Failed to install Wails.
        pause
        exit /b 1
    )
)

echo [1/3] Running wails doctor...
wails doctor

echo [2/3] Installing frontend dependencies...
cd frontend
call npm install
if %ERRORLEVEL% neq 0 (
    echo [!] npm install failed.
    cd ..
    pause
    exit /b 1
)
cd ..

echo [3/4] Building sensor bridge (if LHM DLL present)...
if exist "lhm\LibreHardwareMonitorLib.dll" (
    call build_bridge.bat
    echo.
)

echo [4/4] Building Wails application...

if not exist "build\windows" mkdir "build\windows"
copy /Y "icon.ico" "build\windows\icon.ico" >nul
(
echo ^<?xml version="1.0" encoding="UTF-8" standalone="yes"?^>
echo ^<assembly manifestVersion="1.0" xmlns="urn:schemas-microsoft-com:asm.v1" xmlns:asmv3="urn:schemas-microsoft-com:asm.v3"^>
echo     ^<assemblyIdentity type="win32" name="com.wails.{{.Name}}" version="{{.Info.ProductVersion}}.0" processorArchitecture="*"/^>
echo     ^<dependency^>
echo         ^<dependentAssembly^>
echo             ^<assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/^>
echo         ^</dependentAssembly^>
echo     ^</dependency^>
echo     ^<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3"^>
echo         ^<security^>
echo             ^<requestedPrivileges^>
echo                 ^<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/^>
echo             ^</requestedPrivileges^>
echo         ^</security^>
echo     ^</trustInfo^>
echo     ^<asmv3:application^>
echo         ^<asmv3:windowsSettings^>
echo             ^<dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings"^>true/pm^</dpiAware^>
echo             ^<dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings"^>permonitorv2,permonitor^</dpiAwareness^>
echo         ^</asmv3:windowsSettings^>
echo     ^</asmv3:application^>
echo ^</assembly^>
) > "build\windows\wails.exe.manifest"

wails build -ldflags="-s -w"

echo.
echo Copying files to output directory...

if not exist "build\bin\monitorScreen.exe" (
    echo [!] Build output not found.
    pause
    exit /b 1
)

if not exist "build\bin" mkdir "build\bin"
copy /Y "config.yaml" "build\bin\config.yaml" >nul

if exist "lhm" (
    if not exist "build\bin\lhm" mkdir "build\bin\lhm"
    xcopy /E /Y /Q "lhm\*" "build\bin\lhm\"
    echo   lhm\ -^> build\bin\lhm\
)

echo.
echo ========================================
echo   Build complete!
echo   Output: .\build\bin\monitorScreen.exe
echo ========================================
pause
