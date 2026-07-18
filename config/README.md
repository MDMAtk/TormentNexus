# Configuration Files

This directory contains sample configuration files for both editions of the product.

## Editions

### HyperNexus (Open Source)

- **File**: `hypernexus-branding.json`
- **Config Directory**: `~/.hypernexus`
- **Registry Key**: `HyperNexus`
- **Cloud Connectivity**: None (local only)

### HyperNexus (Corporate)

- **File**: `hypernexus-branding.json`
- **Config Directory**: `~/.hypernexus`
- **Registry Key**: `HyperNexus`
- **Cloud Connectivity**: Streamable HTTP to HyperNexus Cloud

## Usage

### Installing with Branding

1. Copy the appropriate branding file to your config directory:

   ```batch
   copy config\hypernexus-branding.json %USERPROFILE%\.hypernexus\branding.json
   ```

2. Edit the branding file to customize settings

3. Run the application - it will automatically detect the branding configuration

### Environment Variables

Override branding settings with environment variables:

```batch
set HN_EDITION=hypernexus
set HN_CLOUD_ENDPOINT=https://api.hypernexus.io
set HN_CLOUD_AUTH=your-auth-token
```

## Branding Configuration Fields

| Field | Type | Description |
|-------|------|-------------|
| `edition` | string | `hypernexus` or `hypernexus` |
| `product_name` | string | Display name for the product |
| `company_name` | string | Company/organization name |
| `tray_tooltip` | string | System tray tooltip text |
| `dashboard_title` | string | Dashboard page title |
| `config_dir` | string | Configuration directory name |
| `registry_key` | string | Windows registry key name |
| `cloud_endpoint` | string | Cloud API endpoint (corporate only) |
| `cloud_auth` | string | Cloud authentication token (corporate only) |

## Building Custom Editions

To create a custom edition:

1. Create a new branding JSON file with your custom settings
2. Modify the installer script (`installer/hypernexus.nsi`) to include your edition
3. Build the Go binary with your custom edition name
4. Package the installer with your branding

Example:

```batch
cd go
go build -ldflags "-X gitlab.com/robertpelloni/HyperNexus/internal/config.DefaultEdition=myedition" -o ..\bin\myapp.exe ./cmd/hypernexus
```
