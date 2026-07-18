@echo off
cd /d "C:\Users\hyper\workspace\hypernexus"
setlocal enabledelayedexpansion

echo ========================================
echo  HyperNexus Service Registration
echo  Run this as Administrator!
echo ========================================
echo.

echo === Step 1: Windows Services ===
echo.
echo Registering Go Sidecar (port 7778)...
sc create "HyperNexusSidecar" binPath="\"C:\Users\hyper\workspace\hypernexus\hypernexus.exe\" serve" start=auto displayname="HyperNexus Sidecar"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)
echo.

echo Registering Dashboard (port 7779)...
sc create "HyperNexusDashboard" binPath="\"C:\Program Files\nodejs\node.exe\" \"C:\Users\hyper\workspace\hypernexus\apps\web\node_modules\.bin\next.cmd\" dev -p 7779" start=auto displayname="HyperNexus Dashboard"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)
echo.

echo Registering Watchdog...
sc create "HyperNexusWatchdog" binPath="\"C:\Python314\pythonw.exe\" -u \"C:\Users\hyper\workspace\hypernexus\watchdog.py\"" start=auto displayname="HyperNexus Watchdog"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)
echo.

echo === Step 2: Pi Coding Agent Extension ===
echo.
if not exist "%USERPROFILE%\.pi\agent\extensions" mkdir "%USERPROFILE%\.pi\agent\extensions"
copy /Y "C:\Users\hyper\workspace\hypernexus\.pi\extensions\hypernexus.ts" "%USERPROFILE%\.pi\agent\extensions\hypernexus.ts"
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
        echo ✓ MCP config already exists, checking for hypernexus entry...
        findstr /C:"hypernexus" "%USERPROFILE%\.codewhale\mcp.json" >nul 2>nul
        if !errorlevel!==0 (
            echo ✓ hypernexus MCP entry already configured
        ) else (
            echo ⚠️ hypernexus entry not found in mcp.json — adding via codewhale CLI...
            codewhale mcp add "hypernexus" --command "C:\Users\hyper\workspace\hypernexus\hypernexus.exe" --arg "mcp"
        )
    ) else (
        echo Creating new MCP config...
        codewhale mcp add "hypernexus" --command "C:\Users\hyper\workspace\hypernexus\hypernexus.exe" --arg "mcp"
    )
    if %errorlevel%==0 (echo ✅ MCP server configured) else (echo ⚠️ MCP config may need manual setup)

    rem --- Install CodeWhale skill ---
    echo.
    echo Installing CodeWhale skill...
    if not exist "%USERPROFILE%\.codewhale\skills\hypernexus" mkdir "%USERPROFILE%\.codewhale\skills\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.codewhale\skills\hypernexus\SKILL.md" "%USERPROFILE%\.codewhale\skills\hypernexus\SKILL.md"
    if %errorlevel%==0 (echo ✅ CodeWhale skill installed) else (echo ⚠️ Skill copy failed)
) else (
    echo CodeWhale not found — skipping CodeWhale integration.
    echo Install CodeWhale with: npm install -g codewhale
    echo Then re-run this installer.
)
echo.

echo === Step 4: Starting Services ===
echo.
sc start HyperNexusSidecar
sc start HyperNexusDashboard
sc start HyperNexusWatchdog
echo.
echo ========================================
echo  Done!
echo.
echo  HyperNexus Pi extension:    ~\.pi\agent\extensions\hypernexus.ts
echo  CodeWhale skill:              ~\.codewhale\skills\hypernexus\SKILL.md
echo  CodeWhale MCP config:         ~\.codewhale\mcp.json
echo ========================================
pause
