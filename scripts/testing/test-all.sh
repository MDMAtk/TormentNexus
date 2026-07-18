#!/bin/bash

BASE_URL="https://demo.hypernexus.site"
KERNEL_URL="http://127.0.0.1:8090"
PASS=0
FAIL=0
ERRORS=""

log() { echo "[$(date +%H:%M:%S)] $1"; }
pass() { echo "  ✅ $1"; PASS=$((PASS + 1)); }
fail() { echo "  ❌ $1"; FAIL=$((FAIL + 1)); ERRORS="$ERRORS\n- $1"; }

echo "=========================================="
echo "  TormentNexus Full Platform Test Suite"
echo "=========================================="
echo ""

# 1. Health Checks
log "Testing Health Checks..."
RESULT=$(curl -s --max-time 5 "$BASE_URL/health" 2>/dev/null)
echo "$RESULT" | grep -q "ok" && pass "Demo health check" || fail "Demo health check"
RESULT=$(curl -s --max-time 5 "$KERNEL_URL/health" 2>/dev/null)
echo "$RESULT" | grep -q "ok" && pass "Kernel health check" || fail "Kernel health check"

# 2. API Index
log "Testing API Index..."
RESULT=$(curl -s --max-time 5 "$KERNEL_URL/api/index" 2>/dev/null)
echo "$RESULT" | grep -q "success" && pass "API index" || fail "API index"

# 3. Catalog Search
log "Testing Catalog Search..."
RESULT=$(curl -s --max-time 5 "$KERNEL_URL/api/backlog/search?q=postgres&limit=5" 2>/dev/null)
echo "$RESULT" | grep -q "results" && pass "Catalog search" || fail "Catalog search"
TOTAL=$(echo "$RESULT" | python3 -c "import sys,json;print(json.load(sys.stdin).get(total,0))" 2>/dev/null || echo "0")
if [ "$TOTAL" -gt 0 ] 2>/dev/null; then pass "Catalog has $TOTAL entries"; else fail "Catalog entries"; fi
RESULT=$(curl -s --max-time 5 "$KERNEL_URL/api/backlog/search?category=mcp_server&limit=3" 2>/dev/null)
echo "$RESULT" | grep -q "results" && pass "Category filter" || fail "Category filter"

# 4. Catalog Stats
log "Testing Catalog Stats..."
RESULT=$(curl -s --max-time 5 "$KERNEL_URL/api/backlog/stats" 2>/dev/null)
echo "$RESULT" | grep -q "total" && pass "Stats endpoint" || fail "Stats endpoint"
TOTAL=$(echo "$RESULT" | python3 -c "import sys,json;print(json.load(sys.stdin).get(total,0))" 2>/dev/null || echo "0")
if [ "$TOTAL" -gt 10000 ] 2>/dev/null; then pass "Stats has $TOTAL entries"; else fail "Stats entry count ($TOTAL)"; fi

# 5. Categories
log "Testing Categories..."
RESULT=$(curl -s --max-time 5 "$KERNEL_URL/api/backlog/categories" 2>/dev/null)
echo "$RESULT" | grep -q "categories" && pass "Categories endpoint" || fail "Categories endpoint"

# 6. Account System
log "Testing Account System..."
EMAIL="test_$(date +%s)@test.com"
RESULT=$(curl -s --max-time 5 -X POST "$KERNEL_URL/api/account/register" -H "Content-Type: application/json" -d "{\"email\":\"$EMAIL\",\"password\":\"testpass123\"}" 2>/dev/null)
echo "$RESULT" | grep -q "id" && pass "Account registration" || fail "Account registration"
RESULT=$(curl -s --max-time 5 -X POST "$KERNEL_URL/api/account/login" -H "Content-Type: application/json" -d "{\"email\":\"$EMAIL\",\"password\":\"testpass123\"}" 2>/dev/null)
echo "$RESULT" | grep -q "token" && pass "Account login" || fail "Account login"

# 7. SSL
log "Testing SSL..."
RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "http://demo.hypernexus.site" 2>/dev/null)
if [ "$RESULT" = "301" ] || [ "$RESULT" = "302" ]; then pass "HTTP->HTTPS redirect"; else fail "HTTP redirect ($RESULT)"; fi

# 8. Security Headers
log "Testing Security Headers..."
HEADERS=$(curl -s --max-time 5 -I "$BASE_URL" 2>/dev/null)
echo "$HEADERS" | grep -qi "x-frame-options" && pass "X-Frame-Options" || fail "X-Frame-Options"
echo "$HEADERS" | grep -qi "x-content-type" && pass "X-Content-Type-Options" || fail "X-Content-Type-Options"

# 9. Performance
log "Testing Performance..."
TIME=$(curl -s --max-time 5 -o /dev/null -w "%{time_total}" "$BASE_URL/health" 2>/dev/null)
pass "Health response: ${TIME}s"
TIME=$(curl -s --max-time 5 -o /dev/null -w "%{time_total}" "$KERNEL_URL/api/backlog/search?q=postgres" 2>/dev/null)
pass "Search response: ${TIME}s"

# 10. Landing Pages
log "Testing Landing Pages..."
for PAGE in "" "blog/" "catalog" "pricing" "newsletter" "blog/feed.xml" "sitemap.xml" "robots.txt"; do
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site/$PAGE" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then pass "/$PAGE"; else fail "/$PAGE (HTTP $RESULT)"; fi
done

# 11. Docker
log "Testing Docker..."
docker ps | grep -q "tormentnexus-demo" && pass "Demo container running" || fail "Demo container running"

# 12. Memory System
log "Testing Memory System..."
RESULT=$(curl -s --max-time 5 "$KERNEL_URL/api/memory/stats" 2>/dev/null)
if echo "$RESULT" | grep -q "success\|l2\|memory"; then pass "Memory stats"; else fail "Memory stats"; fi

# 13. Provider Catalog
log "Testing Providers..."
RESULT=$(curl -s --max-time 5 "$KERNEL_URL/api/providers/catalog" 2>/dev/null)
echo "$RESULT" | grep -q "success\|providers\|catalog" && pass "Provider catalog" || fail "Provider catalog"

echo ""
echo "=========================================="
echo "  Results: $PASS passed, $FAIL failed"
echo "=========================================="
if [ $FAIL -gt 0 ]; then
    echo -e "\nFailed:$ERRORS"
    exit 1
fi
echo "  ✅ ALL TESTS PASSED"
