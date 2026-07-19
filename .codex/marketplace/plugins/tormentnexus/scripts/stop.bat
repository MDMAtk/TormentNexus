@echo off
setlocal enabledelayedexpansion
:: Stop hook — Log session end to TN

set /p STDIN=
if "!STDIN!"=="" exit /b 0

for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.session_id"`) do set SESSION_ID=%%a

powershell -Command "try { $body=@{content=('{\"content\":\"Session ended: ' + $env:SESSION_ID + '\",\"tags\":[\"system:session_end\"],\"category\":\"session\",\"timestamp\":\"' + (Get-Date -Format o) + '\"}' -replace '\"','\\\"')}; Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/add' -Method Post -Body ($body|ConvertTo-Json) -ContentType 'application/json' -TimeoutSec 3 } catch {}"

echo {"continue": true}
exit /b 0
