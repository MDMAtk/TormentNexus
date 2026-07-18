@echo off
set PORT=7779
cd /d "C:\Users\hyper\workspace\hypernexus\apps\web"
start /MIN /B cmd /c "node ".next-build\standalone\apps\web\server.js" >nul 2>&1"
