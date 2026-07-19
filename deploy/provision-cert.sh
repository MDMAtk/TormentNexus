#!/bin/bash
# =============================================================================
# Auto-provision SSL cert for a tenant subdomain
# Usage: ./provision-cert.sh <tenant-id>
#   Issues a Let's Encrypt cert for <tenant-id>.hypernexus.site via HTTP challenge
# =============================================================================
set -euo pipefail

TENANT_ID="${1:?Usage: $0 <tenant-id>}"
DOMAIN="${TENANT_ID}.hypernexus.site"
WEBROOT="/srv/www/${TENANT_ID}.hypernexus.site"

echo "=== Provisioning SSL cert for ${DOMAIN} ==="

# Create webroot for HTTP challenge
mkdir -p "${WEBROOT}/.well-known/acme-challenge"

# Issue cert via HTTP-01 challenge
certbot certonly --webroot -w "${WEBROOT}" \
	-d "${DOMAIN}" \
	--agree-tos --email "admin@hypernexus.site" \
	--no-eff-email --non-interactive \
	--expand 2>&1

if [ $? -eq 0 ]; then
	echo "✅ Cert issued for ${DOMAIN}"

	# Create Nginx snippet for this tenant
	cat >"/etc/nginx/snippets/hypernexus-${TENANT_ID}.conf" <<SNIPPET_EOF
# Auto-generated for ${TENANT_ID}
# Issued: $(date -u +%Y-%m-%dT%H:%M:%SZ)
ssl_certificate /etc/letsencrypt/live/${DOMAIN}/fullchain.pem;
ssl_certificate_key /etc/letsencrypt/live/${DOMAIN}/privkey.pem;
SNIPPET_EOF

else
	echo "⚠️  Cert issuance failed, check logs"
fi
