# TormentNexus Deployment Instructions

_This document contains the latest deployment instructions for the TormentNexus Universal AI Dashboard and Cognitive Control Plane (TormentNexus)._

## Prerequisites

1.  **Node.js**: >= 24.x (For maximum runtime compatibility)
2.  **pnpm**: Recommended package manager (`npm install -g pnpm@10.28`)
3.  **Go**: >= 1.23 (For the Go control plane sidecar)
4.  **Git**: For submodule fetching and version control.

## Initial Setup

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/NexusSoftMDMA/TormentNexus.git
    cd tormentnexus
    ```

2.  **Initialize Submodules**:
    ```bash
    git submodule update --init --recursive
    ```

3.  **Install Dependencies**:
    ```bash
    pnpm install
    ```
    *Note: If you are running Node 24, you must rebuild `better-sqlite3` bindings post-install:*
    ```bash
    pnpm rebuild better-sqlite3
    ```

4.  **Environment Variables**:
    Copy `.env.example` to `.env` and fill in the required API keys (OpenAI, Anthropic, Gemini, etc.).
    ```bash
    cp .env.example .env
    ```

## Running the Platform

TormentNexus is designed as a long-running service that manages PC memory, CPU, disk, and bandwidth usage.

### Standard Build & Run (Production/Development)

To run the TypeScript monorepo and its dashboard in one shot:
```bash
pnpm run build
pnpm run start
```
This will compile all TypeScript packages, build the web assets, and launch the primary Node.js CLI orchestrator alongside the web dashboard.

### Windows Startup Script

Use the provided startup batch file:
```powershell
.\start.bat
```
`start.bat` defaults to `pnpm run build:workspace` (skipping extension-only build stages for a much faster boot). It also triggers a native-runtime preflight to verify SQLite and Electron bindings.

**Startup Overrides**:
- **Bypass Install**: `set TORMENTNEXUS_SKIP_INSTALL=1`
- **Bypass Build**: `set TORMENTNEXUS_SKIP_BUILD=1`
- **Force Full Monorepo Build**: `set TORMENTNEXUS_FULL_BUILD=1`
- **Bypass Native Preflight**: `set TORMENTNEXUS_SKIP_NATIVE_PREFLIGHT=1`

### Linux/macOS Startup Script

```bash
./start.sh
```

### Start Electron Maestro Separately

Maestro is launched independently from the main control plane:
```bash
pnpm -C apps/maestro start
```

---

## Go Sidecar Kernel

### High-Performance Native Tools
TormentNexus includes native Go implementations for core development tools. Ensure these are in your PATH for maximum performance:
- **Ripgrep (rg)**: Required for high-speed regex search.
- **Anyquery**: SQL interface to various data sources.
- **Supervisor Monitor**: Background watchdog for session inactivity.

To build and run the Go control plane sidecar alongside the main TS engine:
```bash
cd go
go run ./cmd/tormentnexus serve
```
Alternatively, build the binary:
```bash
cd go
go build -buildvcs=false ./cmd/tormentnexus
```

---

## Extensions

To compile VS Code and Chrome/Firefox browser agents:
```bash
pnpm run build:extensions
```

## Production Docker

Build the production bundle inside a container:
```bash
docker build -f Dockerfile.prod -t tormentnexus:latest .
docker run -p 3000:3000 -p 4000:4000 -v tormentnexus-data:/root/.tormentnexus tormentnexus:latest
```

---

## Package Manager Requirement

**pnpm v10 is required.** The root `package.json` locks `packageManager: pnpm@10.28.0`. Using pnpm v9 or below will produce `ERR_PNPM_BAD_PM_VERSION` and fail the build.
```bash
npm install -g pnpm@10
```

## Release Gate Validation

Before checking in code, run the CI release verification pipeline:
```bash
pnpm run check:release-gate:ci
```
This will execute standard placeholder checks, perform type-checking across `packages/core`, and lint the entire workspace.

To perform strict visual and screenshot sync checks:
```bash
pnpm run check:release-gate:ci:strict-visuals
```

---

## Ports & Endpoints

| Service | Default Port | Override | Health Check Endpoint |
|---------|-------------|----------|-----------------------|
| Go Kernel Sidecar | 7778 | `tormentnexus serve --port <n>` | `http://localhost:7778/health` |
| Web Dashboard (Next.js) | 7779 | `PORT` | `http://localhost:7779/dashboard` |
| Socket.io Swarm Server | 3001 | — | — |
| Orchestrator tRPC | 3847 | `TORMENTNEXUS_ORCHESTRATOR_PORT` | `http://localhost:3847` |

## Health Checks
- `http://localhost:7778/health` - Go Kernel Sidecar Health
- `http://localhost:7779/dashboard` - Web Dashboard


---

## Corporate Cloud Deployment (HyperNexus)

HyperNexus supports enterprise cloud requirements including Single Sign-On (SSO), Role-Based Access Control (RBAC), Server-Sent Events (SSE) MCP network connections with authentication, and containerized data isolation.

### 1. SSO & RBAC System Integration
* **OIDC Identity Federation**: HyperNexus interfaces with industry-standard OIDC providers (Okta, Azure AD, Auth0, etc.). Toggle **SSO Active** in the settings console and configure your Client ID, Secret, and Issuer URL.
* **Granular RBAC**: Manage functional capabilities per identity role (Admin, Operator, Viewer) through the interactive role permission matrix.

### 2. Authenticated SSE MCP Connector
To connect external tools/processes to the cloud control plane:
1. Enable **SSE Client Authentication** in the Cloud Settings view.
2. Generate an SSE Client token (`CLOUDMCP_SSE_AUTH_TOKEN`).
3. Connect using standard Model Context Protocol (MCP) clients to:
   ```http
   GET http://YOUR_CLOUD_IP:4300/api/sse?token=hk_prod_YOUR_TOKEN
   ```
4. Client requests are POSTed as JSON-RPC payloads to `/api/sse/message?sessionId=client-xyz`.

### 3. Multi-Tenant Docker Process Isolation
For secure data and process separation per account:
* Spawn account services using the isolated compose blueprint:
  ```bash
  export TENANT_ID="tenant-company-a"
  export TENANT_PORT="3005"
  export CLOUDMCP_SSE_AUTH_TOKEN="hk_prod_secure_token_for_tenant_a"
  docker-compose -f docker-compose.isolated.yml up -d
  ```
* Resource boundaries, directory mounting (`/var/lib/hypernexus/tenants/tenant-company-a/`), and network interfaces are isolated automatically.

### 4. Cloud Autoscaling VM Groups
To configure automated horizontal scaling:
1. **VM Machine Image**: Create a base image (AMI/Snapshot) with Docker and Compose installed. Configure `systemd` to start the hypernexus orchestrator.
2. **Stateless Scale**: Ensure VM instances join a load balancer. Keep state stored in an external shared DB (e.g. Postgres RDS/Cloud SQL) and cache layers (Redis Enterprise) as defined in `docker-compose.yml`.
3. **Autoscaling Policies**: Setup CPU threshold triggers (e.g. scale up at >=70% CPU, scale down at <=30% CPU) using AWS Auto Scaling Groups, Azure VMSS, or Google Compute Engine Managed Instance Groups.

