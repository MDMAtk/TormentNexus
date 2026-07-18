#!/bin/bash
set -e

echo "========================================="
echo "  HyperNexus SSL Complete Setup"
echo "========================================="
echo ""

DOMAIN="cloud.hypernexus.site"
HYPERNEXUS_DIR="/opt/HyperNexus"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root (sudo ./setup-ssl-complete.sh)"
    exit 1
fi

# Step 1: Install certbot if not present
echo "Step 1: Checking certbot installation..."
if ! command -v certbot &> /dev/null; then
    echo "Installing certbot..."
    apt-get update
    apt-get install -y certbot
fi
echo "✓ Certbot installed"
echo ""

# Step 2: Stop services on port 80
echo "Step 2: Stopping services on port 80..."
systemctl stop nginx 2>/dev/null || true
docker stop hypernexus-nginx 2>/dev/null || true
echo "✓ Port 80 freed"
echo ""

# Step 3: Generate SSL certificate
echo "Step 3: Generating SSL certificate..."
echo "Domain: $DOMAIN"
echo ""

certbot certonly \
    --standalone \
    --preferred-challenges http \
    -d "$DOMAIN" \
    --email "admin@hypernexus.site" \
    --agree-tos \
    --non-interactive \
    --force-renewal

echo "✓ SSL certificate generated"
echo ""

# Step 4: Copy certificates to HyperNexus directory
echo "Step 4: Copying certificates..."
mkdir -p "$HYPERNEXUS_DIR/deploy/cloud/ssl"
cp "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" "$HYPERNEXUS_DIR/deploy/cloud/ssl/cert.pem"
cp "/etc/letsencrypt/live/$DOMAIN/privkey.pem" "$HYPERNEXUS_DIR/deploy/cloud/ssl/key.pem"
chmod 644 "$HYPERNEXUS_DIR/deploy/cloud/ssl/cert.pem"
chmod 600 "$HYPERNEXUS_DIR/deploy/cloud/ssl/key.pem"
echo "✓ Certificates copied"
echo ""

# Step 5: Update nginx configuration to use SSL
echo "Step 5: Updating nginx configuration..."
cd "$HYPERNEXUS_DIR/deploy/cloud"

# Backup original nginx.conf
cp nginx.conf nginx.conf.backup

# Use SSL configuration
cp nginx-ssl.conf nginx.conf
echo "✓ Nginx configured for SSL"
echo ""

# Step 6: Set up auto-renewal
echo "Step 6: Setting up auto-renewal..."
cat > /etc/cron.d/hypernexus-ssl << EOF
# Auto-renew SSL certificate for HyperNexus
0 0 1 * * root certbot renew --quiet --deploy-hook "cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem $HYPERNEXUS_DIR/deploy/cloud/ssl/cert.pem && cp /etc/letsencrypt/live/$DOMAIN/privkey.pem $HYPERNEXUS_DIR/deploy/cloud/ssl/key.pem && cd $HYPERNEXUS_DIR/deploy/cloud && docker-compose restart nginx"
EOF
chmod 644 /etc/cron.d/hypernexus-ssl
echo "✓ Auto-renewal configured"
echo ""

# Step 7: Start services
echo "Step 7: Starting services..."
cd "$HYPERNEXUS_DIR/deploy/cloud"
docker-compose up -d
echo "✓ Services started"
echo ""

# Step 8: Test SSL
echo "Step 8: Testing SSL configuration..."
sleep 5
if curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN/api/cloud/health" | grep -q "200"; then
    echo "✓ SSL is working!"
else
    echo "⚠ SSL test failed - checking configuration..."
    docker-compose logs nginx | tail -20
fi
echo ""

echo "========================================="
echo "  SSL Setup Complete!"
echo "========================================="
echo ""
echo "Your site is now available at:"
echo "  https://$DOMAIN"
echo ""
echo "Certificate details:"
openssl x509 -in "$HYPERNEXUS_DIR/deploy/cloud/ssl/cert.pem" -noout -subject -dates
echo ""
echo "Auto-renewal is configured to run monthly."
echo ""
echo "To manually renew:"
echo "  sudo certbot renew"
echo "  sudo docker-compose restart nginx"
echo ""
echo "To check certificate status:"
echo "  sudo certbot certificates"
echo ""
