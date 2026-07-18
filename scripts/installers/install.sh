#!/bin/bash
set -e

echo ""
echo "========================================"
echo "  TormentNexus Installer v1.0.0"
echo "========================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Detect OS
OS="$(uname -s)"
case "${OS}" in
Linux*) MACHINE=Linux ;;
Darwin*) MACHINE=Mac ;;
*) MACHINE="UNKNOWN:${OS}" ;;
esac

echo "[*] Detected platform: ${MACHINE}"

# Set installation directory
INSTALL_DIR="$HOME/.local/bin"
CONFIG_DIR="$HOME/.tormentnexus"

echo "[*] Installing to user directory: ${INSTALL_DIR}"
mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"
mkdir -p "$CONFIG_DIR/memory"

# Download or copy binary
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/tormentnexus" ]; then
	echo "[*] Installing tormentnexus binary..."
	cp "$SCRIPT_DIR/tormentnexus" "$INSTALL_DIR/tormentnexus"
	chmod +x "$INSTALL_DIR/tormentnexus"
elif [ -f "$SCRIPT_DIR/tormentnexus-gui" ]; then
	echo "[*] Installing tormentnexus-gui binary..."
	cp "$SCRIPT_DIR/tormentnexus-gui" "$INSTALL_DIR/tormentnexus"
	chmod +x "$INSTALL_DIR/tormentnexus"
else
	echo -e "${RED}[!] No binary found in $SCRIPT_DIR${NC}"
	echo "    Expected: tormentnexus or tormentnexus-gui"
	exit 1
fi

# Create default configuration
echo "[*] Creating default configuration..."
cat >"$CONFIG_DIR/config.yaml" <<'EOF'
# TormentNexus Configuration
host: 127.0.0.1
port: 7778
workspace: ~/workspace/tormentnexus

# Memory Configuration
memory:
  l2_enabled: true
  l3_enabled: true
  l4_enabled: false

# Provider Configuration
providers:
  deepseek:
    enabled: true
    api_key: ""
  lmstudio:
    enabled: true
    url: http://127.0.0.1:1234
EOF

# Add to PATH if needed
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
	echo "[*] Adding $INSTALL_DIR to PATH..."
	if [ -f "$HOME/.bashrc" ]; then
		echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >>"$HOME/.bashrc"
	fi
	if [ -f "$HOME/.zshrc" ]; then
		echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >>"$HOME/.zshrc"
	fi
	export PATH="$PATH:$INSTALL_DIR"
fi

echo ""
echo "========================================"
echo -e "  ${GREEN}Installation Complete!${NC}"
echo "========================================"
echo ""
echo "TormentNexus has been installed to:"
echo "  $INSTALL_DIR/tormentnexus"
echo ""
echo "Configuration file:"
echo "  $CONFIG_DIR/config.yaml"
echo ""
echo "To start TormentNexus:"
echo "  tormentnexus serve"
echo ""
echo "Dashboard will be available at:"
echo "  http://127.0.0.1:7778"
echo ""

# Ask to start now
read -p "Start TormentNexus now? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
	echo "[*] Starting TormentNexus..."
	nohup "$INSTALL_DIR/tormentnexus" serve >/tmp/tormentnexus.log 2>&1 &
	sleep 2
	echo -e "${GREEN}[✓] TormentNexus started${NC}"
	echo "[*] Dashboard: http://127.0.0.1:7778"

	# Open browser
	if command -v xdg-open &>/dev/null; then
		xdg-open http://127.0.0.1:7778
	elif command -v open &>/dev/null; then
		open http://127.0.0.1:7778
	fi
fi
