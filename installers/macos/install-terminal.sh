#!/bin/bash

# TormentNexus macOS/Linux Installer
# Terminal-based installer with progress UI

set -e

# ── Colors ──
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

# ── Configuration ──
REPO="MDMAtk/TormentNexus"
VERSION="1.0.0-b4"
INSTALL_DIR="$HOME/.tormentnexus"

# ── Detect Platform ──
detect_platform() {
	local os=$(uname -s)
	local arch=$(uname -m)

	case "$os" in
	Darwin)
		if [ "$arch" = "arm64" ]; then
			PLATFORM="darwin-arm64"
		else
			PLATFORM="darwin-amd64"
		fi
		;;
	Linux)
		if [ "$arch" = "aarch64" ] || [ "$arch" = "arm64" ]; then
			PLATFORM="linux-arm64"
		else
			PLATFORM="linux-amd64"
		fi
		;;
	*)
		echo -e "${RED}Unsupported OS: $os${NC}"
		exit 1
		;;
	esac
}

# ── Progress Bar ──
show_progress() {
	local current=$1
	local total=$2
	local width=50
	local percentage=$((current * 100 / total))
	local filled=$((current * width / total))
	local empty=$((width - filled))

	printf "\r  ["
	printf "%${filled}s" | tr ' ' '█'
	printf "%${empty}s" | tr ' ' '░'
	printf "] %3d%%" $percentage
}

# ── Log Function ──
log() {
	echo -e "  ${GREEN}✓${NC} $1"
}

log_info() {
	echo -e "  ${CYAN}→${NC} $1"
}

log_warn() {
	echo -e "  ${YELLOW}⚠${NC} $1"
}

log_error() {
	echo -e "  ${RED}✗${NC} $1"
}

# ── Download Function ──
download() {
	local url=$1
	local dest=$2

	if command -v curl &>/dev/null; then
		curl -L --progress-bar -o "$dest" "$url"
	elif command -v wget &>/dev/null; then
		wget -q --show-progress -O "$dest" "$url"
	else
		log_error "Neither curl nor wget found"
		exit 1
	fi
}

# ── Main Install Function ──
install() {
	clear

	echo ""
	echo -e "${CYAN}╔═══════════════════════════════════════════════════════════╗${NC}"
	echo -e "${CYAN}║                                                           ║${NC}"
	echo -e "${CYAN}║${NC}   ████████╗ ██████╗ ██████╗ ███╗   ███╗███████╗███╗   ██╗${CYAN}║${NC}"
	echo -e "${CYAN}║${NC}   ╚══██╔══╝██╔═══██╗██╔══██╗████╗ ████║██╔════╝████╗  ██║${CYAN}║${NC}"
	echo -e "${CYAN}║${NC}      ██║   ██║   ██║██████╔╝██╔████╔██║█████╗  ██╔██╗ ██║${CYAN}║${NC}"
	echo -e "${CYAN}║${NC}      ██║   ██║   ██║██╔══██╗██║╚██╔╝██║██╔══╝  ██║╚██╗██║${CYAN}║${NC}"
	echo -e "${CYAN}║${NC}      ██║   ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗██║ ╚████║${CYAN}║${NC}"
	echo -e "${CYAN}║${NC}      ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═══╝${CYAN}║${NC}"
	echo -e "${CYAN}║${NC}                  N E X U S   I N S T A L L E R           ${CYAN}║${NC}"
	echo -e "${CYAN}║                                                           ║${NC}"
	echo -e "${CYAN}╚═══════════════════════════════════════════════════════════╝${NC}"
	echo ""
	echo -e "  ${BOLD}AI Control Plane with Persistent Memory${NC}"
	echo -e "  ${DIM}26,000+ MCP Tools · Multi-Agent Orchestration${NC}"
	echo ""
	echo -e "${CYAN}─────────────────────────────────────────────────────────────${NC}"
	echo ""

	# Detect platform
	detect_platform
	log_info "Platform: ${BOLD}$PLATFORM${NC}"
	echo ""

	# Step 1: Create directories
	echo -e "  ${BOLD}[1/5]${NC} Creating installation directory..."
	mkdir -p "$INSTALL_DIR"
	mkdir -p "$INSTALL_DIR/.tormentnexus"
	mkdir -p "$INSTALL_DIR/logs"
	log "$INSTALL_DIR"
	show_progress 20 100
	echo ""
	echo ""

	# Step 2: Download binary
	echo -e "  ${BOLD}[2/5]${NC} Downloading TormentNexus..."
	local archive_name="tormentnexus-${PLATFORM}.tar.gz"
	local download_url="https://github.com/${REPO}/releases/download/${VERSION}/${archive_name}"
	local archive_path="$INSTALL_DIR/$archive_name"

	log_info "URL: $download_url"
	download "$download_url" "$archive_path"
	log "Downloaded"
	show_progress 50 100
	echo ""
	echo ""

	# Step 3: Extract
	echo -e "  ${BOLD}[3/5]${NC} Extracting files..."
	tar -xzf "$archive_path" -C "$INSTALL_DIR"
	rm -f "$archive_path"

	# Find and rename binary
	local binary=$(find "$INSTALL_DIR" -name "tormentnexus-*" -type f -executable 2>/dev/null | head -1)
	if [ -n "$binary" ]; then
		mv "$binary" "$INSTALL_DIR/tormentnexus"
		chmod +x "$INSTALL_DIR/tormentnexus"
	fi

	log "Extracted"
	show_progress 70 100
	echo ""
	echo ""

	# Step 4: Create config
	echo -e "  ${BOLD}[4/5]${NC} Creating configuration..."
	cat >"$INSTALL_DIR/.tormentnexus/config.json" <<'EOF'
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
	log "config.json created"
	show_progress 80 100
	echo ""
	echo ""

	# Step 5: Add to PATH
	echo -e "  ${BOLD}[5/5]${NC} Adding to PATH..."
	local shell_rc=""

	if [ -f "$HOME/.zshrc" ]; then
		shell_rc="$HOME/.zshrc"
	elif [ -f "$HOME/.bash_profile" ]; then
		shell_rc="$HOME/.bash_profile"
	elif [ -f "$HOME/.bashrc" ]; then
		shell_rc="$HOME/.bashrc"
	fi

	if [ -n "$shell_rc" ]; then
		if ! grep -q "$INSTALL_DIR" "$shell_rc"; then
			echo "" >>"$shell_rc"
			echo "# TormentNexus" >>"$shell_rc"
			echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >>"$shell_rc"
			log "Added to $shell_rc"
		else
			log "Already in PATH"
		fi
	else
		log_warn "Could not find shell RC file"
	fi

	# Add to current session
	export PATH="$PATH:$INSTALL_DIR"

	show_progress 100 100
	echo ""
	echo ""

	# Create launchd service (macOS)
	if [ "$(uname)" = "Darwin" ]; then
		local plist="$HOME/Library/LaunchAgents/com.tormentnexus.kernel.plist"
		cat >"$plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.tormentnexus.kernel</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/tormentnexus</string>
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
		log "LaunchAgent created"
	fi

	# Create systemd service (Linux)
	if [ "$(uname)" = "Linux" ]; then
		local service="$HOME/.config/systemd/user/tormentnexus.service"
		mkdir -p "$(dirname "$service")"
		cat >"$service" <<EOF
[Unit]
Description=TormentNexus AI Control Plane
After=network.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/tormentnexus serve
WorkingDirectory=$INSTALL_DIR
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
EOF
		systemctl --user daemon-reload 2>/dev/null || true
		log "Systemd service created"
	fi

	echo ""
	echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
	echo -e "${GREEN}║                                                           ║${NC}"
	echo -e "${GREEN}║${NC}            ${BOLD}INSTALLATION COMPLETE!${NC}                       ${GREEN}║${NC}"
	echo -e "${GREEN}║                                                           ║${NC}"
	echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
	echo ""
	echo -e "  ${BOLD}Location:${NC}  $INSTALL_DIR"
	echo -e "  ${BOLD}Binary:${NC}    $INSTALL_DIR/tormentnexus"
	echo -e "  ${BOLD}Config:${NC}    $INSTALL_DIR/.tormentnexus/config.json"
	echo ""
	echo -e "  ${BOLD}Quick Start:${NC}"
	echo -e "    ${CYAN}tormentnexus serve${NC}"
	echo ""
	echo -e "  ${BOLD}Dashboard:${NC}"
	echo -e "    ${CYAN}http://localhost:7778${NC}"
	echo ""
	echo -e "${CYAN}─────────────────────────────────────────────────────────────${NC}"
	echo ""

	# Ask to start
	read -p "  Start TormentNexus now? (y/n): " -n 1 -r
	echo ""

	if [[ $REPLY =~ ^[Yy]$ ]]; then
		echo ""
		log_info "Starting TormentNexus..."

		"$INSTALL_DIR/tormentnexus" serve >"$INSTALL_DIR/logs/stdout.log" 2>&1 &
		local pid=$!

		sleep 3

		if kill -0 $pid 2>/dev/null; then
			log "Server started (PID: $pid)"
			log_info "Dashboard: http://localhost:7778"

			# Open browser
			if command -v open &>/dev/null; then
				open "http://localhost:7778"
			elif command -v xdg-open &>/dev/null; then
				xdg-open "http://localhost:7778" 2>/dev/null
			fi
		else
			log_error "Failed to start. Check logs: $INSTALL_DIR/logs/"
		fi
	fi

	echo ""
	echo -e "  ${DIM}Press any key to exit...${NC}"
	read -n 1 -s
}

# Run installer
install
