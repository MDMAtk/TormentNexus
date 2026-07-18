# HyperNexus Service Registration
# Run this script to register auto-start scheduled tasks

$action1 = New-ScheduledTaskAction -Execute "C:\Users\hyper\workspace\hypernexus\hypernexus.exe" -Argument "serve" -WorkingDirectory "C:\Users\hyper\workspace\hypernexus"
$action2 = New-ScheduledTaskAction -Execute "C:\Program Files\nodejs\node.exe" -Argument "C:\Users\hyper\workspace\hypernexus\apps\web\node_modules\.bin\next.cmd dev -p 7779" -WorkingDirectory "C:\Users\hyper\workspace\hypernexus\apps\web"
$action3 = New-ScheduledTaskAction -Execute "pythonw" -Argument "-u C:\Users\hyper\workspace\hypernexus\watchdog.py" -WorkingDirectory "C:\Users\hyper\workspace\hypernexus"

$trigger = New-ScheduledTaskTrigger -AtLogOn
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable

try {
    Register-ScheduledTask -TaskName "HyperNexus Kernel" -Action $action1 -Trigger $trigger -Settings $settings -Force
    Write-Host "✅ HyperNexus Kernel scheduled task created"
} catch {
    Write-Host "❌ Kernel: $_"
}

try {
    Register-ScheduledTask -TaskName "HyperNexus Dashboard" -Action $action2 -Trigger $trigger -Settings $settings -Force
    Write-Host "✅ HyperNexus Dashboard scheduled task created"
} catch {
    Write-Host "❌ Dashboard: $_"
}

try {
    Register-ScheduledTask -TaskName "HyperNexus Watchdog" -Action $action3 -Trigger $trigger -Settings $settings -Force
    Write-Host "✅ HyperNexus Watchdog scheduled task created"
} catch {
    Write-Host "❌ Watchdog: $_"
}

Write-Host ""
Write-Host "Done. Tasks will start on next login."
Write-Host "To run now: Start-ScheduledTask -TaskName 'HyperNexus Kernel'"
