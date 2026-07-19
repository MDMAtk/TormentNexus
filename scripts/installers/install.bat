@echo off
setlocal enabledelayedexpansion

echo.
echo ========================================
echo   TormentNexus Installer v1.0.0
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
set "INSTALL_DIR=%ProgramFiles%\TormentNexus"
set "BIN_DIR=%INSTALL_DIR%\bin"
set "CONFIG_DIR=%USERPROFILE%\.tormentnexus"

echo [*] Installing TormentNexus to %INSTALL_DIR%
echo.

:: Create directories
mkdir "%INSTALL_DIR%" 2>nul
mkdir "%BIN_DIR%" 2>nul
mkdir "%CONFIG_DIR%" 2>nul
mkdir "%CONFIG_DIR%\memory" 2>nul

:: Copy binaries
echo [*] Copying binaries...
copy /Y "%~dp0tormentnexus.exe" "%BIN_DIR%\tormentnexus.exe" >nul
if %errorlevel% neq 0 (
    echo [!] Failed to copy tormentnexus.exe
    pause
    exit /b 1
)

:: Create start menu shortcut
echo [*] Creating start menu shortcut...
set "SHORTCUT_DIR=%APPDATA%\Microsoft\Windows\Start Menu\Programs\TormentNexus"
mkdir "%SHORTCUT_DIR%" 2>nul

:: Create VBS script for shortcut creation
echo Set oWS = WScript.CreateObject("WScript.Shell") > "%TEMP%\CreateShortcut.vbs"
echo sLinkFile = "%SHORTCUT_DIR%\TormentNexus.lnk" >> "%TEMP%\CreateShortcut.vbs"
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> "%TEMP%\CreateShortcut.vbs"
echo oLink.TargetPath = "%BIN_DIR%\tormentnexus.exe" >> "%TEMP%\CreateShortcut.vbs"
echo oLink.Arguments = "serve" >> "%TEMP%\CreateShortcut.vbs"
echo oLink.WorkingDirectory = "%INSTALL_DIR%" >> "%TEMP%\CreateShortcut.vbs"
echo oLink.Description = "TormentNexus AI Control Plane" >> "%TEMP%\CreateShortcut.vbs"
echo oLink.Save >> "%TEMP%\CreateShortcut.vbs"
cscript //nologo "%TEMP%\CreateShortcut.vbs"
del "%TEMP%\CreateShortcut.vbs"

:: Add to PATH
echo [*] Adding to PATH...
setx PATH "%PATH%;%BIN_DIR%" /M >nul 2>&1

:: Create default config
echo [*] Creating default configuration...
(
echo # TormentNexus Configuration
echo host: 127.0.0.1
echo port: 7778
echo workspace: %USERPROFILE%\workspace\tormentnexus
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
echo echo Uninstalling TormentNexus...
echo taskkill /f /im tormentnexus.exe 2^>nul
echo rmdir /s /q "%INSTALL_DIR%"
echo rmdir /s /q "%SHORTCUT_DIR%"
echo echo TormentNexus has been uninstalled.
echo pause
) > "%INSTALL_DIR%\uninstall.bat"

:: Create desktop shortcut
echo [*] Creating desktop shortcut...
echo Set oWS = WScript.CreateObject("WScript.Shell") > "%TEMP%\CreateDesktopShortcut.vbs"
echo sLinkFile = "%USERPROFILE%\Desktop\TormentNexus.lnk" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.TargetPath = "%BIN_DIR%\tormentnexus.exe" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.Arguments = "serve" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.WorkingDirectory = "%INSTALL_DIR%" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.Description = "TormentNexus AI Control Plane" >> "%TEMP%\CreateDesktopShortcut.vbs"
echo oLink.Save >> "%TEMP%\CreateDesktopShortcut.vbs"
cscript //nologo "%TEMP%\CreateDesktopShortcut.vbs"
del "%TEMP%\CreateDesktopShortcut.vbs"

echo.
echo ========================================
echo   Installation Complete!
echo ========================================
echo.
echo TormentNexus has been installed to:
echo   %INSTALL_DIR%
echo.
echo Configuration file:
echo   %CONFIG_DIR%\config.yaml
echo.
echo To start TormentNexus:
echo   1. Double-click the desktop shortcut, or
echo   2. Run "tormentnexus serve" from command line
echo.
echo Dashboard will be available at:
echo   http://127.0.0.1:7778
echo.
echo To uninstall:
echo   Run "%INSTALL_DIR%\uninstall.bat"
echo.

:: Ask to start now
set /p START_NOW="Start TormentNexus now? (Y/N): "
if /i "%START_NOW%"=="Y" (
    echo [*] Starting TormentNexus...
    start "" "%BIN_DIR%\tormentnexus.exe" serve
    echo [*] TormentNexus is starting...
    echo [*] Dashboard: http://127.0.0.1:7778
    timeout /t 3 >nul
    start http://127.0.0.1:7778
)

pause
