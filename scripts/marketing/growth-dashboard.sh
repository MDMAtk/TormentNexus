#!/bin/bash
# TormentNexus — Growth Dashboard
# Track all growth metrics in one place

echo "=========================================="
echo "  TormentNexus Growth Dashboard"
echo "  $(date)"
echo "=========================================="
echo ""

# GitHub Stars
STARS=$(curl -s https://api.github.com/repos/MDMAtk/TormentNexus | grep -o '"stargazers_count": [0-9]*' | grep -o '[0-9]*')
echo "⭐ GitHub Stars: $STARS"

# Demo Health
HEALTH=$(curl -s --max-time 5 https://demo.hypernexus.site/health | grep -o '"ok":true' || echo "down")
echo "🟢 Demo Status: $HEALTH"

# Catalog Stats
STATS=$(curl -s --max-time 5 https://demo.hypernexus.site/api/backlog/stats 2>/dev/null)
TOTAL=$(echo "$STATS" | grep -o '"total":[0-9]*' | grep -o '[0-9]*')
echo "📚 Catalog Entries: $TOTAL"

# Enrichment Progress
ENRICHED=$(ssh hetzner 'sqlite3 /opt/tormentnexus/catalog.db "SELECT count(*) FROM links_backlog WHERE description IS NOT NULL AND length(description) > 20;"' 2>/dev/null)
echo "✨ Enriched Entries: $ENRICHED"

# Landing Pages
echo ""
echo "📄 Landing Pages:"
for PAGE in "" "blog/" "catalog" "pricing" "newsletter"; do
	CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://tormentnexus.site/$PAGE" 2>/dev/null)
	echo "  /$PAGE: $CODE"
done

# API Endpoints
echo ""
echo "🔌 API Endpoints:"
for ENDPOINT in "health" "api/index" "api/backlog/search?q=test" "api/backlog/stats"; do
	CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://demo.hypernexus.site/$ENDPOINT" 2>/dev/null)
	echo "  /$ENDPOINT: $CODE"
done

# Docker Status
echo ""
echo "🐳 Docker:"
docker ps | grep tormentnexus | wc -l | xargs -I {} echo "  Running containers: {}"

# PM2 Services
echo ""
echo "⚙️ PM2 Services:"
ssh hetzner 'pm2 list 2>/dev/null | grep -E "name|online" | head -5' 2>/dev/null

# Recent Activity
echo ""
echo "📊 Recent Activity:"
echo "  Last commit: $(git log -1 --format='%h %s' 2>/dev/null)"
echo "  Commits today: $(git log --since='today' --oneline 2>/dev/null | wc -l)"

echo ""
echo "=========================================="
