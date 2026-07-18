#!/bin/bash
# TormentNexus — Catalog DB Tracking & Versioning
# Tracks enrichment progress and creates snapshots

DB_PATH="/opt/tormentnexus/catalog.db"
TRACKING_DIR="/opt/tormentnexus/data/catalog-tracking"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p "$TRACKING_DIR"

echo "=========================================="
echo "  Catalog DB Status — $(date)"
echo "=========================================="
echo ""

# Get counts
TOTAL=$(sqlite3 "$DB_PATH" "SELECT count(*) FROM links_backlog;")
ENRICHED=$(sqlite3 "$DB_PATH" "SELECT count(*) FROM links_backlog WHERE description IS NOT NULL AND length(description) > 20;")
REMAINING=$((TOTAL - ENRICHED))
PERCENTAGE=$(echo "scale=1; $ENRICHED * 100 / $TOTAL" | bc)

echo "Total entries:     $TOTAL"
echo "Enriched:          $ENRICHED ($PERCENTAGE%)"
echo "Remaining:         $REMAINING"
echo ""

# Get counts by source
echo "=== By Source ==="
sqlite3 "$DB_PATH" "SELECT source, count(*) as cnt FROM links_backlog GROUP BY source ORDER BY cnt DESC LIMIT 10;" | while IFS='|' read -r source count; do
	printf "  %-25s %s\n" "$source" "$count"
done
echo ""

# Get counts by category
echo "=== By Category ==="
sqlite3 "$DB_PATH" "SELECT category, count(*) as cnt FROM links_backlog GROUP BY category ORDER BY cnt DESC LIMIT 10;" | while IFS='|' read -r category count; do
	printf "  %-25s %s\n" "$category" "$count"
done
echo ""

# Save tracking snapshot
SNAPSHOT_FILE="$TRACKING_DIR/snapshot_$TIMESTAMP.json"
cat >"$SNAPSHOT_FILE" <<EOF
{
  "timestamp": "$(date -Iseconds)",
  "total": $TOTAL,
  "enriched": $ENRICHED,
  "remaining": $REMAINING,
  "percentage": $PERCENTAGE
}
EOF

echo "Snapshot saved: $SNAPSHOT_FILE"
echo ""

# Show recent snapshots
echo "=== Recent Snapshots ==="
ls -la "$TRACKING_DIR"/snapshot_*.json 2>/dev/null | tail -5 | while read -r line; do
	echo "  $line"
done
echo ""

# Check if enrichment is running
if ps aux | grep -q "[p]ython3 /tmp/enrich.py"; then
	echo "✅ Enrichment is RUNNING"
else
	echo "❌ Enrichment is NOT running"
fi
echo ""

# Show latest enrichment log
if [ -f /tmp/enrich-run.log ]; then
	echo "=== Latest Enrichment Log ==="
	tail -10 /tmp/enrich-run.log
fi
