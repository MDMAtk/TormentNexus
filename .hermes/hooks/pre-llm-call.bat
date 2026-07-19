@echo off
setlocal enabledelayedexpansion
:: pre_llm_call hook — Inject relevant context from TN L2
:: stdin: {"hook_event_name":"pre_llm_call","messages":[...]}
:: Return {} or {"context":"...recent memories..."}
set /p STDIN=
if "!STDIN!"=="" exit /b 0
:: Extract the last user message for semantic search
for /f "usebackq delims=" %%a in (`powershell -Command "$json='!STDIN!'|ConvertFrom-Json; $msgs=$json.messages; $last=$msgs[-1]; if($last -and $last.content){ $q=[System.Net.WebUtility]::UrlEncode($last.content.Substring(0,[Math]::Min(100,$last.content.Length))); try { $r=Invoke-RestMethod -Uri \"http://127.0.0.1:7778/api/memory/search?q=$q\" -TimeoutSec 3; if($r.data -and $r.data.Count -gt 0){ $ctx=''; $i=0; foreach($m in $r.data){ if($i -ge 3){break}; if($m.text){$ctx+='- '+$m.text.Substring(0,[Math]::Min(200,$m.text.Length))+'`n'}; $i++ }; if($ctx){ Write-Host ('{\"context\":\"TormentNexus relevant context:\n' + $ctx + '\"}') } else { Write-Host '{}' } } else { Write-Host '{}' } } catch { Write-Host '{}' } } else { Write-Host '{}' }"`)
if defined %%a echo %%a
exit /b 0
