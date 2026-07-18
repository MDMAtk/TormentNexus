# SSL Certificate Setup for HyperNexus Cloud

This directory contains SSL certificates and setup scripts for `cloud.hypernexus.site`.

## Quick Start

### Option 1: Self-Signed Certificate (Testing)

For local development or testing:

**Linux/Mac:**

```bash
cd deploy/cloud/ssl
chmod +x generate-self-signed.sh
./generate-self-signed.sh
```

**Windows:**

```batch
cd deploy\cloud\ssl
generate-windows.bat
```

### Option 2: Let's Encrypt Certificate (Production)

For production with a real domain:

**Linux:**

```bash
cd deploy/cloud/ssl
chmod +x setup-ssl.sh
sudo ./setup-ssl.sh
```

**Requirements:**

- Domain `cloud.hypernexus.site` must point to your server
- Port 80 must be accessible from the internet
- Root/sudo access required

## Manual Setup

### Using Let's Encrypt (Certbot)

1. Install certbot:

   ```bash
   # Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install certbot

   # CentOS/RHEL
   sudo yum install certbot
   ```

2. Stop any services on port 80:

   ```bash
   sudo systemctl stop nginx
   # or
   docker-compose stop nginx
   ```

3. Generate certificate:

   ```bash
   sudo certbot certonly --standalone -d cloud.hypernexus.site
   ```

4. Copy certificates:

   ```bash
   sudo cp /etc/letsencrypt/live/cloud.hypernexus.site/fullchain.pem ssl/cert.pem
   sudo cp /etc/letsencrypt/live/cloud.hypernexus.site/privkey.pem ssl/key.pem
   ```

5. Set up auto-renewal:

   ```bash
   sudo crontab -e
   # Add this line:
   0 0 1 * * certbot renew --quiet
   ```

### Using OpenSSL (Self-Signed)

1. Generate private key:

   ```bash
   openssl genrsa -out ssl/key.pem 2048
   ```

2. Generate certificate:

   ```bash
   openssl req -x509 -new -nodes -key ssl/key.pem -sha256 -days 365 \
     -out ssl/cert.pem \
     -subj "/C=US/ST=State/L=City/O=HyperNexus/CN=cloud.hypernexus.site"
   ```

### Using mkcert (Local Development)

For local development with trusted certificates:

1. Install mkcert:

   ```bash
   # Mac
   brew install mkcert

   # Windows
   choco install mkcert

   # Linux
   sudo apt install mkcert
   ```

2. Install the local CA:

   ```bash
   mkcert -install
   ```

3. Generate certificates:

   ```bash
   cd deploy/cloud/ssl
   mkcert -cert-file cert.pem -key-file key.pem cloud.hypernexus.site localhost 127.0.0.1
   ```

## Enabling SSL in Nginx

After generating certificates, update `nginx.conf`:

1. Uncomment the SSL server block:

   ```nginx
   server {
       listen 443 ssl http2;
       server_name cloud.hypernexus.site;
       
       ssl_certificate /etc/nginx/ssl/cert.pem;
       ssl_certificate_key /etc/nginx/ssl/key.pem;
       
       # ... rest of configuration
   }
   ```

2. Add HTTP to HTTPS redirect:

   ```nginx
   server {
       listen 80;
       server_name cloud.hypernexus.site;
       return 301 https://$host$request_uri;
   }
   ```

3. Restart nginx:

   ```bash
   docker-compose restart nginx
   ```

## Testing SSL

### Check Certificate Details

```bash
# View certificate information
openssl x509 -in ssl/cert.pem -text -noout

# Check expiration date
openssl x509 -in ssl/cert.pem -enddate -noout
```

### Test HTTPS Connection

```bash
# Using curl
curl -I https://cloud.hypernexus.site

# Using openssl
openssl s_client -connect cloud.hypernexus.site:443 -servername cloud.hypernexus.site
```

### Online Tools

- [SSL Labs](https://www.ssllabs.com/ssltest/) - Comprehensive SSL test
- [SSL Shopper](https://www.sslshopper.com/ssl-checker.html) - Quick SSL check

## Troubleshooting

### Certificate Not Trusted

**Self-Signed Certificates:**

- This is expected for self-signed certificates
- Browsers will show a warning
- Click "Advanced" and "Proceed" to continue
- For production, use Let's Encrypt

**Let's Encrypt Certificates:**

- Ensure the domain points to your server
- Check that port 80 is accessible
- Verify certbot logs: `sudo certbot certificates`

### Connection Refused

1. Check if nginx is running:

   ```bash
   docker-compose ps
   ```

2. Check nginx logs:

   ```bash
   docker-compose logs nginx
   ```

3. Verify ports are open:

   ```bash
   sudo netstat -tlnp | grep -E ':(80|443)'
   ```

### Certificate Expired

1. Renew with certbot:

   ```bash
   sudo certbot renew
   ```

2. Copy new certificates:

   ```bash
   sudo cp /etc/letsencrypt/live/cloud.hypernexus.site/fullchain.pem ssl/cert.pem
   sudo cp /etc/letsencrypt/live/cloud.hypernexus.site/privkey.pem ssl/key.pem
   ```

3. Restart nginx:

   ```bash
   docker-compose restart nginx
   ```

## File Structure

```
ssl/
├── README.md              # This file
├── cert.pem               # SSL certificate
├── key.pem                # Private key
├── setup-ssl.sh           # Let's Encrypt setup script
├── generate-self-signed.sh # Self-signed certificate generator (Linux/Mac)
└── generate-windows.bat   # Self-signed certificate generator (Windows)
```

## Security Best Practices

1. **Keep private keys secure:**
   - Never commit `key.pem` to version control
   - Set restrictive permissions: `chmod 600 ssl/key.pem`

2. **Use strong certificates:**
   - Minimum 2048-bit RSA key
   - Use SHA-256 or stronger
   - Enable HSTS headers

3. **Regular renewal:**
   - Set up auto-renewal with certbot
   - Monitor certificate expiration
   - Test renewal process periodically

4. **Backup certificates:**
   - Store backups securely
   - Include in disaster recovery plan
