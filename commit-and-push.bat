@echo off
setlocal enabledelayedexpansion

echo ========================================
echo   Git Commit and Push to GitLab
echo ========================================
echo.

REM Check if git is available
where git >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Git is not installed or not in PATH
    pause
    exit /b 1
)

REM Show current status
echo Current git status:
echo ----------------------------------------
git status --short
echo ----------------------------------------
echo.

REM Stage all changes
echo Staging all changes...
git add -A
if errorlevel 1 (
    echo [ERROR] Failed to stage changes
    pause
    exit /b 1
)
echo Done staging changes.
echo.

REM Show what will be committed
echo Changes to be committed:
echo ----------------------------------------
git diff --cached --stat
echo ----------------------------------------
echo.

REM Get commit message
set /p COMMIT_MSG="Enter commit message (or press Enter for default): "
if "%COMMIT_MSG%"=="" (
    set "COMMIT_MSG=feat: corporate mode - HyperNexus/TormentNexus dual edition support

- Added branding configuration system (go/internal/config/branding.go)
- Updated system tray to use dynamic branding
- Created unified installer supporting both editions
- Added build script for both editions
- Added corporate mode documentation
- Added sample configuration files
- Removed GPL license, kept MIT license
- Updated all references to use GitLab repository"
)

REM Commit changes
echo.
echo Committing changes...
git commit -m "%COMMIT_MSG%"
if errorlevel 1 (
    echo [ERROR] Failed to commit changes
    pause
    exit /b 1
)
echo Commit successful!
echo.

REM Check remote configuration
echo Checking remote configuration...
echo ----------------------------------------
git remote -v
echo ----------------------------------------
echo.

REM Confirm push
echo.
echo This will push to GitLab ONLY.
set /p CONFIRM="Continue? (Y/N): "
if /i not "%CONFIRM%"=="Y" (
    echo Push cancelled.
    pause
    exit /b 0
)

REM Push to GitLab (origin)
echo.
echo Pushing to GitLab...
git push origin main
if errorlevel 1 (
    echo [WARNING] Push to origin failed, trying with force...
    set /p FORCE="Force push? (Y/N): "
    if /i "!FORCE!"=="Y" (
        git push origin main --force
        if errorlevel 1 (
            echo [ERROR] Force push failed
            pause
            exit /b 1
        )
        echo Force push successful!
    ) else (
        echo Push cancelled.
        pause
        exit /b 0
    )
) else (
    echo Push successful!
)

echo.
echo ========================================
echo   All done! Changes pushed to GitLab.
echo ========================================
echo.
pause
