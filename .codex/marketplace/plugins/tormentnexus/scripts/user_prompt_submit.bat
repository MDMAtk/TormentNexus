@echo off
setlocal enabledelayedexpansion
:: UserPromptSubmit hook — Inject TN context and handle @memory:key expansion
:: Reads JSON from stdin, optionally enriches with TN L2 context

set /p STDIN=
if "!STDIN!"=="" exit /b 0

:: Extract prompt text
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; echo $json.prompt"`) do set PROMPT=%%a

:: Check for @memory:key expansion
echo !PROMPT! | findstr /C:"@memory:" >nul
if !errorlevel!==0 (
    for /f "tokens=2 delims=:" %%k in ('echo !PROMPT! ^| findstr /R "@memory:[a-zA-Z0-9_-]*"') do (
        set MEMORY_KEY=%%k
        powershell -Command "try { $r=Invoke-RestMethod -Uri 'http://127.0.0.1:7778/api/memory/search?q=!MEMORY_KEY!' -TimeoutSec 3; if($r.data -and $r.data.Count -gt 0 -and $r.data[0].text){ $val=$r.data[0].text; $expanded='!PROMPT!'.replace('@memory:!MEMORY_KEY!',$val); echo ('{\"continue\":true,\"hookSpecificOutput\":{\"hookEventName\":\"UserPromptSubmit\",\"additionalContext\":\"' + $expanded + '\"}}') } else { echo '{\"continue\":true}' } } catch { echo '{\"continue\":true}' }"
        exit /b 0
    )
) else (
    :: Inject relevant L2 context for non-@memory prompts
    for /f "usebackq delims=" %%a in (`powershell -Command "try { $q=[System.Net.WebUtility]::UrlEncode('!PROMPT!'.Substring(0,[Math]::Min(100,'!PROMPT!'.Length))); $r=Invoke-RestMethod -Uri \"http://127.0.0.1:7778/api/memory/search?q=$q\" -TimeoutSec 3; if($r.data -and $r.data.Count -gt 0){ $ctx=''; $i=0; foreach($m in $r.data){ if($i -ge 3){break}; if($m.text){$ctx+='  - '+$m.text.Substring(0,[Math]::Min(200,$m.text.Length))+'`n'}; $i++ }; if($ctx){ echo ('{\"continue\":true,\"hookSpecificOutput\":{\"hookEventName\":\"UserPromptSubmit\",\"additionalContext\":\"Relevant TormentNexus context:\n' + $ctx + '\"}}') } else { echo '{\"continue\":true}' } } else { echo '{\"continue\":true}' } } catch { echo '{\"continue\":true}' }"`)
    )
    if defined %%a echo %%a
)

if not defined PROMPT echo {"continue": true}
exit /b 0
