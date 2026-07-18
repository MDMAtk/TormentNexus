# HyperNexus Installer

## Quick Install

### Windows

1. Download `hypernexus-setup.exe` from [Releases](https://github.com/MDMAtk/HyperNexus/releases)
2. Run the installer as Administrator
3. Follow the installation wizard
4. Launch HyperNexus from Start Menu or Desktop shortcut

### Linux/Mac

```bash
# Download and run the installer script
curl -fsSL https://raw.githubusercontent.com/MDMAtk/HyperNexus/main/scripts/install.sh | bash
```

Or manually:

```bash
# Clone the repository
git clone https://github.com/MDMAtk/HyperNexus.git
cd HyperNexus

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
go build -buildvcs=false -o ../bin/hypernexus.exe ./cmd/hypernexus

# Build the NSIS installer
cd ../installer
makensis hypernexus.nsi
```

### Linux/Mac

```bash
# Build the Go binary
cd go
go build -buildvcs=false -o ../bin/hypernexus ./cmd/hypernexus

# Make the installer executable
chmod +x ../scripts/install.sh
```

## Configuration

After installation, the configuration file is located at:

- Windows: `%USERPROFILE%\.hypernexus\config.yaml`
- Linux/Mac: `~/.hypernexus/config.yaml`

### Default Configuration

```yaml
# HyperNexus Configuration
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
2. Or run `C:\Program Files\HyperNexus\uninstall.bat`

## Troubleshooting

### Port Already in Use

If port 7778 is already in use, edit the configuration file:

```yaml
port: 7779  # Use a different port
```

### Permission Denied

On Linux/Mac, ensure the binary is executable:

```bash
chmod +x ~/.local/bin/hypernexus
```

### Firewall Issues

Ensure your firewall allows connections to the HyperNexus port (default: 7778).

## Support

For issues and questions:

- GitHub Issues: <https://github.com/MDMAtk/HyperNexus/issues>
- Documentation: <https://hypernexus.site/docs>
