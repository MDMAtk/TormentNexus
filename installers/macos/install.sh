#!/bin/bash

# HyperNexus macOS Installer
# This script installs HyperNexus and sets it up as a service

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Banner
echo -e "${CYAN}"
echo "  ============================================================"
echo ""
echo "    ████████╗ ██████╗ ██████╗ ███╗   ███╗███████╗███╗   ██╗████████╗"
echo "    ╚══██╔══╝██╔═══██╗██╔══██╗████╗ ████║██╔════╝████╗  ██║╚══██╔══╝"
echo "       ██║   ██║   ██║██████╔╝██╔████╔██║█████╗  ██╔██╗ ██║   ██║"
echo "       ██║   ██║   ██║██╔══██╗██║╚██╔╝██║██╔══╝  ██║╚██╗██║   ██║"
echo "       ██║   ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗██║ ╚████║   ██║"
echo "       ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝"
echo "                   N E X U S   I N S T A L L E R"
echo ""
echo "  ============================================================"
echo -e "${NC}"
echo -e "  ${BLUE}AI Control Plane with Persistent Memory${NC}"
echo -e "  ${BLUE}26,000+ MCP Tools | Multi-Agent Orchestration${NC}"
echo ""
echo "  ============================================================"
echo ""

# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
	BINARY_NAME="hypernexus-darwin-arm64"
	echo -e "  ${GREEN}Detected: Apple Silicon (M1/M2)${NC}"
else
	BINARY_NAME="hypernexus-darwin-amd64"
	echo -e "  ${GREEN}Detected: Intel Mac${NC}"
fi
echo ""

# Set installation directory
INSTALL_DIR="$HOME/.hypernexus"
BINARY_PATH="$INSTALL_DIR/hypernexus"

echo "  [1/6] Creating installation directory..."
mkdir -p "$INSTALL_DIR"
echo -e "        ${GREEN}OK${NC} - $INSTALL_DIR"
echo ""

# Copy binary
echo "  [2/6] Installing HyperNexus binary..."
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp "$SCRIPT_DIR/$BINARY_NAME" "$BINARY_PATH"
chmod +x "$BINARY_PATH"
echo -e "        ${GREEN}OK${NC} - $BINARY_PATH"
echo ""

# Create config directory
echo "  [3/6] Creating configuration directory..."
mkdir -p "$INSTALL_DIR/.hypernexus"
echo -e "        ${GREEN}OK${NC} - $INSTALL_DIR/.hypernexus"
echo ""

# Create default config
echo "  [4/6] Creating default configuration..."
cat >"$INSTALL_DIR/.hypernexus/config.json" <<'EOF'
{
  "version": "1.0.0",
  "server": {
    "host": "127.0.0.1",
    "port": 7778
  },
  "memory": {
    "enabled": true,
    "tiers": ["L1", "L2", "L3", "L4"]
  },
  "mcp": {
    "catalog": true,
    "autoInstall": false
  }
}
EOF
echo -e "        ${GREEN}OK${NC} - config.json created"
echo ""

# Add to PATH
echo "  [5/6] Adding to PATH..."
SHELL_RC=""
if [ -f "$HOME/.zshrc" ]; then
	SHELL_RC="$HOME/.zshrc"
elif [ -f "$HOME/.bash_profile" ]; then
	SHELL_RC="$HOME/.bash_profile"
elif [ -f "$HOME/.bashrc" ]; then
	SHELL_RC="$HOME/.bashrc"
fi

if [ -n "$SHELL_RC" ]; then
	if ! grep -q "$INSTALL_DIR" "$SHELL_RC"; then
		echo "" >>"$SHELL_RC"
		echo "# HyperNexus" >>"$SHELL_RC"
		echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >>"$SHELL_RC"
		echo -e "        ${GREEN}OK${NC} - Added to $SHELL_RC"
	else
		echo -e "        ${GREEN}OK${NC} - Already in PATH"
	fi
else
	echo -e "        ${YELLOW}WARN${NC} - Could not find shell RC file"
	echo "              Add $INSTALL_DIR to your PATH manually"
fi
echo ""

# Create launchd service (macOS)
echo "  [6/6] Creating launchd service..."
PLIST_FILE="$HOME/Library/LaunchAgents/com.hypernexus.kernel.plist"
cat >"$PLIST_FILE" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.hypernexus.kernel</string>
    <key>ProgramArguments</key>
    <array>
        <string>$BINARY_PATH</string>
        <string>serve</string>
    </array>
    <key>RunAtLoad</key>
    <false/>
    <key>KeepAlive</key>
    <false/>
    <key>WorkingDirectory</key>
    <string>$INSTALL_DIR</string>
    <key>StandardOutPath</key>
    <string>$INSTALL_DIR/logs/stdout.log</string>
    <key>StandardErrorPath</key>
    <string>$INSTALL_DIR/logs/stderr.log</string>
</dict>
</plist>
EOF

mkdir -p "$INSTALL_DIR/logs"
echo -e "        ${GREEN}OK${NC} - Service created"
echo ""

echo "  ============================================================"
echo ""
echo -e "  ${GREEN}INSTALLATION COMPLETE!${NC}"
echo ""
echo "  ============================================================"
echo ""
echo "  HyperNexus has been installed to:"
echo "    $INSTALL_DIR"
echo ""
echo "  To start HyperNexus:"
echo "    1. Open a new terminal"
echo "    2. Run: hypernexus serve"
echo "    3. Open: http://localhost:7778"
echo ""
echo "  Or start the background service:"
echo "    launchctl load $PLIST_FILE"
echo ""
echo "  ============================================================"
echo ""

# Ask to start now
read -p "  Start HyperNexus now? (y/n): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
	echo ""
	echo "  Starting HyperNexus..."
	mkdir -p "$INSTALL_DIR/logs"
	"$BINARY_PATH" serve >"$INSTALL_DIR/logs/stdout.log" 2>&1 &
	SERV_PID=$!
	sleep 2
	if kill -0 $SERV_PID 2>/dev/null; then
		echo -e "  ${GREEN}HyperNexus is running (PID: $SERV_PID)${NC}"
		echo "  Dashboard: http://localhost:7778"
		open http://localhost:7778
	else
		echo -e "  ${RED}Failed to start. Check logs: $INSTALL_DIR/logs/${NC}"
	fi
fi

echo ""
echo "  Press any key to exit..."
read -n 1 -s
