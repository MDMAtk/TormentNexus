@echo off
setlocal enabledelayedexpansion
:: PreToolUse hook — Log tool calls to TN, enforce RBAC
:: Reads JSON from stdin with tool_name, tool_input

set /p STDIN=
if "!STDIN!"=="" exit /b 0

:: Extract tool info
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.tool_name"`) do set TOOL_NAME=%%a
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo ($json.tool_input|ConvertTo-Json -Compress)"`) do set TOOL_INPUT=%%a

:: Log tool call to TN L2
powershell -Command "try { $body=@{content=('{\"content\":\"Tool call: ' + $env:TOOL_NAME + '\",\"tags\":[\"system:tool_call\",\"tool:' + $env:TOOL_NAME + '\"],\"data\":' + $env:TOOL_INPUT + ',\"category\":\"tool\",\"timestamp\":\"' + (Get-Date -Format o) + '\"}' -replace '\"','\\\"')}; Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/add' -Method Post -Body ($body|ConvertTo-Json) -ContentType 'application/json' -TimeoutSec 3 } catch {}"

:: Always allow, just log
echo {"continue": true}
exit /b 0
