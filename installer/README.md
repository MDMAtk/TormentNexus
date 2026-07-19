# TormentNexus Installer

## Quick Install

### Windows

1. Download `tormentnexus-setup.exe` from [Releases](https://github.com/MDMAtk/TormentNexus/releases)
2. Run the installer as Administrator
3. Follow the installation wizard
4. Launch TormentNexus from Start Menu or Desktop shortcut

### Linux/Mac

```bash
# Download and run the installer script
curl -fsSL https://raw.githubusercontent.com/MDMAtk/TormentNexus/main/scripts/install.sh | bash
```

Or manually:

```bash
# Clone the repository
git clone https://github.com/MDMAtk/TormentNexus.git
cd TormentNexus

# Run the installer
chmod +x scripts/install.sh
./scripts/install.sh
```

## Build Installer from Source

### Prerequisites

- Go 1.25+
- Node.js 24+
- pnpm
- NSIS (Windows only)

### Windows Installer

```bash
# Build the Go binary
cd go
go build -buildvcs=false -o ../bin/tormentnexus.exe ./cmd/tormentnexus

# Build the NSIS installer
cd ../installer
makensis tormentnexus.nsi
```

### Linux/Mac

```bash
# Build the Go binary
cd go
go build -buildvcs=false -o ../bin/tormentnexus ./cmd/tormentnexus

# Make the installer executable
chmod +x ../scripts/install.sh
```

## Configuration

After installation, the configuration file is located at:

- Windows: `%USERPROFILE%\.tormentnexus\config.yaml`
- Linux/Mac: `~/.tormentnexus/config.yaml`

### Default Configuration

```yaml
# TormentNexus Configuration
host: 127.0.0.1
port: 7778

# Memory Configuration
memory:
  l2_enabled: true
  l3_enabled: true
  l4_enabled: false

# Provider Configuration
providers:
  deepseek:
    enabled: true
    api_key: ""
  lmstudio:
    enabled: true
    url: http://127.0.0.1:1234
```

## Uninstalling

### Windows

1. Use "Add or Remove Programs" in Windows Settings
2. Or run `C:\Program Files\TormentNexus\uninstall.bat`

## Troubleshooting

### Port Already in Use

If port 7778 is already in use, edit the configuration file:

```yaml
port: 7779  # Use a different port
```

### Permission Denied

On Linux/Mac, ensure the binary is executable:

```bash
chmod +x ~/.local/bin/tormentnexus
```

### Firewall Issues

Ensure your firewall allows connections to the TormentNexus port (default: 7778).

## Support

For issues and questions:

- GitHub Issues: <https://github.com/MDMAtk/TormentNexus/issues>
- Documentation: <https://tormentnexus.site/docs>
