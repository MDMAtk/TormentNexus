# HyperNexus Cloud Deployment Guide

## Overview

HyperNexus Cloud provides a hosted MCP server with Streamable HTTP transport, user account management, Docker container provisioning, and automatic backups.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     cloud.hypernexus.site                    │
├─────────────────────────────────────────────────────────────┤
│                         Nginx                               │
│                    (SSL + Rate Limiting)                    │
├─────────────────────────────────────────────────────────────┤
│                    HyperNexus Cloud API                     │
│              (Account + Container Management)               │
├─────────────────────────────────────────────────────────────┤
│               MCP Streamable HTTP Transport                 │
│                 (JSON-RPC over HTTP/SSE)                    │
├─────────────────────────────────────────────────────────────┤
│                   Docker Provisioner                        │
│              (User Container Management)                    │
├─────────────────────────────────────────────────────────────┤
│                     Backup System                           │
│               (Automatic + Manual Backups)                  │
└─────────────────────────────────────────────────────────────┘
```

## Features

### 1. Streamable HTTP Transport

The MCP server uses the Streamable HTTP transport protocol:

- **Endpoint**: `https://cloud.hypernexus.site/mcp/v1`
- **SSE Endpoint**: `https://cloud.hypernexus.site/mcp/v1/sse`
- **Authentication**: API key via `X-API-Key` header or `Bearer` token

### 2. Account Management

- User registration and login
- API key management
- Plan-based resource allocation
- Account status management

### 3. Docker Container Provisioning

Each user gets an isolated Docker container:

| Plan       | Memory | CPU  | Storage |
|------------|--------|------|---------|
| Free       | 512MB  | 0.5  | 1GB     |
| Starter    | 1GB    | 1    | 5GB     |
| Pro        | 2GB    | 2    | 10GB    |
| Enterprise | 4GB    | 4    | 50GB    |

### 4. Backup System

- Automatic daily backups
- Manual backup creation
- Backup restoration
- Backup listing and deletion

## Deployment

### Prerequisites

- Docker and Docker Compose
- Domain name (cloud.hypernexus.site)
- SSL certificate (Let's Encrypt recommended)

### Quick Start

1. Clone the repository:

   ```bash
   git clone https://gitlab.com/robertpelloni/HyperNexus.git
   cd HyperNexus
   ```

2. Configure environment:

   ```bash
   cp deploy/cloud/.env.example deploy/cloud/.env
   # Edit .env with your configuration
   ```

3. Start the cloud services:

   ```bash
   cd deploy/cloud
   docker-compose up -d
   ```

4. Access the dashboard:
   - Open `https://cloud.hypernexus.site`
   - Create an account
   - Start using the MCP server

### SSL Configuration

1. Generate SSL certificate:

   ```bash
   # Using Let's Encrypt
   certbot certonly --standalone -d cloud.hypernexus.site
   ```

2. Copy certificates:

   ```bash
   cp /etc/letsencrypt/live/cloud.hypernexus.site/fullchain.pem deploy/cloud/ssl/cert.pem
   cp /etc/letsencrypt/live/cloud.hypernexus.site/privkey.pem deploy/cloud/ssl/key.pem
   ```

3. Uncomment SSL section in `nginx.conf`

## API Reference

### Authentication

All API requests require authentication via API key:

```bash
# Using X-API-Key header
curl -H "X-API-Key: hn_your_api_key" https://cloud.hypernexus.site/api/cloud/auth/me

# Using Bearer token
curl -H "Authorization: Bearer hn_your_api_key" https://cloud.hypernexus.site/api/cloud/auth/me
```

### Endpoints

#### Account Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/api/cloud/auth/register` | Register new account |
| POST   | `/api/cloud/auth/login` | Login |
| GET    | `/api/cloud/auth/me` | Get current user |
| POST   | `/api/cloud/auth/api-key` | Regenerate API key |
| PUT    | `/api/cloud/account` | Update account |
| POST   | `/api/cloud/account/plan` | Update plan |

#### Container Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/api/cloud/container` | Get container info |
| POST   | `/api/cloud/container/start` | Start container |
| POST   | `/api/cloud/container/stop` | Stop container |
| POST   | `/api/cloud/container/destroy` | Destroy container |

#### Backup Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/api/cloud/backup` | Create backup |
| POST   | `/api/cloud/backup/restore` | Restore backup |
| GET    | `/api/cloud/backup/list` | List backups |

#### MCP Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/mcp/v1` | MCP JSON-RPC endpoint |
| GET    | `/mcp/v1/sse` | MCP SSE endpoint |
| GET    | `/mcp/v1/health` | Health check |

### MCP Protocol

#### Initialize

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "MyClient",
      "version": "1.0.0"
    }
  }
}
```

#### List Tools

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

#### Call Tool

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "memory_store",
    "arguments": {
      "content": "Remember this important fact",
      "tags": ["important", "fact"]
    }
  }
}
```

## Local Installation (When Using Cloud)

When using HyperNexus Cloud, you only need a lightweight local client:

### Option 1: CLI Client

```bash
npm install -g @hypernexus/cli
hypernexus connect --api-key hn_your_api_key
```

### Option 2: Browser Extension

Install the HyperNexus browser extension from the Chrome Web Store or Firefox Add-ons.

### Option 3: IDE Extension

Install the HyperNexus extension for:

- VS Code
- JetBrains IDEs
- Cursor

### Local Client Configuration

Create `~/.hypernexus/config.yaml`:

```yaml
cloud:
  endpoint: https://cloud.hypernexus.site
  api_key: hn_your_api_key
  
local:
  memory_cache: true
  offline_mode: false
```

## Maintenance

### Container Management

```bash
# List all containers
docker ps --filter "name=hypernexus-"

# View container logs
docker logs hypernexus-account_id

# Restart container
docker restart hypernexus-account_id

# Access container shell
docker exec -it hypernexus-account_id sh
```

### Database Backup

```bash
# Backup database
docker exec hypernexus-cloud sqlite3 /data/hypernexus.db ".backup /data/backup.db"

# Copy backup out of container
docker cp hypernexus-cloud:/data/backup.db ./backup.db
```

### Monitoring

```bash
# Check health
curl https://cloud.hypernexus.site/api/cloud/health

# View metrics
curl https://cloud.hypernexus.site/api/cloud/metrics
```

## Troubleshooting

### Container Won't Start

```bash
# Check container status
docker inspect hypernexus-account_id

# View container logs
docker logs hypernexus-account_id

# Check Docker daemon
systemctl status docker
```

### Connection Issues

```bash
# Test API connectivity
curl -v https://cloud.hypernexus.site/api/cloud/health

# Check DNS
nslookup cloud.hypernexus.site

# Test MCP endpoint
curl -X POST https://cloud.hypernexus.site/mcp/v1 \
  -H "X-API-Key: hn_your_api_key" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize"}'
```

### Backup Restoration

```bash
# List available backups
curl -H "X-API-Key: hn_your_api_key" https://cloud.hypernexus.site/api/cloud/backup/list

# Restore backup
curl -X POST https://cloud.hypernexus.site/api/cloud/backup/restore \
  -H "X-API-Key: hn_your_api_key" \
  -H "Content-Type: application/json" \
  -d '{"backup_file":"/var/backups/hypernexus/account_id_20260101_120000.tar.gz"}'
```

## Security Considerations

1. **API Keys**: Store securely, rotate regularly
2. **Network**: Use HTTPS only, configure firewall
3. **Containers**: Isolated per user, resource limits enforced
4. **Backups**: Encrypted at rest, regular rotation
5. **Monitoring**: Log all API access, alert on anomalies

## Support

- **Documentation**: <https://hypernexus.site/docs>
- **API Reference**: <https://cloud.hypernexus.site/api/docs>
- **Status Page**: <https://status.hypernexus.site>
- **Support Email**: <support@hypernexus.site>
