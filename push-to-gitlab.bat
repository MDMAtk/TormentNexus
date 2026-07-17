@echo off
echo Pushing all changes to GitLab...
echo.

REM Stage all changes
git add -A

REM Commit with descriptive message
git commit -m "feat: corporate mode - HyperNexus/TormentNexus dual edition support

- Added branding configuration system (go/internal/config/branding.go)
- Updated system tray to use dynamic branding
- Created unified installer supporting both editions (installer/hypernexus.nsi)
- Added build script for both editions (build-editions.bat)
- Added corporate mode documentation (docs/corporate-mode.md)
- Added sample configuration files (config/)
- Removed GPL license, kept MIT license
- Updated all references to use GitLab repository (gitlab.com/robertpelloni/HyperNexus)
- Updated Go module path to GitLab
- Updated landing pages, scripts, and tools for GitLab"

REM Push to GitLab
git push origin main

echo.
echo Done!
pause
