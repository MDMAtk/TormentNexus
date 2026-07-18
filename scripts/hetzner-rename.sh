#!/bin/bash
set -euo pipefail

# HyperNexus Hetzner Migration Script
# Renames TormentNexus -> HyperNexus without disturbing running services

echo "========================================="
echo "  HyperNexus Hetzner Migration"
echo "========================================="
echo ""

# Configuration
OLD_NAME="tormentnexus"
NEW_NAME="hypernexus"
OLD_GIT_URL="https://github.com/NexusSoftMDMA/TormentNexus.git"
NEW_GIT_URL="https://gitlab.com/robertpelloni/HyperNexus.git"
OLD_MODULE="github.com/MDMAtk/TormentNexus"
NEW_MODULE="gitlab.com/robertpelloni/HyperNexus"
WORK_DIR="/opt/tormentnexus"

# Safety check
echo "WARNING: This script will rename TormentNexus to HyperNexus"
echo "Running services will NOT be stopped, but will be renamed."
echo ""
echo "Work directory: $WORK_DIR"
echo ""

# Step 1: Update git remote
echo "Step 1: Updating git remote..."
cd "$WORK_DIR"
CURRENT_REMOTE=$(git remote get-url origin 2>/dev/null || echo "none")
echo "  Current remote: $CURRENT_REMOTE"
echo "  New remote: $NEW_GIT_URL"
git remote set-url origin "$NEW_GIT_URL"
echo "  ✓ Git remote updated"
echo ""

# Step 2: Update Go module path (if go.mod exists)
echo "Step 2: Updating Go module path..."
if [ -f go/go.mod ]; then
	sed -i "s|$OLD_MODULE|$NEW_MODULE|g" go/go.mod
	sed -i "s|$OLD_MODULE|$NEW_MODULE|g" go/go.sum 2>/dev/null || true
	echo "  ✓ Go module path updated"
else
	echo "  ⚠ go.mod not found, skipping"
fi
echo ""

# Step 3: Update start.sh
echo "Step 3: Updating start.sh..."
if [ -f start.sh ]; then
	cp start.sh start.sh.backup
	sed -i "s|TormentNexus TORMENTNEXUS|HyperNexus HYPERNEXUS|g" start.sh
	sed -i "s|github.com/MDMAtk/TormentNexus|$NEW_MODULE|g" start.sh
	sed -i "s|bin/tormentnexus|bin/hypernexus|g" start.sh
	echo "  ✓ start.sh updated"
else
	echo "  ⚠ start.sh not found, skipping"
fi
echo ""

# Step 4: Rename binary (if exists)
echo "Step 4: Checking binary..."
if [ -f "$WORK_DIR/tormentnexus" ]; then
	cp "$WORK_DIR/tormentnexus" "$WORK_DIR/hypernexus"
	echo "  ✓ Binary copied to hypernexus"
	echo "  ⚠ Old binary kept for compatibility"
else
	echo "  ⚠ Binary not found, skipping"
fi
echo ""

# Step 5: Create new service files (without stopping old ones)
echo "Step 5: Creating new service files..."

# Create hypernexus.service
cat >/etc/systemd/system/hypernexus.service <<'EOF'
[Unit]
Description=HyperNexus AI Agent
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/bin/script -q -c /usr/local/bin/hypernexus /dev/null
Restart=always
RestartSec=5
Environment=PORT=3000
WorkingDirectory=/opt/tormentnexus
StandardInput=null
StandardOutput=inherit
StandardError=inherit

[Install]
WantedBy=multi-user.target
EOF
echo "  ✓ Created hypernexus.service"

# Create hypernexus-bot.service
cat >/etc/systemd/system/hypernexus-bot.service <<'EOF'
[Unit]
Description=HyperNexus Autonomous Sales Bot
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
WorkingDirectory=/opt/tormentnexus/bot
ExecStart=/opt/tormentnexus/bot/sales_bot_linux
Restart=always
RestartSec=10
EnvironmentFile=/opt/tormentnexus/bot/.env
User=root

[Install]
WantedBy=multi-user.target
EOF
echo "  ✓ Created hypernexus-bot.service"

# Update marketing-agent.service description
if [ -f /etc/systemd/system/marketing-agent.service ]; then
	sed -i 's|TormentNexus Autonomous Marketing Agent|HyperNexus Autonomous Marketing Agent|g' /etc/systemd/system/marketing-agent.service
	echo "  ✓ Updated marketing-agent.service description"
fi
echo ""

# Step 6: Reload systemd (doesn't affect running services)
echo "Step 6: Reloading systemd..."
systemctl daemon-reload
echo "  ✓ Systemd reloaded"
echo ""

# Step 7: Update nginx configs (comments only, not server names)
echo "Step 7: Updating nginx config comments..."
if [ -f /etc/nginx/sites-enabled/hypernexus.site ]; then
	sed -i 's|TN deployment|HN deployment|g' /etc/nginx/sites-enabled/hypernexus.site
	echo "  ✓ Updated hypernexus.site comments"
fi
echo ""

# Step 8: Update .env files
echo "Step 8: Checking .env files..."
if [ -f "$WORK_DIR/.env" ]; then
	if grep -q "tormentnexus" "$WORK_DIR/.env"; then
		cp "$WORK_DIR/.env" "$WORK_DIR/.env.backup"
		sed -i 's|tormentnexus|hypernexus|g' "$WORK_DIR/.env"
		echo "  ✓ Updated .env references"
	else
		echo "  ✓ .env already clean"
	fi
fi
echo ""

# Step 9: Update AGENTS.md and README.md
echo "Step 9: Updating documentation..."
for f in AGENTS.md README.md CHANGELOG.md; do
	if [ -f "$WORK_DIR/$f" ]; then
		sed -i 's|TormentNexus|HyperNexus|g' "$WORK_DIR/$f"
		sed -i 's|tormentnexus|hypernexus|g' "$WORK_DIR/$f"
		echo "  ✓ Updated $f"
	fi
done
echo ""

# Step 10: Create symlink for backward compatibility
echo "Step 10: Creating compatibility symlink..."
if [ -f "$WORK_DIR/hypernexus" ] && [ ! -f "/usr/local/bin/hypernexus" ]; then
	ln -sf "$WORK_DIR/hypernexus" /usr/local/bin/hypernexus
	echo "  ✓ Created /usr/local/bin/hypernexus symlink"
fi
if [ -f "$WORK_DIR/tormentnexus" ] && [ ! -f "/usr/local/bin/tormentnexus" ]; then
	ln -sf "$WORK_DIR/tormentnexus" /usr/local/bin/tormentnexus
	echo "  ✓ Kept /usr/local/bin/tormentnexus symlink for compatibility"
fi
echo ""

# Step 11: Summary
echo "========================================="
echo "  Migration Summary"
echo "========================================="
echo ""
echo "✓ Git remote updated to GitLab"
echo "✓ Go module path updated"
echo "✓ start.sh updated"
echo "✓ New service files created (not started)"
echo "✓ Systemd reloaded"
echo "✓ Documentation updated"
echo ""
echo "IMPORTANT: Running services are UNCHANGED."
echo ""
echo "To activate the new services:"
echo "  1. Stop old services:"
echo "     systemctl stop tormentnexus tormentnexus-bot"
echo ""
echo "  2. Start new services:"
echo "     systemctl start hypernexus hypernexus-bot"
echo ""
echo "  3. Enable new services on boot:"
echo "     systemctl enable hypernexus hypernexus-bot"
echo ""
echo "  4. (Optional) Disable old services:"
echo "     systemctl disable tormentnexus tormentnexus-bot"
echo ""
echo "To pull latest code from GitLab:"
echo "  cd /opt/tormentnexus"
echo "  git pull origin main"
echo ""
