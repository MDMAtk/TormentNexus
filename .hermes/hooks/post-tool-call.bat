@echo off
setlocal enabledelayedexpansion
:: post_tool_call hook — Auto-store tool results under 2KB
:: stdin: {"hook_event_name":"post_tool_call","tool_name":"...","tool_input":{...},"action":"...","result":"...","session_id":"..."}
set /p STDIN=
if "!STDIN!"=="" exit /b 0
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.tool_name"`) do set TNAME=%%a
if /i "!TNAME!"=="read" exit /b 0
if /i "!TNAME!"=="ls" exit /b 0
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; $r=$json.result; if($r -and $r.Length -lt 2000){ echo ($r|ConvertTo-Json -Compress) } else { echo '--SKIP--' }"`) do set RESULT=%%a
if "!RESULT!"=="--SKIP--" exit /b 0
powershell -Command "try { $body=@{content=('{\"content\":\"Hermes ' + $env:TNAME + ' completed\",\"tags\":[\"system:tool_result\",\"tool:' + $env:TNAME + '\",\"agent:hermes\"],\"category\":\"tool_result\",\"timestamp\":\"' + (Get-Date -Format o) + '\"}' -replace '\"','\\\"')}; Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/add' -Method Post -Body ($body|ConvertTo-Json) -ContentType 'application/json' -TimeoutSec 3 } catch {}"
echo {}
exit /b 0
