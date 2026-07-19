<# 
    TormentNexus Windows Installer
    Graphical installer with progress UI
#>

Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing
[System.Windows.Forms.Application]::EnableVisualStyles()

$ErrorActionPreference = "Stop"

# ── Configuration ──
$REPO = "MDMAtk/TormentNexus"
$VERSION = "1.0.0-b4"
$INSTALL_DIR = "$env:USERPROFILE\.tormentnexus"
$BINARY_NAME = "tormentnexus.exe"
$SERVICE_NAME = "TormentNexus"

# ── Main Form ──
$form = New-Object System.Windows.Forms.Form
$form.Text = "TormentNexus Installer"
$form.Size = New-Object System.Drawing.Size(520, 420)
$form.StartPosition = "CenterScreen"
$form.FormBorderStyle = "FixedDialog"
$form.MaximizeBox = $false
$form.BackColor = [System.Drawing.Color]::FromArgb(15, 15, 25)
$form.ForeColor = [System.Drawing.Color]::White

# ── Logo / Title ──
$lblTitle = New-Object System.Windows.Forms.Label
$lblTitle.Text = "TORMENTNEXUS"
$lblTitle.Font = New-Object System.Drawing.Font("Consolas", 24, [System.Drawing.FontStyle]::Bold)
$lblTitle.ForeColor = [System.Drawing.Color]::FromArgb(0, 255, 136)
$lblTitle.AutoSize = $true
$lblTitle.Location = New-Object System.Drawing.Point(120, 20)
$form.Controls.Add($lblTitle)

$lblSubtitle = New-Object System.Windows.Forms.Label
$lblSubtitle.Text = "AI Control Plane with Persistent Memory"
$lblSubtitle.Font = New-Object System.Drawing.Font("Segoe UI", 10)
$lblSubtitle.ForeColor = [System.Drawing.Color]::FromArgb(150, 150, 170)
$lblSubtitle.AutoSize = $true
$lblSubtitle.Location = New-Object System.Drawing.Point(140, 60)
$form.Controls.Add($lblSubtitle)

# ── Stats Panel ──
$panelStats = New-Object System.Windows.Forms.Panel
$panelStats.Size = New-Object System.Drawing.Size(460, 50)
$panelStats.Location = New-Object System.Drawing.Point(25, 95)
$panelStats.BackColor = [System.Drawing.Color]::FromArgb(25, 25, 40)
$form.Controls.Add($panelStats)

@(
    @{X=10; Text="26,180"; Label="MCP Tools"},
    @{X=160; Text="4-Tier"; Label="Memory"},
    @{X=310; Text="100%"; Label="Open Source"}
) | ForEach-Object {
    $lbl = New-Object System.Windows.Forms.Label
    $lbl.Text = $_.Text
    $lbl.Font = New-Object System.Drawing.Font("Segoe UI", 16, [System.Drawing.FontStyle]::Bold)
    $lbl.ForeColor = [System.Drawing.Color]::FromArgb(0, 255, 136)
    $lbl.AutoSize = $true
    $lbl.Location = New-Object System.Drawing.Point($_.X, 5)
    $panelStats.Controls.Add($lbl)
    
    $lbl2 = New-Object System.Windows.Forms.Label
    $lbl2.Text = $_.Label
    $lbl2.Font = New-Object System.Drawing.Font("Segoe UI", 8)
    $lbl2.ForeColor = [System.Drawing.Color]::FromArgb(100, 100, 120)
    $lbl2.AutoSize = $true
    $lbl2.Location = New-Object System.Drawing.Point($_.X, 30)
    $panelStats.Controls.Add($lbl2)
}

# ── Progress Bar ──
$progressBar = New-Object System.Windows.Forms.ProgressBar
$progressBar.Style = "Continuous"
$progressBar.Size = New-Object System.Drawing.Size(460, 25)
$progressBar.Location = New-Object System.Drawing.Point(25, 160)
$progressBar.ForeColor = [System.Drawing.Color]::FromArgb(0, 255, 136)
$progressBar.BackColor = [System.Drawing.Color]::FromArgb(30, 30, 50)
$form.Controls.Add($progressBar)

# ── Status Label ──
$lblStatus = New-Object System.Windows.Forms.Label
$lblStatus.Text = "Ready to install"
$lblStatus.Font = New-Object System.Drawing.Font("Segoe UI", 9)
$lblStatus.ForeColor = [System.Drawing.Color]::FromArgb(180, 180, 200)
$lblStatus.Size = New-Object System.Drawing.Size(460, 20)
$lblStatus.Location = New-Object System.Drawing.Point(25, 195)
$form.Controls.Add($lblStatus)

# ── Log TextBox ──
$txtLog = New-Object System.Windows.Forms.TextBox
$txtLog.Multiline = $true
$txtLog.ScrollBars = "Vertical"
$txtLog.ReadOnly = $true
$txtLog.Size = New-Object System.Drawing.Size(460, 100)
$txtLog.Location = New-Object System.Drawing.Point(25, 220)
$txtLog.BackColor = [System.Drawing.Color]::FromArgb(10, 10, 20)
$txtLog.ForeColor = [System.Drawing.Color]::FromArgb(0, 200, 100)
$txtLog.Font = New-Object System.Drawing.Font("Consolas", 9)
$form.Controls.Add($txtLog)

# ── Buttons ──
$btnInstall = New-Object System.Windows.Forms.Button
$btnInstall.Text = "Install TormentNexus"
$btnInstall.Size = New-Object System.Drawing.Size(200, 40)
$btnInstall.Location = New-Object System.Drawing.Point(25, 330)
$btnInstall.FlatStyle = "Flat"
$btnInstall.BackColor = [System.Drawing.Color]::FromArgb(0, 200, 100)
$btnInstall.ForeColor = [System.Drawing.Color]::Black
$btnInstall.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$form.Controls.Add($btnInstall)

$btnLaunch = New-Object System.Windows.Forms.Button
$btnLaunch.Text = "Launch Dashboard"
$btnLaunch.Size = New-Object System.Drawing.Size(150, 40)
$btnLaunch.Location = New-Object System.Drawing.Point(235, 330)
$btnLaunch.FlatStyle = "Flat"
$btnLaunch.BackColor = [System.Drawing.Color]::FromArgb(50, 50, 80)
$btnLaunch.ForeColor = [System.Drawing.Color]::White
$btnLaunch.Font = New-Object System.Drawing.Font("Segoe UI", 10)
$btnLaunch.Enabled = $false
$form.Controls.Add($btnLaunch)

$btnExit = New-Object System.Windows.Forms.Button
$btnExit.Text = "Exit"
$btnExit.Size = New-Object System.Drawing.Size(75, 40)
$btnExit.Location = New-Object System.Drawing.Point(400, 330)
$btnExit.FlatStyle = "Flat"
$btnExit.BackColor = [System.Drawing.Color]::FromArgb(60, 30, 30)
$btnExit.ForeColor = [System.Drawing.Color]::White
$btnExit.Font = New-Object System.Drawing.Font("Segoe UI", 10)
$form.Controls.Add($btnExit)

# ── Helper Functions ──
function Log($msg) {
    $txtLog.AppendText("  $msg`r`n")
    $txtLog.SelectionStart = $txtLog.TextLength
    $txtLog.ScrollToCaret()
    [System.Windows.Forms.Application]::DoEvents()
}

function SetStatus($msg) {
    $lblStatus.Text = $msg
    [System.Windows.Forms.Application]::DoEvents()
}

function SetProgress($value) {
    $progressBar.Value = [Math]::Min(100, $value)
    [System.Windows.Forms.Application]::DoEvents()
}

function DownloadFile($url, $dest) {
    $wc = New-Object System.Net.WebClient
    $wc.add_DownloadProgressChanged({
        param($sender, $e)
        SetProgress $e.ProgressPercentage
    })
    $wc.DownloadFile($url, $dest)
    $wc.Dispose()
}

# ── Install Function ──
$btnInstall.Add_Click({
    $btnInstall.Enabled = $false
    
    try {
        Log "=========================================="
        Log "  TormentNexus Installer v$VERSION"
        Log "=========================================="
        Log ""
        
        # Step 1: Create directories
        SetStatus "Creating installation directory..."
        Log "[1/5] Creating installation directory..."
        
        if (!(Test-Path $INSTALL_DIR)) {
            New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
        }
        if (!(Test-Path "$INSTALL_DIR\.tormentnexus")) {
            New-Item -ItemType Directory -Path "$INSTALL_DIR\.tormentnexus" -Force | Out-Null
        }
        if (!(Test-Path "$INSTALL_DIR\logs")) {
            New-Item -ItemType Directory -Path "$INSTALL_DIR\logs" -Force | Out-Null
        }
        Log "       OK - $INSTALL_DIR"
        SetProgress 20
        
        # Step 2: Download binary
        SetStatus "Downloading TormentNexus..."
        Log "[2/5] Downloading TormentNexus binary..."
        
        $downloadUrl = "https://github.com/$REPO/releases/download/$VERSION/tormentnexus-windows-amd64.zip"
        $zipPath = "$INSTALL_DIR\tormentnexus.zip"
        
        Log "       URL: $downloadUrl"
        DownloadFile $downloadUrl $zipPath
        Log "       OK - Downloaded"
        SetProgress 50
        
        # Step 3: Extract
        SetStatus "Extracting files..."
        Log "[3/5] Extracting files..."
        
        Expand-Archive -Path $zipPath -DestinationPath $INSTALL_DIR -Force
        Remove-Item $zipPath -Force -ErrorAction SilentlyContinue
        Log "       OK - Extracted"
        SetProgress 70
        
        # Step 4: Create config
        SetStatus "Creating configuration..."
        Log "[4/5] Creating configuration..."
        
        $config = @{
            version = "1.0.0"
            server = @{
                host = "127.0.0.1"
                port = 7778
            }
            memory = @{
                enabled = $true
                tiers = @("L1", "L2", "L3", "L4")
            }
            mcp = @{
                catalog = $true
                autoInstall = $false
            }
        } | ConvertTo-Json -Depth 5
        
        $config | Out-File -FilePath "$INSTALL_DIR\.tormentnexus\config.json" -Encoding UTF8
        Log "       OK - config.json created"
        SetProgress 80
        
        # Step 5: Add to PATH
        SetStatus "Adding to PATH..."
        Log "[5/5] Adding to system PATH..."
        
        $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($currentPath -notlike "*$INSTALL_DIR*") {
            [Environment]::SetEnvironmentVariable("Path", "$currentPath;$INSTALL_DIR", "User")
            $env:Path = "$env:Path;$INSTALL_DIR"
            Log "       OK - Added to PATH"
        } else {
            Log "       OK - Already in PATH"
        }
        SetProgress 90
        
        # Create start menu shortcuts
        $startMenu = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\TormentNexus"
        if (!(Test-Path $startMenu)) {
            New-Item -ItemType Directory -Path $startMenu -Force | Out-Null
        }
        
        $shell = New-Object -ComObject WScript.Shell
        
        # Start shortcut
        $shortcut = $shell.CreateShortcut("$startMenu\Start TormentNexus.lnk")
        $shortcut.TargetPath = "$INSTALL_DIR\$BINARY_NAME"
        $shortcut.Arguments = "serve"
        $shortcut.WorkingDirectory = $INSTALL_DIR
        $shortcut.Description = "Start TormentNexus Server"
        $shortcut.Save()
        
        # Dashboard shortcut
        $shortcut = $shell.CreateShortcut("$startMenu\Dashboard.lnk")
        $shortcut.TargetPath = "http://localhost:7778"
        $shortcut.Description = "Open TormentNexus Dashboard"
        $shortcut.Save()
        
        # Desktop shortcut
        $shortcut = $shell.CreateShortcut("$env:USERPROFILE\Desktop\TormentNexus.lnk")
        $shortcut.TargetPath = "$INSTALL_DIR\$BINARY_NAME"
        $shortcut.Arguments = "serve"
        $shortcut.WorkingDirectory = $INSTALL_DIR
        $shortcut.Save()
        
        Log "       OK - Shortcuts created"
        SetProgress 100
        
        Log ""
        Log "=========================================="
        Log "  INSTALLATION COMPLETE!"
        Log "=========================================="
        Log ""
        Log "  Location: $INSTALL_DIR"
        Log "  Binary:   $INSTALL_DIR\$BINARY_NAME"
        Log ""
        Log "  Run: tormentnexus serve"
        Log "  Dashboard: http://localhost:7778"
        Log ""
        
        SetStatus "Installation complete!"
        $btnLaunch.Enabled = $true
        $btnInstall.Text = "Reinstall"
        $btnInstall.Enabled = $true
        
    } catch {
        Log ""
        Log "ERROR: $($_.Exception.Message)"
        Log ""
        Log "Please download manually from:"
        Log "https://github.com/$REPO/releases"
        SetStatus "Installation failed!"
        $btnInstall.Enabled = $true
    }
})

# ── Launch Function ──
$btnLaunch.Add_Click({
    Log "Starting TormentNexus..."
    SetStatus "Starting server..."
    
    $process = Start-Process -FilePath "$INSTALL_DIR\$BINARY_NAME" -ArgumentList "serve" -WorkingDirectory $INSTALL_DIR -PassThru -WindowStyle Hidden
    
    Start-Sleep -Seconds 3
    
    if (!$process.HasExited) {
        Log "Server started (PID: $($process.Id))"
        Log "Opening dashboard..."
        Start-Process "http://localhost:7778"
        SetStatus "Server running - Dashboard opened!"
    } else {
        Log "Server failed to start. Check logs at:"
        Log "$INSTALL_DIR\logs\"
        SetStatus "Server failed to start"
    }
})

# ── Exit Function ──
$btnExit.Add_Click({
    $form.Close()
})

# ── Show Form ──
[System.Windows.Forms.Application]::Run($form)
