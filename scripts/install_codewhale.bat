@echo off
cd /d "C:\Users\hyper\workspace\hypernexus"
setlocal enabledelayedexpansion

echo === CodeWhale: HyperNexus Integration ===
echo.

where codewhale >nul 2>nul
if %errorlevel% neq 0 (
    echo ⏭️  CodeWhale not found. Skipping.
    exit /b 0
)

echo ✅ CodeWhale detected.

:: ── Step 1: Install the skill ──
echo.
echo Installing HyperNexus skill...
if not exist "%USERPROFILE%\.codewhale\skills\hypernexus" mkdir "%USERPROFILE%\.codewhale\skills\hypernexus"
copy /Y "C:\Users\hyper\workspace\hypernexus\.codewhale\plugins\hypernexus\skills\SKILL.md" "%USERPROFILE%\.codewhale\skills\hypernexus\SKILL.md"
if %errorlevel%==0 (echo ✅ Skill installed) else (echo ⚠️ Skill copy failed)

:: ── Step 2: Install the plugin config ──
echo.
echo Installing plugin configuration...
if not exist "%USERPROFILE%\.codewhale\plugins\hypernexus" mkdir "%USERPROFILE%\.codewhale\plugins\hypernexus"
if not exist "%USERPROFILE%\.codewhale\plugins\hypernexus\skills" mkdir "%USERPROFILE%\.codewhale\plugins\hypernexus\skills"
copy /Y "C:\Users\hyper\workspace\hypernexus\.codewhale\plugins\hypernexus\plugin.toml" "%USERPROFILE%\.codewhale\plugins\hypernexus\plugin.toml"
if %errorlevel%==0 (echo ✅ Plugin config installed) else (echo ⚠️ Plugin config copy failed)

:: ── Step 3: Ensure MCP server is registered ──
echo.
echo Checking MCP server registration...
codewhale mcp list 2>nul | findstr /I "hypernexus" >nul
if %errorlevel%==0 (
    echo ✅ MCP server already registered
) else (
    echo Registering HyperNexus MCP server...
    codewhale mcp add "hypernexus" ^
        --command "C:\Users\hyper\workspace\hypernexus\hypernexus.exe" ^
        --arg "mcp" ^
        --env "HYPERNEXUS_WORKSPACE_ROOT=C:\Users\hyper\workspace\hypernexus" >nul 2>nul
    if !errorlevel!==0 (echo ✅ MCP server registered) else (echo ⚠️ MCP registration had issues)
)

:: ── Step 4: Verify ──
echo.
echo Verifying installation...
echo.
codewhale mcp list 2>nul | findstr /I "hypernexus"
if %errorlevel%==0 (
    echo ✅ HyperNexus MCP is connected
) else (
    echo ⚠️ Cannot verify MCP connection
)

echo.
echo ========================================
echo  CodeWhale HyperNexus Install Complete
echo ========================================
echo  ✅ Skill:    ~\.codewhale\skills\hypernexus\SKILL.md
echo  ✅ Plugin:   ~\.codewhale\plugins\hypernexus\plugin.toml
echo  ✅ MCP:      ~\.codewhale\mcp.json (hypernexus)
echo.
endlocal
