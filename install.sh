#!/bin/bash
# HyperNexus Installer — Works on Linux, macOS, and Windows (via Git Bash/WSL)
# Usage: curl -fsSL https://hypernexus.site/install.sh | bash

set -e

REPO="robertpelloni/HyperNexus"
VERSION="1.0.0"
BASE_URL="https://releases.hypernexus.site/${VERSION}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║   HyperNexus Universal Installer        ║${NC}"
echo -e "${CYAN}║   38 AI Clients • One Command           ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════╝${NC}"
echo ""

# Detect platform
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)
            echo -e "${RED}Error: Unsupported platform $(uname -s)${NC}"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)   arch="arm64" ;;
        *)
            echo -e "${RED}Error: Unsupported architecture $(uname -m)${NC}"
            exit 1
            ;;
    esac

    echo "${os}-${arch}"
}

# Get download URL
get_download_url() {
    local platform="$1"

    case "$platform" in
        windows-*)
            echo "${BASE_URL}/hypernexus-setup.exe"
            ;;
        darwin-arm64)
            echo "${BASE_URL}/hypernexus-darwin-arm64"
            ;;
        darwin-*)
            echo "${BASE_URL}/hypernexus-darwin-amd64"
            ;;
        linux-arm64)
            echo "${BASE_URL}/hypernexus-linux-arm64"
            ;;
        linux-*)
            echo "${BASE_URL}/hypernexus-linux-amd64"
            ;;
    esac
}

# Download file
download() {
    local url="$1"
    local dest="$2"

    if command -v curl &> /dev/null; then
        curl -fsSL -o "$dest" "$url"
    elif command -v wget &> /dev/null; then
        wget -qO "$dest" "$url"
    else
        echo -e "${RED}Error: curl or wget required${NC}"
        exit 1
    fi
}

# Main installation
main() {
    local platform url dest

    platform=$(detect_platform)
    url=$(get_download_url "$platform")

    echo -e "Platform: ${GREEN}${platform}${NC}"
    echo -e "Download: ${BLUE}${url}${NC}"
    echo ""

    # Create temp directory
    local tmpdir
    tmpdir=$(mktemp -d)
    trap 'rm -rf "$tmpdir"' EXIT

    # Determine filename
    if [[ "$platform" == windows-* ]]; then
        dest="${tmpdir}/hypernexus-setup.exe"
    else
        dest="${tmpdir}/hypernexus"
    fi

    # Download
    echo "Downloading HyperNexus..."
    download "$url" "$dest"
    echo -e "${GREEN}Download complete.${NC}"
    echo ""

    # Make executable on Unix
    if [[ "$platform" != windows-* ]]; then
        chmod +x "$dest"
    fi

    # Run installer
    echo "Running installer..."
    echo ""

    if [[ "$platform" == windows-* ]]; then
        # Check if running in Git Bash or WSL
        if command -v cmd.exe &> /dev/null; then
            # Convert to Windows path and run
            local win_dest
            win_dest=$(cygpath -w "$dest" 2>/dev/null || echo "$dest")
            cmd.exe /c "$win_dest" /S
        else
            echo -e "${RED}Error: Please run the .exe installer directly on Windows${NC}"
            echo "Download from: $url"
            exit 1
        fi
    else
        # Run binary installer
        "$dest"
    fi

    echo ""
    echo -e "${GREEN}✅ HyperNexus installed successfully!${NC}"
    echo ""
    echo "   Run 'hypernexus serve' to start the server."
    echo "   Run 'hypernexus --help' for more options."
    echo ""
    echo "   Or use npm: npx @hypernexus/install"
    echo ""
}

main "$@"
