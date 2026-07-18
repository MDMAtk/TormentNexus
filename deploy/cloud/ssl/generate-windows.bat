@echo off
setlocal

echo =========================================
echo   SSL Certificate Generator for Windows
echo =========================================
echo.
echo NOTE: Self-signed certificates are for TESTING only!
echo       Browsers will show security windows.
echo       Use Let's Encrypt for production certificates.
echo.

set DOMAIN=cloud.hypernexus.site
set DAYS=365

REM Create ssl directory if it doesn't exist
if not exist ssl mkdir ssl

REM Check if OpenSSL is available
where openssl >nul 2>nul
if errorlevel 1 (
    echo OpenSSL not found!
    echo.
    echo Options:
    echo 1. Install OpenSSL: https://slproweb.com/products/Win32OpenSSL.html
    echo 2. Use Git Bash which includes OpenSSL
    echo 3. Use WSL (Windows Subsystem for Linux)
    echo.
    echo After installing OpenSSL, add it to your PATH and try again.
    pause
    exit /b 1
)

echo Generating self-signed certificate for %DOMAIN%...
echo.

REM Generate private key
openssl genrsa -out ssl\key.pem 2048

REM Generate certificate signing request
openssl req -new -key ssl\key.pem -out ssl\csr.pem -subj "/C=US/ST=State/L=City/O=HyperNexus/CN=%DOMAIN%"

REM Generate self-signed certificate
openssl x509 -req -days %DAYS% -in ssl\csr.pem -signkey ssl\key.pem -out ssl\cert.pem -extfile <(echo subjectAltName=DNS:%DOMAIN%,DNS:www.%DOMAIN%)

REM Clean up CSR
del ssl\csr.pem 2>nul

echo.
echo =========================================
echo   Certificate Generated!
echo =========================================
echo.
echo Files created:
echo   - ssl\cert.pem (certificate)
echo   - ssl\key.pem (private key)
echo.
echo Validity: %DAYS% days
echo.
echo To use with docker-compose:
echo   cd deploy\cloud
echo   docker-compose up -d
echo.
echo To test:
echo   curl -k https://%DOMAIN%
echo.
pause
