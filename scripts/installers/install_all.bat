@echo off
cd /d "C:\Users\hyper\workspace\tormentnexus"
setlocal enabledelayedexpansion

echo ========================================
echo  TormentNexus Multi-Agent Installer
echo  Run this as Administrator for services!
echo ========================================
echo.

echo === Step 1: Windows Services ===
echo.
echo Registering Go Sidecar (port 7778)...
sc create "TormentNexusSidecar" binPath="\"C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe\" serve" start=auto displayname="TormentNexus Sidecar"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)

echo Registering Dashboard (port 7779)...
sc create "TormentNexusDashboard" binPath="\"C:\Program Files\nodejs\node.exe\" \"C:\Users\hyper\workspace\tormentnexus\apps\web\node_modules\.bin\next.cmd\" dev -p 7779" start=auto displayname="TormentNexus Dashboard"

echo Registering Watchdog...
sc create "TormentNexusWatchdog" binPath="\"C:\Python314\pythonw.exe\" -u \"C:\Users\hyper\workspace\tormentnexus\watchdog.py\"" start=auto displayname="TormentNexus Watchdog"
echo.

echo === Step 2: Pi Coding Agent ===
echo.
if not exist "%USERPROFILE%\.pi\agent\extensions" mkdir "%USERPROFILE%\.pi\agent\extensions"
if exist "%USERPROFILE%\.pi\agent\extensions\tormentnexus.ts" (
    echo Pi extension already exists. Skipping.
) else (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.pi\extensions\tormentnexus.ts" "%USERPROFILE%\.pi\agent\extensions\tormentnexus.ts"
    if !errorlevel!==0 (echo ✅ Pi extension) else (echo ⚠️)
)
echo.

echo === Step 3: Ollama / vLLM (Tool Prediction Engine) ===
echo.
echo TormentNexus uses a local LLM for tool prediction (ConversationalToolInjector).
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
echo Set TORMENTNEXUS_OLLAMA_URL=http://127.0.0.1:11434
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
call "C:\Users\hyper\workspace\tormentnexus\scripts\install_codewhale.bat"
echo.

echo === Step 5: Gemini CLI ===
echo.
where gemini >nul 2>nul
if %errorlevel%==0 (
    if not exist "%USERPROFILE%\.gemini\extensions" mkdir "%USERPROFILE%\.gemini\extensions"
    xcopy /E /I /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\extensions\tormentnexus" "%USERPROFILE%\.gemini\extensions\tormentnexus" >nul
    gemini extensions link "%USERPROFILE%\.gemini\extensions\tormentnexus" >nul 2>nul
    if not exist "%USERPROFILE%\.gemini\skills\tormentnexus" mkdir "%USERPROFILE%\.gemini\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.gemini\skills\tormentnexus\SKILL.md" >nul
    echo ✅ Gemini CLI extension + skill
) else (echo ⏭️)
echo.

echo === Step 6: Claude Desktop ===
echo.
if exist "%APPDATA%\Claude\claude_desktop_config.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\claude-desktop-mcp.json" "%APPDATA%\Claude\claude_desktop_config.json.tn-template" >nul
    echo ✅ Claude Desktop template saved
)
echo.

echo === Step 7: Claude Code CLI ===
echo.
where claude >nul 2>nul
if %errorlevel%==0 (
    claude mcp add --transport stdio tormentnexus -- "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" "mcp" >nul 2>nul
    if !errorlevel!==0 (echo ✅ Claude CLI MCP) else (echo May already exist)
) else (echo ⏭️)
echo.

echo === Step 8: Codex CLI ===
echo.
where codex >nul 2>nul
if %errorlevel%==0 (
    codex mcp add "tormentnexus" --env TORMENTNEXUS_WORKSPACE_ROOT="C:\Users\hyper\workspace\tormentnexus" -- "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" "mcp" >nul 2>nul
    codex plugin marketplace add "C:\Users\hyper\workspace\tormentnexus\.codex\marketplace" >nul 2>nul
    codex plugin add tormentnexus@tormentnexus-marketplace >nul 2>nul
    if not exist "%USERPROFILE%\.codex\skills\tormentnexus" mkdir "%USERPROFILE%\.codex\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.codex\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.codex\skills\tormentnexus\SKILL.md" >nul
    echo ✅ Codex CLI plugin + MCP + skill
) else (echo ⏭️)
echo.

echo === Step 9: Cursor Extension + MCP ===
echo.
if exist "%USERPROFILE%\.cursor\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\cursor-mcp.json" "%USERPROFILE%\.cursor\mcp.json.tn-template" >nul
)
if not exist "%USERPROFILE%\.cursor\extensions\tormentnexus" mkdir "%USERPROFILE%\.cursor\extensions\tormentnexus"
copy /Y "C:\Users\hyper\workspace\tormentnexus\.cursor\extensions\tormentnexus\package.json" "%USERPROFILE%\.cursor\extensions\tormentnexus\package.json" >nul
copy /Y "C:\Users\hyper\workspace\tormentnexus\.cursor\extensions\tormentnexus\extension.js" "%USERPROFILE%\.cursor\extensions\tormentnexus\extension.js" >nul
if not exist "%USERPROFILE%\.cursor\rules" mkdir "%USERPROFILE%\.cursor\rules"
copy /Y "C:\Users\hyper\workspace\tormentnexus\.cursor\rules\tormentnexus.mdc" "%USERPROFILE%\.cursor\rules\tormentnexus.mdc" >nul
if not exist "%USERPROFILE%\.cursor\commands" mkdir "%USERPROFILE%\.cursor\commands"
copy /Y "C:\Users\hyper\workspace\tormentnexus\.cursor\commands\tn-store.md" "%USERPROFILE%\.cursor\commands\tn-store.md" >nul
copy /Y "C:\Users\hyper\workspace\tormentnexus\.cursor\commands\tn-search.md" "%USERPROFILE%\.cursor\commands\tn-search.md" >nul
copy /Y "C:\Users\hyper\workspace\tormentnexus\.cursor\commands\tn-status.md" "%USERPROFILE%\.cursor\commands\tn-status.md" >nul
echo ✅ Cursor: extension + MCP + rules + commands
echo.

echo === Step 10: Windsurf ===
echo.
where windsurf >nul 2>nul
if %errorlevel%==0 (
    windsurf --add-mcp "{"""name""":"""tormentnexus""","""command""":"""C:\\Users\\hyper\\workspace\\tormentnexus\\tormentnexus.exe""","""args""":["""mcp"""]}" >nul 2>nul
    if !errorlevel!==0 (echo ✅ Windsurf MCP) else (echo ⚠️)
) else (echo ⏭️)
echo.

echo === Step 11: VS Code Extension + MCP ===
echo.
if exist "%USERPROFILE%\.vscode\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\vscode-mcp.json" "%USERPROFILE%\.vscode\mcp.json.tn-template" >nul
) else (
    mkdir "%USERPROFILE%\.vscode" >nul 2>nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\vscode-mcp.json" "%USERPROFILE%\.vscode\mcp.json" >nul
)
if not exist "%USERPROFILE%\.vscode\extensions\tormentnexus" mkdir "%USERPROFILE%\.vscode\extensions\tormentnexus"
copy /Y "C:\Users\hyper\workspace\tormentnexus\.vscode\extensions\tormentnexus\package.json" "%USERPROFILE%\.vscode\extensions\tormentnexus\package.json" >nul
copy /Y "C:\Users\hyper\workspace\tormentnexus\.vscode\extensions\tormentnexus\extension.js" "%USERPROFILE%\.vscode\extensions\tormentnexus\extension.js" >nul
echo ✅ VS Code extension + MCP
echo.

echo === Step 12: Copilot CLI ===
echo.
if exist "%USERPROFILE%\.copilot\mcp-config.json" (
    echo Copilot CLI detected — installing extension + MCP...
    if not exist "%USERPROFILE%\.copilot\extensions\tormentnexus" mkdir "%USERPROFILE%\.copilot\extensions\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.copilot\extensions\tormentnexus\extension.mjs" "%USERPROFILE%\.copilot\extensions\tormentnexus\extension.mjs" >nul
    if %errorlevel%==0 (echo ✅ Copilot extension: 5 hooks + 2 tools) else (echo ⚠️)
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\copilot-mcp.json" "%USERPROFILE%\.copilot\tormentnexus-mcp-merge.json.tn" >nul
) else (echo ⏭️)
echo.

echo === Step 13: OpenCode ===
echo.
if exist "%USERPROFILE%\.opencode\mcp.json" (echo ✅ Already configured) else (echo ⏭️)
echo.

echo === Step 14: Continue ===
echo.
if exist "%USERPROFILE%\.continue\config.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\continue-mcp.json" "%USERPROFILE%\.continue\tormentnexus-mcp-merge.json.tn" >nul
    echo ✅ Template saved
) else (echo ⏭️)
echo.

echo === Step 15: Mavis / MiniMax Code ===
echo.
if exist "%USERPROFILE%\.mavis\mcp\mcp.json" (
    if not exist "%USERPROFILE%\.mavis\skills\tormentnexus" mkdir "%USERPROFILE%\.mavis\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.mavis\mcp.json" "%USERPROFILE%\.mavis\mcp\tormentnexus-mcp-merge.json.tn" >nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.mavis\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.mavis\skills\tormentnexus\SKILL.md" >nul
    echo ✅ Mavis MCP + skill
) else (echo ⏭️)
echo.

echo === Step 16: Antigravity IDE ===
echo.
if exist "%USERPROFILE%\.gemini\antigravity\mcp" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\antigravity\mcp\mcp.json" "%USERPROFILE%\.gemini\antigravity\mcp\tormentnexus-mcp-merge.json.tn" >nul
    if not exist "%USERPROFILE%\.gemini\antigravity-ide\extensions\tormentnexus" mkdir "%USERPROFILE%\.gemini\antigravity-ide\extensions\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\antigravity-ide\extensions\tormentnexus\SKILL.md" "%USERPROFILE%\.gemini\antigravity-ide\extensions\tormentnexus\SKILL.md" >nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\antigravity-ide\extensions\tormentnexus\agent.md" "%USERPROFILE%\.gemini\antigravity-ide\extensions\tormentnexus\agent.md" >nul
    echo ✅ Antigravity 2.0 ADE
) else if exist "%USERPROFILE%\.antigravity" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\mcp_config.json" "%USERPROFILE%\.antigravity\tormentnexus-mcp-merge.json.tn" >nul
    if not exist "%USERPROFILE%\.antigravity\extensions\tormentnexus" mkdir "%USERPROFILE%\.antigravity\extensions\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\extensions\tormentnexus\SKILL.md" "%USERPROFILE%\.antigravity\extensions\tormentnexus\SKILL.md" >nul
    if not exist "%USERPROFILE%\.antigravity\agents\tormentnexus" mkdir "%USERPROFILE%\.antigravity\agents\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\agents\tormentnexus\agent.md" "%USERPROFILE%\.antigravity\agents\tormentnexus\agent.md" >nul
    echo ✅ Antigravity 1.0 IDE
) else (echo ⏭️)
echo.

echo === Step 17: Kimi Desktop ===
echo.
if exist "%USERPROFILE%\.kimi-code" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.kimi-code\mcp.json" "%USERPROFILE%\.kimi-code\mcp.json.tn-template" >nul
    echo ✅ Template saved
) else (echo ⏭️)
echo.

echo === Step 18: ZCode Desktop ===
echo.
if exist "%USERPROFILE%\.zcode" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.zcode\mcp.json" "%USERPROFILE%\.zcode\mcp.json.tn-template" >nul
    echo ✅ Template saved
) else (echo ⏭️)
echo.

echo === Step 19: Hermes Agent ===
echo.
if exist "%USERPROFILE%\.hermes\config.yaml" (
    if not exist "%USERPROFILE%\.hermes\optional-mcps\tormentnexus" mkdir "%USERPROFILE%\.hermes\optional-mcps\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.hermes\optional-mcps\tormentnexus\manifest.yaml" "%USERPROFILE%\.hermes\optional-mcps\tormentnexus\manifest.yaml" >nul
    if not exist "%USERPROFILE%\.hermes\skills\tormentnexus" mkdir "%USERPROFILE%\.hermes\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.hermes\skills\tormentnexus\DESCRIPTION.md" "%USERPROFILE%\.hermes\skills\tormentnexus\DESCRIPTION.md" >nul
    if not exist "%USERPROFILE%\.hermes\hooks" mkdir "%USERPROFILE%\.hermes\hooks"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.hermes\hooks\*.bat" "%USERPROFILE%\.hermes\hooks\" >nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.hermes\hooks-config.yaml" "%USERPROFILE%\.hermes\tormentnexus-hooks-merge.yaml.tn" >nul
    echo ✅ Hermes MCP + 5 hooks + skill
) else (echo ⏭️)
echo.

echo === Step 20: Aider ===
echo.
if exist "%USERPROFILE%\.aider.conf.yml" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.aider\mcp.json" "%USERPROFILE%\.aider.mcp.json.tn-template" >nul
    echo ✅ Aider MCP template
) else (echo ⏭️)
echo.

echo === Step 21: Cline ===
echo.
if exist "%USERPROFILE%\.cline" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.cline\mcp.json" "%USERPROFILE%\.cline\mcp.json.tn-template" >nul
    echo ✅ Cline MCP template
) else (echo ⏭️)
echo.

echo === Step 22: Roo Code ===
echo.
if exist "%USERPROFILE%\.roo" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.roo\mcp.json" "%USERPROFILE%\.roo\mcp.json.tn-template" >nul
    echo ✅ Roo Code MCP template
) else (echo ⏭️)
echo.

echo === Step 23: Kilo Code ===
echo.
if exist "%USERPROFILE%\.kilo" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.kilo\mcp.json" "%USERPROFILE%\.kilo\mcp.json.tn-template" >nul
    echo ✅ Kilo Code MCP template
) else (echo ⏭️)
echo.

echo === Step 24: OpenHands ===
echo.
if exist "%USERPROFILE%\.openhands" (
    if not exist "%USERPROFILE%\.openhands\microagents" mkdir "%USERPROFILE%\.openhands\microagents"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.openhands\microagents\tormentnexus.md" "%USERPROFILE%\.openhands\microagents\tormentnexus.md" >nul
    echo ✅ OpenHands micro-agent installed
) else (echo ⏭️)
echo.

echo === Step 25: Goose ===
echo.
if exist "%USERPROFILE%\.goose" (
    if not exist "%USERPROFILE%\.goose\extensions\tormentnexus" mkdir "%USERPROFILE%\.goose\extensions\tormentnexus"
    echo ✅ Goose extensions directory ready
) else (echo ⏭️)
echo.

echo === Step 26: Starting Services ===
echo.
sc start TormentNexusSidecar >nul 2>nul
sc start TormentNexusDashboard >nul 2>nul
sc start TormentNexusWatchdog >nul 2>nul
echo.
echo ========================================
echo  TormentNexus Multi-Agent Installer
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
