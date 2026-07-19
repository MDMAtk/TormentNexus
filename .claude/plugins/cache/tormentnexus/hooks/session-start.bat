@echo off
setlocal enabledelayedexpansion
:: SessionStart hook — Log session to TN L2 memory
:: stdin: {"session_id":"...","reason":"startup|resume|..."}

set /p STDIN=
if "!STDIN!"=="" exit /b 0

for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.session_id"`) do set SESSION_ID=%%a
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.reason"`) do set REASON=%%a

powershell -Command "try { $body=@{content=('{\"content\":\"Claude session ' + $env:REASON + ': ' + $env:SESSION_ID + '\",\"tags\":[\"system:session\",\"agent:claude\"],\"category\":\"session\",\"timestamp\":\"' + (Get-Date -Format o) + '\"}' -replace '\"','\\\"')}; Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/add' -Method Post -Body ($body|ConvertTo-Json) -ContentType 'application/json' -TimeoutSec 3 } catch {}"

echo {}
exit /b 0
