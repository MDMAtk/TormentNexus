@echo off
setlocal

echo Building HyperNexus...

REM Build Go sidecar
where go >nul 2>nul
if errorlevel 1 (
    echo [SKIP] Go not found
    goto :done
)

echo Building Go sidecar...
for /f "tokens=*" %%v in ('type VERSION') do set VER=%%v
cd go
go build -ldflags "-s -w -X gitlab.com/robertpelloni/HyperNexus/internal/buildinfo.Version=%VER%" -buildvcs=false -o ..\bin\hypernexus.exe ./cmd/hypernexus
if errorlevel 1 (
    echo [FAIL] Go build failed
    cd ..
    exit /b 1
)
cd ..
echo OK - bin\hypernexus.exe built

:done
echo.
echo Build complete.
