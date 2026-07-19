@echo off
setlocal enabledelayedexpansion
:: on_session_end hook — Log session end to TN L2
set /p STDIN=
if "!STDIN!"=="" exit /b 0
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.session_id"`) do set SID=%%a
powershell -Command "try { $body=@{content=('{\"content\":\"Hermes session ended: ' + $env:SID + '\",\"tags\":[\"system:session_end\",\"agent:hermes\"],\"category\":\"session\",\"timestamp\":\"' + (Get-Date -Format o) + '\"}' -replace '\"','\\\"')}; Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/add' -Method Post -Body ($body|ConvertTo-Json) -ContentType 'application/json' -TimeoutSec 3 } catch {}"
echo {}
exit /b 0
