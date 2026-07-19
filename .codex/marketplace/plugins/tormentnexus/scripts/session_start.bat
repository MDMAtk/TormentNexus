@echo off
setlocal enabledelayedexpansion
:: SessionStart hook — Log session to TN L2 memory
:: Reads JSON from stdin, sends session info to TN API

set /p STDIN=
if "!STDIN!"=="" exit /b 0

:: Extract session_id using PowerShell for JSON parsing
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.session_id"`) do set SESSION_ID=%%a

:: Log session start to TN L2
powershell -Command "try { $body=@{content=('{\"content\":\"Session started: ' + $env:SESSION_ID + '\",\"tags\":[\"system:session\",\"reason:codex\"],\"category\":\"session\",\"timestamp\":\"' + (Get-Date -Format o) + '\"}' -replace '\"','\\\"')}; Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/add' -Method Post -Body ($body|ConvertTo-Json) -ContentType 'application/json' -TimeoutSec 3 } catch {}"

echo {"continue": true}
exit /b 0
