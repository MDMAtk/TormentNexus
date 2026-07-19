# TormentNexus Universal Installer and Service Supervisor Setup
# =============================================================
# Run in PowerShell (Administrator recommended for Windows Service registration).

$ErrorActionPreference = "Stop"
$originalDir = Get-Location

# Ensure we're in the repository root
$scriptPath = $MyInvocation.MyCommand.Path
if ($scriptPath) {
    $repoRoot = Split-Path -Parent $scriptPath
    Set-Location $repoRoot
} else {
    $repoRoot = (Get-Item .).FullName
}

Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host "         TormentNexus System Installer & Supervisord     " -ForegroundColor Cyan
Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host "Working Directory: $repoRoot" -ForegroundColor DarkGray
Write-Host ""

# 1. Helper: Check Administrative Privileges
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if ($isAdmin) {
    Write-Host "[Info] Running with Administrator privileges. Service registration is enabled." -ForegroundColor Green
} else {
    Write-Host "[Warning] Not running as Administrator. Service registration will be skipped." -ForegroundColor Yellow
    Write-Host "          (Local user-space execution fallback will be used)." -ForegroundColor Yellow
}
Write-Host ""

# 2. Kill Stale Services to prevent lock contention
Write-Host "--- 1. Pruning Existing TormentNexus Processes ---" -ForegroundColor Cyan
$ports = @(7778, 7779, 4000)
foreach ($port in $ports) {
    Write-Host "Checking port $port..." -ForegroundColor DarkGray
    $netstat = netstat -ano | Select-String "LISTENING" | Select-String ":$port "
    if ($netstat) {
        foreach ($line in $netstat) {
            if ($line -match '\s+(\d+)$') {
                $targetPid = $matches[1]
                Write-Host "  Found process $targetPid listening on port $port. Terminating..." -ForegroundColor Yellow
                Stop-Process -Id $targetPid -Force -ErrorAction SilentlyContinue
            }
        }
    }
}

# Kill processes by name
$names = @("tormentnexus", "watchdog")
foreach ($name in $names) {
    $procs = Get-Process -Name $name -ErrorAction SilentlyContinue
    if ($procs) {
        Write-Host "Terminating active process instances for: $name" -ForegroundColor Yellow
        $procs | Stop-Process -Force -ErrorAction SilentlyContinue
    }
}
Write-Host "Clean state initialized." -ForegroundColor Green
Write-Host ""

# 3. Installing dependencies & building workspaces
Write-Host "--- 2. Building Monorepo Workspaces & Packages ---" -ForegroundColor Cyan
Write-Host "Installing NPM dependencies..." -ForegroundColor DarkGray
pnpm install --ignore-scripts

Write-Host "Compiling package workspaces..." -ForegroundColor DarkGray
node scripts/build_all.mjs --workspace-only

Write-Host "Compiling native Go sidecar daemon..." -ForegroundColor DarkGray
if (Test-Path .\build.bat) {
    cmd.exe /c .\build.bat
} else {
    Set-Location go
    go build -buildvcs=false -o ..\bin\tormentnexus.exe ./cmd/tormentnexus
    Set-Location $repoRoot
}
Write-Host "Workspace build completed successfully." -ForegroundColor Green
Write-Host ""

# 4. Building extensions & plugins
Write-Host "--- 3. Compiling System Plugins & Extensions ---" -ForegroundColor Cyan
Write-Host "Building browser extensions (Chromium & Firefox)..." -ForegroundColor DarkGray
node scripts/build_all.mjs --extensions-only

Write-Host "Packaging VS Code Extension..." -ForegroundColor DarkGray
if (Test-Path .\apps\vscode) {
    Set-Location .\apps\vscode
    pnpm run build
    pnpm run package
    Set-Location $repoRoot
}
Write-Host "Plugins and Extensions built successfully." -ForegroundColor Green
Write-Host ""

# 5. Installing MCP Client configurations
Write-Host "--- 4. Registering MCP Configurations ---" -ForegroundColor Cyan
if (Test-Path .\scripts\install-mcp-clients.py) {
    python .\scripts\install-mcp-clients.py
} else {
    Write-Host "[Warning] install-mcp-clients.py script not found. Skipping MCP setup." -ForegroundColor Yellow
}
Write-Host ""

# 6. Service Installation/Execution Configuration
Write-Host "--- 5. Launching TormentNexus Services ---" -ForegroundColor Cyan
if ($isAdmin) {
    Write-Host "Registering and starting Windows Services..." -ForegroundColor DarkGray
    if (Test-Path .\install_services.bat) {
        cmd.exe /c .\install_services.bat
    } else {
        Write-Host "install_services.bat not found. Spawning services in background..." -ForegroundColor Yellow
        $runServicesLocally = $true
    }
} else {
    $runServicesLocally = $true
}

if ($runServicesLocally) {
    Write-Host "Starting Go Sidecar (port 7778) in background..." -ForegroundColor DarkGray
    Start-Process -FilePath "bin\tormentnexus.exe" -ArgumentList "serve --port 7778" -WindowStyle Hidden

    Write-Host "Starting Next.js Dashboard (port 7779) in background..." -ForegroundColor DarkGray
    $dashPath = Join-Path $repoRoot "apps\web"
    Start-Process -FilePath "cmd.exe" -ArgumentList "/c cd /d `"$dashPath`" && set NEXT_PRIVATE_DISABLE_TURBOPACK_CACHE=1 && npx next dev --port 7779" -WindowStyle Hidden

    Write-Host "Starting Watchdog Supervisor in background..." -ForegroundColor DarkGray
    Start-Process -FilePath "python.exe" -ArgumentList "watchdog.py" -WindowStyle Hidden
}

Write-Host ""
Write-Host "=========================================================" -ForegroundColor Green
Write-Host "  Installation, Compilation & Service Start Complete! ✅  " -ForegroundColor Green
Write-Host "=========================================================" -ForegroundColor Green
Write-Host "Go sidecar:  http://127.0.0.1:7778" -ForegroundColor Green
Write-Host "Dashboard:   http://127.0.0.1:7779/dashboard" -ForegroundColor Green
Write-Host ""

Set-Location $originalDir
