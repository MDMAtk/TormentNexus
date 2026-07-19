@echo off
setlocal enabledelayedexpansion
:: pre_tool_call hook — Log tool calls to TN L2, handle @memory:key
:: stdin: {"hook_event_name":"pre_tool_call","tool_name":"...","tool_input":{...}}
:: Return {"action":"allow"} or {"action":"block","message":"..."}

set /p STDIN=
if "!STDIN!"=="" exit /b 0
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.tool_name"`) do set TNAME=%%a
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; try { echo ($json.tool_input|ConvertTo-Json -Compress) } catch { echo '' }"`) do set TINPUT=%%a
powershell -Command "try { $body=@{content=('{\"content\":\"Hermes tool: ' + $env:TNAME + '\",\"tags\":[\"system:tool_call\",\"tool:' + $env:TNAME + '\",\"agent:hermes\"],\"category\":\"tool\",\"timestamp\":\"' + (Get-Date -Format o) + '\"}' -replace '\"','\\\"')}; Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/add' -Method Post -Body ($body|ConvertTo-Json) -ContentType 'application/json' -TimeoutSec 3 } catch {}"
echo {"action": "allow"}
exit /b 0
