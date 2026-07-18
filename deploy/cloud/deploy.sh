#!/bin/bash
set -e

echo "========================================="
echo "  HyperNexus Cloud Deployment"
echo "========================================="
echo ""

# Check if Docker is installed
if ! command -v docker &>/dev/null; then
	echo "Error: Docker is not installed"
	echo "Please install Docker first: https://docs.docker.com/get-docker/"
	exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &>/dev/null; then
	echo "Error: Docker Compose is not installed"
	echo "Please install Docker Compose first: https://docs.docker.com/compose/install/"
	exit 1
fi

# Create necessary directories
echo "Creating directories..."
mkdir -p /var/lib/hypernexus
mkdir -p /var/backups/hypernexus
mkdir -p ssl

# Check for SSL certificates
if [ ! -f ssl/cert.pem ] || [ ! -f ssl/key.pem ]; then
	echo ""
	echo "Warning: SSL certificates not found"
	echo "The server will run on HTTP only"
	echo ""
	echo "To enable HTTPS:"
	echo "1. Generate SSL certificates (e.g., using Let's Encrypt)"
	echo "2. Copy them to ssl/cert.pem and ssl/key.pem"
	echo "3. Uncomment the SSL section in nginx.conf"
	echo "4. Restart the deployment"
	echo ""
fi

# Build and start services
echo "Building and starting services..."
docker-compose build
docker-compose up -d

echo ""
echo "========================================="
echo "  Deployment Complete!"
echo "========================================="
echo ""
echo "Services:"
echo "  - Cloud API: http://localhost:7778"
echo "  - Dashboard: http://localhost:3000"
echo "  - Nginx: http://localhost:80"
echo ""
echo "Next steps:"
echo "1. Configure DNS for cloud.hypernexus.site"
echo "2. Set up SSL certificates"
echo "3. Create your first account"
echo ""
echo "Useful commands:"
echo "  - View logs: docker-compose logs -f"
echo "  - Stop services: docker-compose down"
echo "  - Restart services: docker-compose restart"
echo ""
