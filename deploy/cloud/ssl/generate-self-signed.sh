#!/bin/bash
set -e

echo "========================================="
echo "  Self-Signed SSL Certificate Generator"
echo "========================================="
echo ""
echo "NOTE: Self-signed certificates are for TESTING only!"
echo "      Browsers will show security warnings."
echo "      Use certbot for production certificates."
echo ""

DOMAIN="cloud.hypernexus.site"
DAYS=365

# Create ssl directory if it doesn't exist
mkdir -p ssl

echo "Generating self-signed certificate for $DOMAIN..."
echo ""

# Generate private key and certificate
openssl req -x509 \
	-nodes \
	-days $DAYS \
	-newkey rsa:2048 \
	-keyout ssl/key.pem \
	-out ssl/cert.pem \
	-subj "/C=US/ST=State/L=City/O=HyperNexus/CN=$DOMAIN" \
	-addext "subjectAltName=DNS:$DOMAIN,DNS:www.$DOMAIN"

echo ""
echo "========================================="
echo "  Certificate Generated!"
echo "========================================="
echo ""
echo "Files created:"
echo "  - ssl/cert.pem (certificate)"
echo "  - ssl/key.pem (private key)"
echo ""
echo "Validity: $DAYS days"
echo ""
echo "To use with docker-compose:"
echo "  cd deploy/cloud"
echo "  docker-compose up -d"
echo ""
echo "To test:"
echo "  curl -k https://$DOMAIN"
echo ""
