@echo off
cd /d "C:\Users\hyper\workspace\tormentnexus"
setlocal enabledelayedexpansion

echo === CodeWhale: TormentNexus Integration ===
echo.

where codewhale >nul 2>nul
if %errorlevel% neq 0 (
    echo ⏭️  CodeWhale not found. Skipping.
    exit /b 0
)

echo ✅ CodeWhale detected.

:: ── Step 1: Install the skill ──
echo.
echo Installing TormentNexus skill...
if not exist "%USERPROFILE%\.codewhale\skills\tormentnexus" mkdir "%USERPROFILE%\.codewhale\skills\tormentnexus"
copy /Y "C:\Users\hyper\workspace\tormentnexus\.codewhale\plugins\tormentnexus\skills\SKILL.md" "%USERPROFILE%\.codewhale\skills\tormentnexus\SKILL.md"
if %errorlevel%==0 (echo ✅ Skill installed) else (echo ⚠️ Skill copy failed)

:: ── Step 2: Install the plugin config ──
echo.
echo Installing plugin configuration...
if not exist "%USERPROFILE%\.codewhale\plugins\tormentnexus" mkdir "%USERPROFILE%\.codewhale\plugins\tormentnexus"
if not exist "%USERPROFILE%\.codewhale\plugins\tormentnexus\skills" mkdir "%USERPROFILE%\.codewhale\plugins\tormentnexus\skills"
copy /Y "C:\Users\hyper\workspace\tormentnexus\.codewhale\plugins\tormentnexus\plugin.toml" "%USERPROFILE%\.codewhale\plugins\tormentnexus\plugin.toml"
if %errorlevel%==0 (echo ✅ Plugin config installed) else (echo ⚠️ Plugin config copy failed)

:: ── Step 3: Ensure MCP server is registered ──
echo.
echo Checking MCP server registration...
codewhale mcp list 2>nul | findstr /I "tormentnexus" >nul
if %errorlevel%==0 (
    echo ✅ MCP server already registered
) else (
    echo Registering TormentNexus MCP server...
    codewhale mcp add "tormentnexus" ^
        --command "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" ^
        --arg "mcp" ^
        --env "TORMENTNEXUS_WORKSPACE_ROOT=C:\Users\hyper\workspace\tormentnexus" >nul 2>nul
    if !errorlevel!==0 (echo ✅ MCP server registered) else (echo ⚠️ MCP registration had issues)
)

:: ── Step 4: Verify ──
echo.
echo Verifying installation...
echo.
codewhale mcp list 2>nul | findstr /I "tormentnexus"
if %errorlevel%==0 (
    echo ✅ TormentNexus MCP is connected
) else (
    echo ⚠️ Cannot verify MCP connection
)

echo.
echo ========================================
echo  CodeWhale TormentNexus Install Complete
echo ========================================
echo  ✅ Skill:    ~\.codewhale\skills\tormentnexus\SKILL.md
echo  ✅ Plugin:   ~\.codewhale\plugins\tormentnexus\plugin.toml
echo  ✅ MCP:      ~\.codewhale\mcp.json (tormentnexus)
echo.
endlocal
