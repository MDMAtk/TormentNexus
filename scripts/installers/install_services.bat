@echo off
cd /d "C:\Users\hyper\workspace\tormentnexus"
setlocal enabledelayedexpansion

echo ========================================
echo  TormentNexus Service Registration
echo  Run this as Administrator!
echo ========================================
echo.

echo === Step 1: Windows Services ===
echo.
echo Registering Go Sidecar (port 7778)...
sc create "TormentNexusSidecar" binPath="\"C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe\" serve" start=auto displayname="TormentNexus Sidecar"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)
echo.

echo Registering Dashboard (port 7779)...
sc create "TormentNexusDashboard" binPath="\"C:\Program Files\nodejs\node.exe\" \"C:\Users\hyper\workspace\tormentnexus\apps\web\node_modules\.bin\next.cmd\" dev -p 7779" start=auto displayname="TormentNexus Dashboard"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)
echo.

echo Registering Watchdog...
sc create "TormentNexusWatchdog" binPath="\"C:\Python314\pythonw.exe\" -u \"C:\Users\hyper\workspace\tormentnexus\watchdog.py\"" start=auto displayname="TormentNexus Watchdog"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)
echo.

echo === Step 2: Pi Coding Agent Extension ===
echo.
if not exist "%USERPROFILE%\.pi\agent\extensions" mkdir "%USERPROFILE%\.pi\agent\extensions"
copy /Y "C:\Users\hyper\workspace\tormentnexus\.pi\extensions\tormentnexus.ts" "%USERPROFILE%\.pi\agent\extensions\tormentnexus.ts"
if %errorlevel%==0 (echo ✅ Pi extension installed) else (echo ⚠️ Pi extension copy failed)
echo.

echo === Step 3: CodeWhale Integration ===
echo.
where codewhale >nul 2>nul
if %errorlevel%==0 (
    echo CodeWhale detected at:
    where codewhale

    rem --- Install MCP server config ---
    echo.
    echo Installing MCP server config...
    mkdir "%USERPROFILE%\.codewhale" >nul 2>nul
    if exist "%USERPROFILE%\.codewhale\mcp.json" (
        echo ✓ MCP config already exists, checking for tormentnexus entry...
        findstr /C:"tormentnexus" "%USERPROFILE%\.codewhale\mcp.json" >nul 2>nul
        if !errorlevel!==0 (
            echo ✓ tormentnexus MCP entry already configured
        ) else (
            echo ⚠️ tormentnexus entry not found in mcp.json — adding via codewhale CLI...
            codewhale mcp add "tormentnexus" --command "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" --arg "mcp"
        )
    ) else (
        echo Creating new MCP config...
        codewhale mcp add "tormentnexus" --command "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" --arg "mcp"
    )
    if %errorlevel%==0 (echo ✅ MCP server configured) else (echo ⚠️ MCP config may need manual setup)

    rem --- Install CodeWhale skill ---
    echo.
    echo Installing CodeWhale skill...
    if not exist "%USERPROFILE%\.codewhale\skills\tormentnexus" mkdir "%USERPROFILE%\.codewhale\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.codewhale\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.codewhale\skills\tormentnexus\SKILL.md"
    if %errorlevel%==0 (echo ✅ CodeWhale skill installed) else (echo ⚠️ Skill copy failed)
) else (
    echo CodeWhale not found — skipping CodeWhale integration.
    echo Install CodeWhale with: npm install -g codewhale
    echo Then re-run this installer.
)
echo.

echo === Step 4: Starting Services ===
echo.
sc start TormentNexusSidecar
sc start TormentNexusDashboard
sc start TormentNexusWatchdog
echo.
echo ========================================
echo  Done!
echo.
echo  TormentNexus Pi extension:    ~\.pi\agent\extensions\tormentnexus.ts
echo  CodeWhale skill:              ~\.codewhale\skills\tormentnexus\SKILL.md
echo  CodeWhale MCP config:         ~\.codewhale\mcp.json
echo ========================================
pause
