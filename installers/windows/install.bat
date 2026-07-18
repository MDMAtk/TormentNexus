@echo off
setlocal enabledelayedexpansion

title HyperNexus Installer
color 0B

echo.
echo  ============================================================
echo  ============================================================
echo.
echo     ████████╗ ██████╗ ██████╗ ███╗   ███╗███████╗███╗   ██╗████████╗
echo     ╚══██╔══╝██╔═══██╗██╔══██╗████╗ ████║██╔════╝████╗  ██║╚══██╔══╝
echo        ██║   ██║   ██║██████╔╝██╔████╔██║█████╗  ██╔██╗ ██║   ██║
echo        ██║   ██║   ██║██╔══██╗██║╚██╔╝██║██╔══╝  ██║╚██╗██║   ██║
echo        ██║   ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗██║ ╚████║   ██║
echo        ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝
echo                    N E X U S   I N S T A L L E R
echo.
echo  ============================================================
echo.
echo  AI Control Plane with Persistent Memory
echo  26,000+ MCP Tools ^| Multi-Agent Orchestration
echo.
echo  ============================================================
echo.

:: Check if running as administrator
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo  [!] This installer requires administrator privileges.
    echo      Right-click and select "Run as administrator"
    echo.
    pause
    exit /b 1
)

:: Set installation directory
set "INSTALL_DIR=%USERPROFILE%\.hypernexus"
set "BINARY_NAME=hypernexus.exe"

echo  [1/6] Creating installation directory...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
echo        OK - %INSTALL_DIR%
echo.

:: Copy binary
echo  [2/6] Installing HyperNexus binary...
copy /y "%~dp0%BINARY_NAME%" "%INSTALL_DIR%\%BINARY_NAME%" >nul 2>&1
if %errorlevel% neq 0 (
    echo  [!] Failed to copy binary. Please run as administrator.
    pause
    exit /b 1
)
echo        OK - %INSTALL_DIR%\%BINARY_NAME%
echo.

:: Create config directory
echo  [3/6] Creating configuration directory...
if not exist "%INSTALL_DIR%\.hypernexus" mkdir "%INSTALL_DIR%\.hypernexus"
echo        OK - %INSTALL_DIR%\.hypernexus
echo.

:: Create default config
echo  [4/6] Creating default configuration...
(
echo {
echo   "version": "1.0.0",
echo   "server": {
echo     "host": "127.0.0.1",
echo     "port": 7778
echo   },
echo   "memory": {
echo     "enabled": true,
echo     "tiers": ["L1", "L2", "L3", "L4"]
echo   },
echo   "mcp": {
echo     "catalog": true,
echo     "autoInstall": false
echo   }
echo }
) > "%INSTALL_DIR%\.hypernexus\config.json"
echo        OK - config.json created
echo.

:: Add to PATH
echo  [5/6] Adding to system PATH...
setx PATH "%PATH%;%INSTALL_DIR%" /M >nul 2>&1
echo        OK - Added to PATH
echo.

:: Create start menu shortcut
echo  [6/6] Creating start menu shortcut...
set "START_MENU=%APPDATA%\Microsoft\Windows\Start Menu\Programs\HyperNexus"
if not exist "%START_MENU%" mkdir "%START_MENU%"

(
echo @echo off
echo title HyperNexus Server
echo cd /d "%INSTALL_DIR%"
echo "%INSTALL_DIR%\%BINARY_NAME%" serve
echo pause
) > "%START_MENU%\Start HyperNexus.bat"

(
echo @echo off
echo title HyperNexus Dashboard
echo start http://localhost:7778
echo.
) > "%START_MENU%\Open Dashboard.bat"

echo        OK - Start menu shortcuts created
echo.

echo  ============================================================
echo.
echo  INSTALLATION COMPLETE!
echo.
echo  ============================================================
echo.
echo  HyperNexus has been installed to:
echo    %INSTALL_DIR%
echo.
echo  To start HyperNexus:
echo    1. Open a new command prompt (PATH must refresh)
echo    2. Run: hypernexus serve
echo    3. Open: http://localhost:7778
echo.
echo  Or use the Start Menu shortcuts:
echo    - Start HyperNexus
echo    - Open Dashboard
echo.
echo  ============================================================
echo.

:: Ask to start now
set /p START_NOW="  Start HyperNexus now? (Y/N): "
if /i "%START_NOW%"=="Y" (
    echo.
    echo  Starting HyperNexus...
    start "HyperNexus" cmd /k "cd /d "%INSTALL_DIR%" && "%INSTALL_DIR%\%BINARY_NAME%" serve"
    timeout /t 3 >nul
    start http://localhost:7778
)

echo.
echo  Press any key to exit...
pause >nul
