# Corporate Mode (HyperNexus)

HyperNexus supports two editions:

1. **HyperNexus** - Open source, free edition
2. **HyperNexus** - Corporate, closed source edition with cloud connectivity

## Overview

The corporate mode feature allows the same codebase to be branded and configured differently based on the installation type. This enables:

- Different product names and branding
- Cloud connectivity for corporate edition
- Separate configuration directories
- Different Windows registry keys

## Building the Editions

### Using the Build Script

Run `build-editions.bat` and select the edition to build:

```batch
build-editions.bat
```

### Manual Build

#### HyperNexus (Open Source)

```batch
cd go
go build -ldflags "-X gitlab.com/robertpelloni/HyperNexus/internal/config.DefaultEdition=hypernexus" -o ..\bin\hypernexus.exe ./cmd/hypernexus

cd ..\installer
makensis /DHYPERNEXUS=1 hypernexus.nsi
```

#### HyperNexus (Corporate)

```batch
cd go
go build -ldflags "-X gitlab.com/robertpelloni/HyperNexus/internal/config.DefaultEdition=hypernexus" -o ..\bin\hypernexus.exe ./cmd/hypernexus

cd ..\installer
makensis /DCORPORATE=1 hypernexus.nsi
```

## Configuration

### Branding Configuration File

The branding configuration is stored in `branding.json` in the user's config directory:

- **HyperNexus**: `~/.hypernexus/branding.json`
- **HyperNexus**: `~/.hypernexus/branding.json`

Example configuration:

```json
{
  "edition": "hypernexus",
  "product_name": "HyperNexus",
  "company_name": "HyperNexus Corp",
  "tray_tooltip": "HyperNexus (Running)",
  "dashboard_title": "HyperNexus Dashboard",
  "config_dir": ".hypernexus",
  "registry_key": "HyperNexus",
  "cloud_endpoint": "https://api.hypernexus.io",
  "cloud_auth": ""
}
```

### Environment Variables

The following environment variables can override the configuration:

| Variable | Description |
|----------|-------------|
| `HN_EDITION` | Set to `hypernexus` or `corporate` for corporate mode |
| `HN_CLOUD_ENDPOINT` | Cloud API endpoint URL |
| `HN_CLOUD_AUTH` | Authentication token for cloud connection |

## Cloud Connectivity (Corporate Edition)

The corporate edition supports Streamable HTTP connectivity to the HyperNexus Cloud:

### Configuration

1. Set the cloud endpoint:

   ```batch
   setx HN_CLOUD_ENDPOINT "https://api.hypernexus.io"
   ```

2. Set the authentication token:

   ```batch
   setx HN_CLOUD_AUTH "your-auth-token"
   ```

3. Restart the application to apply changes.

### Using the Installer

The corporate installer includes a "Connect to Cloud" shortcut that guides users through the cloud connection setup.

## API Reference

### config.GetBranding()

Returns the current branding configuration:

```go
branding := config.GetBranding()
fmt.Println(branding.ProductName) // "HyperNexus" or "HyperNexus"
```

### config.IsCorporateMode()

Returns `true` if running in corporate mode:

```go
if config.IsCorporateMode() {
    // Corporate-specific logic
}
```

### config.IsCloudConnected()

Returns `true` if connected to a cloud endpoint:

```go
if config.IsCloudConnected() {
    // Cloud-specific logic
}
```

### config.SetBranding()

Updates the branding configuration at runtime:

```go
config.SetBranding(&config.BrandingConfig{
    Edition:     config.EditionHyperNexus,
    ProductName: "HyperNexus",
    // ... other fields
})
```

## Directory Structure

```
~/.hypernexus/          # HyperNexus config directory
├── branding.json         # Branding configuration
├── config.yaml           # Main configuration
├── memory/               # Memory storage
└── ...

~/.hypernexus/            # HyperNexus config directory
├── branding.json         # Branding configuration
├── config.yaml           # Main configuration
├── connect-cloud.bat     # Cloud connection script
├── memory/               # Memory storage
└── ...
```

## Windows Registry

The edition information is stored in the Windows registry:

```
HKCU\Software\HyperNexus     # HyperNexus
HKCU\Software\HyperNexus       # HyperNexus
```

The registry key contains:

- Installation path
- Edition type
- Version information

## Troubleshooting

### Wrong Edition Installed

If the wrong edition is installed:

1. Uninstall the current edition
2. Delete the config directory:
   - HyperNexus: `~/.hypernexus`
   - HyperNexus: `~/.hypernexus`
3. Reinstall the correct edition

### Cloud Connection Issues

If the corporate edition cannot connect to the cloud:

1. Verify environment variables are set:

   ```batch
   echo %HN_CLOUD_ENDPOINT%
   echo %HN_CLOUD_AUTH%
   ```

2. Check network connectivity to the cloud endpoint

3. Verify the authentication token is valid

### Switching Editions

To switch from HyperNexus to HyperNexus:

1. Export any important data from HyperNexus
2. Uninstall HyperNexus
3. Install HyperNexus
4. Import data if needed

## Security Considerations

### Corporate Edition

- All cloud communications use HTTPS
- Authentication tokens are stored securely
- Local data is encrypted at rest

### Open Source Edition

- All data remains local
- No cloud connectivity
- Full control over data

## Support

- **HyperNexus**: Community support via GitHub
- **HyperNexus**: Enterprise support via HyperNexus Corp
