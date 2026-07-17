@echo off
setlocal enabledelayedexpansion

echo ========================================
echo   HyperNexus / TormentNexus Build Script
echo ========================================
echo.

:menu
echo Select edition to build:
echo.
echo   1. TormentNexus (Open Source)
echo   2. HyperNexus (Corporate)
echo   3. Both editions
echo   4. Exit
echo.
set /p CHOICE="Enter choice (1-4): "

if "%CHOICE%"=="1" goto build_tormentnexus
if "%CHOICE%"=="2" goto build_hypernexus
if "%CHOICE%"=="3" goto build_both
if "%CHOICE%"=="4" goto end
echo Invalid choice. Please try again.
goto menu

:build_tormentnexus
echo.
echo Building TormentNexus (Open Source Edition)...
echo.

REM Get version
for /f "tokens=*" %%v in ('type VERSION') do set VER=%%v
echo Version: %VER%

REM Build Go binary
echo Building Go kernel...
cd go
go build -ldflags "-s -w -X gitlab.com/robertpelloni/HyperNexus/internal/buildinfo.Version=%VER% -X gitlab.com/robertpelloni/HyperNexus/internal/config.DefaultEdition=tormentnexus" -buildvcs=false -o ..\bin\tormentnexus.exe ./cmd/tormentnexus
if errorlevel 1 (
    echo [FAIL] Go build failed
    cd ..
    exit /b 1
)
cd ..
echo OK - bin\tormentnexus.exe built

REM Build installer
echo Building installer...
cd installer
makensis /DTORMENTNEXUS=1 hypernexus.nsi
if errorlevel 1 (
    echo [FAIL] Installer build failed
    cd ..
    exit /b 1
)
cd ..
echo OK - tormentnexus-setup.exe created

echo.
echo ========================================
echo   TormentNexus build complete!
echo ========================================
goto end

:build_hypernexus
echo.
echo Building HyperNexus (Corporate Edition)...
echo.

REM Get version
for /f "tokens=*" %%v in ('type VERSION') do set VER=%%v
echo Version: %VER%

REM Build Go binary
echo Building Go kernel...
cd go
go build -ldflags "-s -w -X gitlab.com/robertpelloni/HyperNexus/internal/buildinfo.Version=%VER% -X gitlab.com/robertpelloni/HyperNexus/internal/config.DefaultEdition=hypernexus" -buildvcs=false -o ..\bin\hypernexus.exe ./cmd/tormentnexus
if errorlevel 1 (
    echo [FAIL] Go build failed
    cd ..
    exit /b 1
)
cd ..
echo OK - bin\hypernexus.exe built

REM Build installer
echo Building installer...
cd installer
makensis /DCORPORATE=1 hypernexus.nsi
if errorlevel 1 (
    echo [FAIL] Installer build failed
    cd ..
    exit /b 1
)
cd ..
echo OK - hypernexus-setup.exe created

echo.
echo ========================================
echo   HyperNexus build complete!
echo ========================================
goto end

:build_both
call :build_tormentnexus
call :build_hypernexus
goto end

:end
echo.
pause
