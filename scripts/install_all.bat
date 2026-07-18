@echo off
cd /d "C:\Users\hyper\workspace\hypernexus"
setlocal enabledelayedexpansion

echo ========================================
echo  HyperNexus Multi-Agent Installer
echo  Run this as Administrator for services!
echo ========================================
echo.

echo === Step 1: Windows Services ===
echo.
echo Registering Go Sidecar (port 7778)...
sc create "HyperNexusSidecar" binPath="\"C:\Users\hyper\workspace\hypernexus\hypernexus.exe\" serve" start=auto displayname="HyperNexus Sidecar"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)

echo Registering Dashboard (port 7779)...
sc create "HyperNexusDashboard" binPath="\"C:\Program Files\nodejs\node.exe\" \"C:\Users\hyper\workspace\hypernexus\apps\web\node_modules\.bin\next.cmd\" dev -p 7779" start=auto displayname="HyperNexus Dashboard"

echo Registering Watchdog...
sc create "HyperNexusWatchdog" binPath="\"C:\Python314\pythonw.exe\" -u \"C:\Users\hyper\workspace\hypernexus\watchdog.py\"" start=auto displayname="HyperNexus Watchdog"
echo.

echo === Step 2: Pi Coding Agent ===
echo.
if not exist "%USERPROFILE%\.pi\agent\extensions" mkdir "%USERPROFILE%\.pi\agent\extensions"
if exist "%USERPROFILE%\.pi\agent\extensions\hypernexus.ts" (
    echo Pi extension already exists. Skipping.
) else (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.pi\extensions\hypernexus.ts" "%USERPROFILE%\.pi\agent\extensions\hypernexus.ts"
    if !errorlevel!==0 (echo ✅ Pi extension) else (echo ⚠️)
)
echo.

echo === Step 3: Ollama / vLLM (Tool Prediction Engine) ===
echo.
echo HyperNexus uses a local LLM for tool prediction (ConversationalToolInjector).
echo.
echo Choose an option:
echo   [1] Ollama — easiest, auto-start as Windows service (recommended)
echo   [2] vLLM  — faster inference, GPU-accelerated
echo   [S] Skip — tool prediction degrades to keyword matching
echo.
choice /C 12S /N /M "Select [1], [2], or [S]: "
if errorlevel 3 goto :skip_llm
if errorlevel 2 goto :install_vllm
if errorlevel 1 goto :install_ollama

:install_ollama
echo Installing Ollama...
curl -sL -o "%TEMP%\ollama_windows.exe" "https://github.com/ollama/ollama/releases/latest/download/OllamaSetup.exe"
if exist "%TEMP%\ollama_windows.exe" (
    start /wait "" "%TEMP%\ollama_windows.exe" /S
    echo Installing Gemma 4 model (this downloads ~8GB, may take a while)...
    ollama pull gemma4 2>nul || ollama pull gemma3:12b 2>nul
    echo.
    echo Setting up Ollama as auto-start service...
    sc config ollama start=auto >nul 2>nul
    sc start ollama >nul 2>nul
    echo ✅ Ollama + Gemma 4 installed at http://127.0.0.1:11434
) else (
    echo ⚠️ Download failed. Install manually from https://ollama.ai
)
goto :end_llm

:install_vllm
echo vLLM installation requires Python + CUDA.
echo.
echo pip install vllm
echo vllm serve gemma-4 --port 11434 --api-key token-abc123
echo.
echo Set HYPERNEXUS_OLLAMA_URL=http://127.0.0.1:11434
echo.
echo Manual setup required — see https://github.com/vllm-project/vllm
goto :end_llm

:skip_llm
echo Skipping LLM install. Tool prediction will use BM25 keyword matching.
goto :end_llm

:end_llm
echo.

echo === Step 4: CodeWhale Integration ===
echo.
call "C:\Users\hyper\workspace\hypernexus\scripts\install_codewhale.bat"
echo.

echo === Step 5: Gemini CLI ===
echo.
where gemini >nul 2>nul
if %errorlevel%==0 (
    if not exist "%USERPROFILE%\.gemini\extensions" mkdir "%USERPROFILE%\.gemini\extensions"
    xcopy /E /I /Y "C:\Users\hyper\workspace\hypernexus\.gemini\extensions\hypernexus" "%USERPROFILE%\.gemini\extensions\hypernexus" >nul
    gemini extensions link "%USERPROFILE%\.gemini\extensions\hypernexus" >nul 2>nul
    if not exist "%USERPROFILE%\.gemini\skills\hypernexus" mkdir "%USERPROFILE%\.gemini\skills\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.gemini\skills\hypernexus\SKILL.md" "%USERPROFILE%\.gemini\skills\hypernexus\SKILL.md" >nul
    echo ✅ Gemini CLI extension + skill
) else (echo ⏭️)
echo.

echo === Step 6: Claude Desktop ===
echo.
if exist "%APPDATA%\Claude\claude_desktop_config.json" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.editor-configs\claude-desktop-mcp.json" "%APPDATA%\Claude\claude_desktop_config.json.tn-template" >nul
    echo ✅ Claude Desktop template saved
)
echo.

echo === Step 7: Claude Code CLI ===
echo.
where claude >nul 2>nul
if %errorlevel%==0 (
    claude mcp add --transport stdio hypernexus -- "C:\Users\hyper\workspace\hypernexus\hypernexus.exe" "mcp" >nul 2>nul
    if !errorlevel!==0 (echo ✅ Claude CLI MCP) else (echo May already exist)
) else (echo ⏭️)
echo.

echo === Step 8: Codex CLI ===
echo.
where codex >nul 2>nul
if %errorlevel%==0 (
    codex mcp add "hypernexus" --env HYPERNEXUS_WORKSPACE_ROOT="C:\Users\hyper\workspace\hypernexus" -- "C:\Users\hyper\workspace\hypernexus\hypernexus.exe" "mcp" >nul 2>nul
    codex plugin marketplace add "C:\Users\hyper\workspace\hypernexus\.codex\marketplace" >nul 2>nul
    codex plugin add hypernexus@hypernexus-marketplace >nul 2>nul
    if not exist "%USERPROFILE%\.codex\skills\hypernexus" mkdir "%USERPROFILE%\.codex\skills\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.codex\skills\hypernexus\SKILL.md" "%USERPROFILE%\.codex\skills\hypernexus\SKILL.md" >nul
    echo ✅ Codex CLI plugin + MCP + skill
) else (echo ⏭️)
echo.

echo === Step 9: Cursor Extension + MCP ===
echo.
if exist "%USERPROFILE%\.cursor\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.editor-configs\cursor-mcp.json" "%USERPROFILE%\.cursor\mcp.json.tn-template" >nul
)
if not exist "%USERPROFILE%\.cursor\extensions\hypernexus" mkdir "%USERPROFILE%\.cursor\extensions\hypernexus"
copy /Y "C:\Users\hyper\workspace\hypernexus\.cursor\extensions\hypernexus\package.json" "%USERPROFILE%\.cursor\extensions\hypernexus\package.json" >nul
copy /Y "C:\Users\hyper\workspace\hypernexus\.cursor\extensions\hypernexus\extension.js" "%USERPROFILE%\.cursor\extensions\hypernexus\extension.js" >nul
if not exist "%USERPROFILE%\.cursor\rules" mkdir "%USERPROFILE%\.cursor\rules"
copy /Y "C:\Users\hyper\workspace\hypernexus\.cursor\rules\hypernexus.mdc" "%USERPROFILE%\.cursor\rules\hypernexus.mdc" >nul
if not exist "%USERPROFILE%\.cursor\commands" mkdir "%USERPROFILE%\.cursor\commands"
copy /Y "C:\Users\hyper\workspace\hypernexus\.cursor\commands\tn-store.md" "%USERPROFILE%\.cursor\commands\tn-store.md" >nul
copy /Y "C:\Users\hyper\workspace\hypernexus\.cursor\commands\tn-search.md" "%USERPROFILE%\.cursor\commands\tn-search.md" >nul
copy /Y "C:\Users\hyper\workspace\hypernexus\.cursor\commands\tn-status.md" "%USERPROFILE%\.cursor\commands\tn-status.md" >nul
echo ✅ Cursor: extension + MCP + rules + commands
echo.

echo === Step 10: Windsurf ===
echo.
where windsurf >nul 2>nul
if %errorlevel%==0 (
    windsurf --add-mcp "{"""name""":"""hypernexus""","""command""":"""C:\\Users\\hyper\\workspace\\hypernexus\\hypernexus.exe""","""args""":["""mcp"""]}" >nul 2>nul
    if !errorlevel!==0 (echo ✅ Windsurf MCP) else (echo ⚠️)
) else (echo ⏭️)
echo.

echo === Step 11: VS Code Extension + MCP ===
echo.
if exist "%USERPROFILE%\.vscode\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.editor-configs\vscode-mcp.json" "%USERPROFILE%\.vscode\mcp.json.tn-template" >nul
) else (
    mkdir "%USERPROFILE%\.vscode" >nul 2>nul
    copy /Y "C:\Users\hyper\workspace\hypernexus\.editor-configs\vscode-mcp.json" "%USERPROFILE%\.vscode\mcp.json" >nul
)
if not exist "%USERPROFILE%\.vscode\extensions\hypernexus" mkdir "%USERPROFILE%\.vscode\extensions\hypernexus"
copy /Y "C:\Users\hyper\workspace\hypernexus\.vscode\extensions\hypernexus\package.json" "%USERPROFILE%\.vscode\extensions\hypernexus\package.json" >nul
copy /Y "C:\Users\hyper\workspace\hypernexus\.vscode\extensions\hypernexus\extension.js" "%USERPROFILE%\.vscode\extensions\hypernexus\extension.js" >nul
echo ✅ VS Code extension + MCP
echo.

echo === Step 12: Copilot CLI ===
echo.
if exist "%USERPROFILE%\.copilot\mcp-config.json" (
    echo Copilot CLI detected — installing extension + MCP...
    if not exist "%USERPROFILE%\.copilot\extensions\hypernexus" mkdir "%USERPROFILE%\.copilot\extensions\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.copilot\extensions\hypernexus\extension.mjs" "%USERPROFILE%\.copilot\extensions\hypernexus\extension.mjs" >nul
    if %errorlevel%==0 (echo ✅ Copilot extension: 5 hooks + 2 tools) else (echo ⚠️)
    copy /Y "C:\Users\hyper\workspace\hypernexus\.editor-configs\copilot-mcp.json" "%USERPROFILE%\.copilot\hypernexus-mcp-merge.json.tn" >nul
) else (echo ⏭️)
echo.

echo === Step 13: OpenCode ===
echo.
if exist "%USERPROFILE%\.opencode\mcp.json" (echo ✅ Already configured) else (echo ⏭️)
echo.

echo === Step 14: Continue ===
echo.
if exist "%USERPROFILE%\.continue\config.json" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.editor-configs\continue-mcp.json" "%USERPROFILE%\.continue\hypernexus-mcp-merge.json.tn" >nul
    echo ✅ Template saved
) else (echo ⏭️)
echo.

echo === Step 15: Mavis / MiniMax Code ===
echo.
if exist "%USERPROFILE%\.mavis\mcp\mcp.json" (
    if not exist "%USERPROFILE%\.mavis\skills\hypernexus" mkdir "%USERPROFILE%\.mavis\skills\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.mavis\mcp.json" "%USERPROFILE%\.mavis\mcp\hypernexus-mcp-merge.json.tn" >nul
    copy /Y "C:\Users\hyper\workspace\hypernexus\.mavis\skills\hypernexus\SKILL.md" "%USERPROFILE%\.mavis\skills\hypernexus\SKILL.md" >nul
    echo ✅ Mavis MCP + skill
) else (echo ⏭️)
echo.

echo === Step 16: Antigravity IDE ===
echo.
if exist "%USERPROFILE%\.gemini\antigravity\mcp" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.gemini\antigravity\mcp\mcp.json" "%USERPROFILE%\.gemini\antigravity\mcp\hypernexus-mcp-merge.json.tn" >nul
    if not exist "%USERPROFILE%\.gemini\antigravity-ide\extensions\hypernexus" mkdir "%USERPROFILE%\.gemini\antigravity-ide\extensions\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.gemini\antigravity-ide\extensions\hypernexus\SKILL.md" "%USERPROFILE%\.gemini\antigravity-ide\extensions\hypernexus\SKILL.md" >nul
    copy /Y "C:\Users\hyper\workspace\hypernexus\.gemini\antigravity-ide\extensions\hypernexus\agent.md" "%USERPROFILE%\.gemini\antigravity-ide\extensions\hypernexus\agent.md" >nul
    echo ✅ Antigravity 2.0 ADE
) else if exist "%USERPROFILE%\.antigravity" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.antigravity\mcp_config.json" "%USERPROFILE%\.antigravity\hypernexus-mcp-merge.json.tn" >nul
    if not exist "%USERPROFILE%\.antigravity\extensions\hypernexus" mkdir "%USERPROFILE%\.antigravity\extensions\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.antigravity\extensions\hypernexus\SKILL.md" "%USERPROFILE%\.antigravity\extensions\hypernexus\SKILL.md" >nul
    if not exist "%USERPROFILE%\.antigravity\agents\hypernexus" mkdir "%USERPROFILE%\.antigravity\agents\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.antigravity\agents\hypernexus\agent.md" "%USERPROFILE%\.antigravity\agents\hypernexus\agent.md" >nul
    echo ✅ Antigravity 1.0 IDE
) else (echo ⏭️)
echo.

echo === Step 17: Kimi Desktop ===
echo.
if exist "%USERPROFILE%\.kimi-code" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.kimi-code\mcp.json" "%USERPROFILE%\.kimi-code\mcp.json.tn-template" >nul
    echo ✅ Template saved
) else (echo ⏭️)
echo.

echo === Step 18: ZCode Desktop ===
echo.
if exist "%USERPROFILE%\.zcode" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.zcode\mcp.json" "%USERPROFILE%\.zcode\mcp.json.tn-template" >nul
    echo ✅ Template saved
) else (echo ⏭️)
echo.

echo === Step 19: Hermes Agent ===
echo.
if exist "%USERPROFILE%\.hermes\config.yaml" (
    if not exist "%USERPROFILE%\.hermes\optional-mcps\hypernexus" mkdir "%USERPROFILE%\.hermes\optional-mcps\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.hermes\optional-mcps\hypernexus\manifest.yaml" "%USERPROFILE%\.hermes\optional-mcps\hypernexus\manifest.yaml" >nul
    if not exist "%USERPROFILE%\.hermes\skills\hypernexus" mkdir "%USERPROFILE%\.hermes\skills\hypernexus"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.hermes\skills\hypernexus\DESCRIPTION.md" "%USERPROFILE%\.hermes\skills\hypernexus\DESCRIPTION.md" >nul
    if not exist "%USERPROFILE%\.hermes\hooks" mkdir "%USERPROFILE%\.hermes\hooks"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.hermes\hooks\*.bat" "%USERPROFILE%\.hermes\hooks\" >nul
    copy /Y "C:\Users\hyper\workspace\hypernexus\.hermes\hooks-config.yaml" "%USERPROFILE%\.hermes\hypernexus-hooks-merge.yaml.tn" >nul
    echo ✅ Hermes MCP + 5 hooks + skill
) else (echo ⏭️)
echo.

echo === Step 20: Aider ===
echo.
if exist "%USERPROFILE%\.aider.conf.yml" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.aider\mcp.json" "%USERPROFILE%\.aider.mcp.json.tn-template" >nul
    echo ✅ Aider MCP template
) else (echo ⏭️)
echo.

echo === Step 21: Cline ===
echo.
if exist "%USERPROFILE%\.cline" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.cline\mcp.json" "%USERPROFILE%\.cline\mcp.json.tn-template" >nul
    echo ✅ Cline MCP template
) else (echo ⏭️)
echo.

echo === Step 22: Roo Code ===
echo.
if exist "%USERPROFILE%\.roo" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.roo\mcp.json" "%USERPROFILE%\.roo\mcp.json.tn-template" >nul
    echo ✅ Roo Code MCP template
) else (echo ⏭️)
echo.

echo === Step 23: Kilo Code ===
echo.
if exist "%USERPROFILE%\.kilo" (
    copy /Y "C:\Users\hyper\workspace\hypernexus\.kilo\mcp.json" "%USERPROFILE%\.kilo\mcp.json.tn-template" >nul
    echo ✅ Kilo Code MCP template
) else (echo ⏭️)
echo.

echo === Step 24: OpenHands ===
echo.
if exist "%USERPROFILE%\.openhands" (
    if not exist "%USERPROFILE%\.openhands\microagents" mkdir "%USERPROFILE%\.openhands\microagents"
    copy /Y "C:\Users\hyper\workspace\hypernexus\.openhands\microagents\hypernexus.md" "%USERPROFILE%\.openhands\microagents\hypernexus.md" >nul
    echo ✅ OpenHands micro-agent installed
) else (echo ⏭️)
echo.

echo === Step 25: Goose ===
echo.
if exist "%USERPROFILE%\.goose" (
    if not exist "%USERPROFILE%\.goose\extensions\hypernexus" mkdir "%USERPROFILE%\.goose\extensions\hypernexus"
    echo ✅ Goose extensions directory ready
) else (echo ⏭️)
echo.

echo === Step 26: Starting Services ===
echo.
sc start HyperNexusSidecar >nul 2>nul
sc start HyperNexusDashboard >nul 2>nul
sc start HyperNexusWatchdog >nul 2>nul
echo.
echo ========================================
echo  HyperNexus Multi-Agent Installer
echo  Complete!
echo.
echo  Installed for:
echo   ✅ Pi Coding Agent        ~\.pi\agent\extensions\
echo   ✅ Ollama/vLLM            Tool prediction engine
echo   ✅ CodeWhale              ~\.codewhale\plugins\ + MCP
echo   ✅ Gemini CLI             ~\.gemini\extensions\ + skills
echo   ✅ Claude Desktop         template saved
echo   ✅ Claude Code CLI        MCP configured
echo   ✅ Codex CLI              plugin + MCP + skill
echo   ✅ Cursor                 extension + MCP
echo   ✅ Windsurf               MCP added
echo   ✅ VS Code                extension + MCP
echo   ✅ Copilot CLI            extension + MCP
echo   ✅ Mavis / MiniMax Code   .mavis\skills\ + MCP
echo   ✅ Antigravity IDE        MCP + extension + agent
echo   ✅ Kimi Desktop           MCP template
echo   ✅ ZCode Desktop          MCP template
echo   ✅ Hermes Agent           MCP + 5 hooks + skill
echo ========================================
echo.
pause
