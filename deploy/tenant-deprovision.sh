#!/bin/bash
# =============================================================================
# TormentNexus Multi-Tenant Deprovisioning Script
# Usage: ./tenant-deprovision.sh <tenant-id>
# =============================================================================
set -euo pipefail

TENANT_ID="${1:?Usage: $0 <tenant-id>}"
DATA_ROOT="/var/lib/hypernexus/tenants"

echo "=== Deprovisioning tenant: $TENANT_ID ==="

# 1. Stop and remove containers
echo "Stopping containers..."
docker compose -f "/tmp/docker-compose.$TENANT_ID.yml" down -v 2>/dev/null || true
docker rm -f tn-core-${TENANT_ID} tn-web-${TENANT_ID} tn-sidecar-${TENANT_ID} 2>/dev/null || true

# 2. Remove Nginx config
rm -f "/etc/nginx/sites-enabled/hypernexus.${TENANT_ID}.conf"
nginx -t && systemctl reload nginx 2>/dev/null || true

# 3. Remove docker-compose override
rm -f "/tmp/docker-compose.$TENANT_ID.yml"

# 4. Archive data (keep for 30 days)
BACKUP_DIR="${DATA_ROOT}/_deprovisioned"
mkdir -p "$BACKUP_DIR"
if [ -d "${DATA_ROOT}/${TENANT_ID}" ]; then
	mv "${DATA_ROOT}/${TENANT_ID}" "${BACKUP_DIR}/${TENANT_ID}-$(date +%Y%m%d)"
	echo "Data archived to ${BACKUP_DIR}/${TENANT_ID}-$(date +%Y%m%d)"
fi

echo "=== Deprovisioning complete for $TENANT_ID ==="
