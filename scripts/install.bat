@echo off
setlocal enabledelayedexpansion

echo.
echo ========================================
echo   HyperNexus Installer v1.0.0
echo ========================================
echo.

:: Check if running as administrator
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo [!] This installer requires administrator privileges.
    echo     Right-click and select "Run as administrator"
    pause
    exit /b 1
)

:: Set installation directory
set "INSTALL_DIR=%ProgramFiles%\HyperNexus"
set "BIN_DIR=%INSTALL_DIR%\bin"
set "CONFIG_DIR=%USERPROFILE%\.hypernexus"

echo [*] Installing HyperNexus to %INSTALL_DIR%
echo.

:: Create directories
mkdir "%INSTALL_DIR%" 2>nul
mkdir "%BIN_DIR%" 2>nul
mkdir "%CONFIG_DIR%" 2>nul
mkdir "%CONFIG_DIR%\memory" 2>nul

:: Copy binaries
echo [*] Copying binaries...
copy /Y "%~dp0hypernexus.exe" "%BIN_DIR%\hypernexus.exe" >nul
if %errorlevel% neq 0 (
    echo [!] Failed to copy hypernexus.exe
    pause
    exit /b 1
)

:: Create start menu shortcut
echo [*] Creating start menu shortcut...
set "SHORTCUT_DIR=%APPDATA%\Microsoft\Windows\Start Menu\Programs\HyperNexus"
mkdir "%SHORTCUT_DIR%" 2>nul

:: Create VBS script for shortcut creation
echo Set oWS = WScript.CreateObject("WScript.Shell") > "%TEMP%\CreateShortcut.vbs"
echo sLinkFile = "%SHORTCUT_DIR%\HyperNexus.lnk" >> "%TEMP%\CreateShortcut.vbs"
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> "%TEMP%\CreateShortcut.vbs"
echo oLink.TargetPath = "%BIN_DIR%\hypernexus.exe" >> "%TEMP%\CreateShortcut.vbs"
echo oLink.Arguments = "serve" >> "%TEMP%\CreateShortcut.vbs"
echo oLink.WorkingDirectory = "%INSTALL_DIR%" >> "%TEMP%\CreateShortcut.vbs"
echo oLink.Description = "HyperNexus AI Control Plane" >> "%TEMP%\CreateShortcut.vbs"
echo oLink.Save >> "%TEMP%\CreateShortcut.vbs"
cscript //nologo "%TEMP%\CreateShortcut.vbs"
del "%TEMP%\CreateShortcut.vbs"

:: Add to PATH
echo [*] Adding to PATH...
setx PATH "%PATH%;%BIN_DIR%" /M >nul 2>&1

:: Create default config
echo [*] Creating default configuration...
(
echo # HyperNexus Configuration
echo host: 127.0.0.1
echo port: 7778
echo workspace: %USERPROFILE%\workspace\hypernexus
echo.
echo # Memory Configuration
echo memory:
echo   l2_enabled: true
echo   l3_enabled: true
echo   l4_enabled: false
echo.
echo # Provider Configuration  
echo providers:
echo   deepseek:
echo     enabled: true
echo     api_key: ""
echo   lmstudio:
echo     enabled: true
echo     url: http://127.0.0.1:1234
) > "%CONFIG_DIR%\config.yaml"

:: Create uninstaller
echo [*] Creating uninstaller...
(
echo @echo off
echo echo Uninstalling HyperNexus...
echo taskkill /f /im hypernexus.exe 2^>nul
echo rmdir /s /q "%INSTALL_DIR%"
echo rmdir /s /q "%SHORTCUT_DIR%"
echo echo HyperNexus has been uninstalled.
echo pause
) > "%INSTALL_DIR%\uninstall.bat"

:: Create desktop shortcut
echo [*] Creating desktop shortcut...
echo Set oWS = WScript.CreateObject("WScript.Shell") > "%TEMP%\CreateDesktopShortcut.vbs"
echo sLinkFile = "%USERPROFILE%\Desktop\HyperNexus.lnk" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.TargetPath = "%BIN_DIR%\hypernexus.exe" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.Arguments = "serve" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.WorkingDirectory = "%INSTALL_DIR%" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.Description = "HyperNexus AI Control Plane" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.Save >> "%TEMP%\CreateDesktopShortcut.vbs"
cscript //nologo "%TEMP%\CreateDesktopShortcut.vbs"
del "%TEMP%\CreateDesktopShortcut.vbs"

echo.
echo ========================================
echo   Installation Complete!
echo ========================================
echo.
echo HyperNexus has been installed to:
echo   %INSTALL_DIR%
echo.
echo Configuration file:
echo   %CONFIG_DIR%\config.yaml
echo.
echo To start HyperNexus:
echo   1. Double-click the desktop shortcut, or
echo   2. Run "hypernexus serve" from command line
echo.
echo Dashboard will be available at:
echo   http://127.0.0.1:7778
echo.
echo To uninstall:
echo   Run "%INSTALL_DIR%\uninstall.bat"
echo.

:: Ask to start now
set /p START_NOW="Start HyperNexus now? (Y/N): "
if /i "%START_NOW%"=="Y" (
    echo [*] Starting HyperNexus...
    start "" "%BIN_DIR%\hypernexus.exe" serve
    echo [*] HyperNexus is starting...
    echo [*] Dashboard: http://127.0.0.1:7778
    timeout /t 3 >nul
    start http://127.0.0.1:7778
)

pause
