#!/bin/bash
# =============================================================================
# TormentNexus Stripe Webhook → Auto-Provision
# Called by the Go sidecar when a new subscription is created.
# Arguments passed as env vars from the billing webhook handler:
#   TENANT_ID, TENANT_EMAIL, TENANT_ORG, PLAN_TYPE
# =============================================================================
set -euo pipefail

TENANT_ID="${1:?Usage: $0 <tenant-id> <email> <org>}"
TENANT_EMAIL="${2:-unknown@unknown.com}"
TENANT_ORG="${3:-$TENANT_ID}"
PLAN_TYPE="${4:-basic}"

DATA_ROOT="/var/lib/hypernexus/tenants"
PORT_START=3001

echo "=== Auto-provisioning from Stripe webhook ==="
echo "   Tenant:  $TENANT_ID"
echo "   Email:   $TENANT_EMAIL"
echo "   Org:     $TENANT_ORG"
echo "   Plan:    $PLAN_TYPE"

# Find next available port
LAST_PORT=$(ls -1 "${DATA_ROOT}/" 2>/dev/null | grep -v '_deprovisioned' | while read dir; do
	grep -l "${dir}" /etc/nginx/sites-enabled/*.conf 2>/dev/null || true
done | grep -oP '(?<=proxy_pass http://127.0.0.1:)[0-9]+' | sort -n | tail -1)
NEXT_PORT=${LAST_PORT:-$PORT_START}
if [ "$NEXT_PORT" -eq "$LAST_PORT" ] 2>/dev/null; then
	NEXT_PORT=$((LAST_PORT + 1))
fi

# Allocate resource limits based on plan
case "$PLAN_TYPE" in
enterprise)
	CORE_CPU="2.0"
	CORE_MEM="4096M"
	WEB_CPU="1.5"
	WEB_MEM="2048M"
	;;
pro)
	CORE_CPU="1.5"
	CORE_MEM="2048M"
	WEB_CPU="1.0"
	WEB_MEM="1024M"
	;;
basic | *)
	CORE_CPU="0.75"
	CORE_MEM="1024M"
	WEB_CPU="0.5"
	WEB_MEM="512M"
	;;
esac

# Run provisioning
/opt/tormentnexus/deploy/tenant-provision.sh "$TENANT_ID" "$NEXT_PORT"

# Log to Stripe tracking
echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) | $TENANT_ID | $TENANT_EMAIL | $TENANT_ORG | $PLAN_TYPE | $NEXT_PORT" >>"${DATA_ROOT}/provisioning.log"

echo "=== Auto-provisioning complete ==="
echo "   Dashboard: https://${TENANT_ID}.hypernexus.site"
echo "   Port:      $NEXT_PORT"
echo "   Plan:      $PLAN_TYPE"
