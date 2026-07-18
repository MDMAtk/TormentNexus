#!/bin/bash
set -e

echo "========================================="
echo "  SSL Certificate Setup for HyperNexus"
echo "========================================="
echo ""

DOMAIN="cloud.hypernexus.site"
EMAIL="admin@hypernexus.site"

# Check if certbot is installed
if ! command -v certbot &> /dev/null; then
    echo "Certbot not found. Installing..."
    
    # Detect OS and install certbot
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get &> /dev/null; then
            sudo apt-get update
            sudo apt-get install -y certbot python3-certbot-nginx
        elif command -v yum &> /dev/null; then
            sudo yum install -y certbot python3-certbot-nginx
        else
            echo "Error: Unsupported Linux distribution"
            echo "Please install certbot manually: https://certbot.eff.org/"
            exit 1
        fi
    else
        echo "Error: Unsupported operating system"
        echo "Please install certbot manually: https://certbot.eff.org/"
        exit 1
    fi
fi

echo "Stopping nginx temporarily..."
docker-compose stop nginx 2>/dev/null || true

echo ""
echo "Generating SSL certificate for $DOMAIN..."
echo ""

# Generate certificate using standalone mode
sudo certbot certonly \
    --standalone \
    --preferred-challenges http \
    -d "$DOMAIN" \
    --email "$EMAIL" \
    --agree-tos \
    --non-interactive

echo ""
echo "Copying certificates..."

# Copy certificates to ssl directory
sudo cp "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ssl/cert.pem
sudo cp "/etc/letsencrypt/live/$DOMAIN/privkey.pem" ssl/key.pem
sudo chown $(whoami):$(whoami) ssl/*.pem

echo ""
echo "Setting up auto-renewal..."

# Create renewal hook
sudo tee /etc/letsencrypt/renewal-hooks/deploy/hypernexus.sh > /dev/null << 'EOF'
#!/bin/bash
# Copy renewed certificates
cp /etc/letsencrypt/live/cloud.hypernexus.site/fullchain.pem /path/to/HyperNexus/deploy/cloud/ssl/cert.pem
cp /etc/letsencrypt/live/cloud.hypernexus.site/privkey.pem /path/to/HyperNexus/deploy/cloud/ssl/key.pem

# Restart nginx
cd /path/to/HyperNexus/deploy/cloud
docker-compose restart nginx
EOF

sudo chmod +x /etc/letsencrypt/renewal-hooks/deploy/hypernexus.sh

echo ""
echo "Enabling SSL in nginx configuration..."

# Update nginx.conf to enable SSL
sed -i 's/# server {/server {/g' nginx.conf
sed -i 's/#     listen 443/    listen 443/g' nginx.conf
sed -i 's/#     server_name/    server_name/g' nginx.conf
sed -i 's/#     ssl_certificate/    ssl_certificate/g' nginx.conf
sed -i 's/#     ssl_certificate_key/    ssl_certificate_key/g' nginx.conf
sed -i 's/#     location/    location/g' nginx.conf
sed -i 's/#         proxy_pass/        proxy_pass/g' nginx.conf
sed -i 's/#         proxy_set_header/        proxy_set_header/g' nginx.conf
sed -i 's/#     }/    }/g' nginx.conf
sed -i 's/# }/}/g' nginx.conf

echo ""
echo "Restarting nginx..."
docker-compose start nginx

echo ""
echo "========================================="
echo "  SSL Setup Complete!"
echo "========================================="
echo ""
echo "Your site is now available at:"
echo "  https://$DOMAIN"
echo ""
echo "Certificate will auto-renew via certbot."
echo ""
echo "To test SSL configuration:"
echo "  curl -I https://$DOMAIN"
echo ""
echo "To check certificate expiration:"
echo "  sudo certbot certificates"
echo ""
