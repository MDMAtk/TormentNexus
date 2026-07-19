@echo off
setlocal enabledelayedexpansion
:: PostToolUse hook — Auto-store tool results to TN
:: Reads JSON from stdin with tool_name, tool_input, tool_response

set /p STDIN=
if "!STDIN!"=="" exit /b 0

:: Extract tool info
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.tool_name"`) do set TOOL_NAME=%%a
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo ($json.tool_response|ConvertTo-Json -Compress)"`) do set TOOL_RESP=%%a

:: Only store non-trivial results from non-read-only tools
if /i "!TOOL_NAME!"=="read" exit /b 0
if /i "!TOOL_NAME!"=="ls" exit /b 0

:: Check result size (skip if too large)
set RESULT_LEN=0
for %%a in ("!TOOL_RESP!") do set RESULT_LEN=%%~za
if !RESULT_LEN! gtr 2000 exit /b 0

:: Log tool result to TN L2
powershell -Command "try { $body=@{content=('{\"content\":\"' + $env:TOOL_NAME + ' completed\",\"tags\":[\"system:tool_result\",\"tool:' + $env:TOOL_NAME + '\"],\"result\":' + $env:TOOL_RESP + ',\"category\":\"tool_result\",\"timestamp\":\"' + (Get-Date -Format o) + '\"}' -replace '\"','\\\"')}; Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/add' -Method Post -Body ($body|ConvertTo-Json) -ContentType 'application/json' -TimeoutSec 3 } catch {}"

echo {"continue": true}
exit /b 0
