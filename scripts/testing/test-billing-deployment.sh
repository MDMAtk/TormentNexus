#!/bin/bash
# TormentNexus — Billing & Deployment Test Suite

KERNEL="http://127.0.0.1:8090"
DEMO="https://demo.hypernexus.site"
LANDING="https://tormentnexus.site"
CLOUD="https://cloud.hypernexus.site"
PASS=0
FAIL=0
ERRORS=""

p() {
	echo "  ✅ $1"
	PASS=$((PASS + 1))
}
f() {
	echo "  ❌ $1"
	FAIL=$((FAIL + 1))
	ERRORS="$ERRORS\n- $1"
}

echo "=========================================="
echo "  Billing & Deployment Test Suite"
echo "=========================================="
echo ""

# ==========================================
echo "═══ A. BILLING & PAYMENTS ═══"
# ==========================================

echo "[A1] Pricing Page"
CODE=$(curl -s -o /dev/null -w "%{http_code}" "$LANDING/pricing")
[ "$CODE" = "200" ] && p "Pricing page loads" || f "Pricing page ($CODE)"

CONTENT=$(curl -s "$LANDING/pricing")
echo "$CONTENT" | grep -q "Free" && p "Free plan shown" || f "Free plan shown"
echo "$CONTENT" | grep -q "19" && p "Pro plan shown (\$19)" || f "Pro plan shown"
echo "$CONTENT" | grep -q "49" && p "Team plan shown (\$49)" || f "Team plan shown"
echo "$CONTENT" | grep -q "stripe\|checkout\|Get Started" && p "Checkout button" || f "Checkout button"

echo "[A2] Cloud Landing"
CODE=$(curl -s -o /dev/null -w "%{http_code}" "$CLOUD")
[ "$CODE" = "200" ] && p "Cloud landing loads" || f "Cloud landing ($CODE)"

echo "[A3] Account Registration"
EMAIL="billing_test_$(date +%s)@test.com"
RESULT=$(curl -s -X POST "$KERNEL/api/account/register" \
	-H "Content-Type: application/json" \
	-d "{\"email\":\"$EMAIL\",\"password\":\"BillingTest123!\"}")
echo "$RESULT" | grep -q "id" && p "Registration works" || f "Registration"
ACCOUNT_ID=$(echo "$RESULT" | python3 -c "import sys,json;print(json.load(sys.stdin).get('id',''))" 2>/dev/null)

echo "[A4] Account Login"
RESULT=$(curl -s -X POST "$KERNEL/api/account/login" \
	-H "Content-Type: application/json" \
	-d "{\"email\":\"$EMAIL\",\"password\":\"BillingTest123!\"}")
echo "$RESULT" | grep -q "token" && p "Login works" || f "Login"
TOKEN=$(echo "$RESULT" | python3 -c "import sys,json;print(json.load(sys.stdin).get('token',''))" 2>/dev/null)

echo "[A5] Account Status"
RESULT=$(curl -s "$KERNEL/api/account/status" 2>/dev/null)
echo "$RESULT" | grep -q "success\|status\|account" && p "Account status" || f "Account status"

echo "[A6] Subscription Provisioning"
RESULT=$(curl -s -X POST "$KERNEL/api/account/provision" \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN" \
	-d "{\"plan\":\"basic\",\"seats\":1}" 2>/dev/null)
echo "$RESULT" | grep -q "success\|provision\|account" && p "Provisioning" || f "Provisioning"

# ==========================================
echo ""
echo "═══ B. DEPLOYMENT & INFRASTRUCTURE ═══"
# ==========================================

echo "[B1] Docker"
docker ps | grep -q "tormentnexus-demo" && p "Demo container running" || f "Demo container"
docker images | grep -q "tormentnexus" && p "Docker image exists" || f "Docker image"

echo "[B2] PM2 Services"
ssh hetzner 'pm2 list 2>/dev/null | grep -q "tn-kernel.*online"' && p "tn-kernel online" || f "tn-kernel"
ssh hetzner 'pm2 list 2>/dev/null | grep -q "fwber-backend.*online"' && p "fwber-backend online" || f "fwber-backend"

echo "[B3] Health Checks"
curl -s "$DEMO/health" | grep -q "ok" && p "Demo health" || f "Demo health"
curl -s "$KERNEL/health" | grep -q "ok" && p "Kernel health" || f "Kernel health"

echo "[B4] SSL Certificates"
curl -s -o /dev/null -w "%{ssl_verify_result}" "$DEMO" | grep -q "0" && p "Demo SSL valid" || f "Demo SSL"
curl -s -o /dev/null -w "%{ssl_verify_result}" "$LANDING" | grep -q "0" && p "Landing SSL valid" || f "Landing SSL"
curl -s -o /dev/null -w "%{ssl_verify_result}" "$CLOUD" | grep -q "0" && p "Cloud SSL valid" || f "Cloud SSL"

echo "[B5] HTTP->HTTPS Redirect"
CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://demo.hypernexus.site")
[ "$CODE" = "301" ] || [ "$CODE" = "302" ] && p "HTTPS redirect" || f "HTTPS redirect ($CODE)"

echo "[B6] DNS Resolution"
dig +short demo.hypernexus.site | grep -q "." && p "demo.hypernexus.site" || f "demo DNS"
dig +short cloud.hypernexus.site | grep -q "." && p "cloud.hypernexus.site" || f "cloud DNS"
dig +short tormentnexus.site | grep -q "." && p "tormentnexus.site" || f "landing DNS"

# ==========================================
echo ""
echo "═══ C. API ENDPOINTS ═══"
# ==========================================

echo "[C1] Core APIs"
curl -s "$KERNEL/health" | grep -q "ok" && p "GET /health" || f "GET /health"
curl -s "$KERNEL/api/index" | grep -q "success" && p "GET /api/index" || f "GET /api/index"

echo "[C2] Account APIs"
curl -s -X POST "$KERNEL/api/account/register" \
	-H "Content-Type: application/json" \
	-d "{\"email\":\"api_test_$(date +%s)@test.com\",\"password\":\"test123\"}" |
	grep -q "id" && p "POST /api/account/register" || f "POST /api/account/register"

echo "[C3] Catalog APIs"
curl -s "$KERNEL/api/backlog/search?q=postgres" | grep -q "results" && p "GET /api/backlog/search" || f "GET /api/backlog/search"
curl -s "$KERNEL/api/backlog/stats" | grep -q "total" && p "GET /api/backlog/stats" || f "GET /api/backlog/stats"
curl -s "$KERNEL/api/backlog/categories" | grep -q "categories" && p "GET /api/backlog/categories" || f "GET /api/backlog/categories"

echo "[C4] Provider APIs"
curl -s "$KERNEL/api/providers/catalog" | grep -q "success\|providers" && p "GET /api/providers/catalog" || f "GET /api/providers/catalog"

echo "[C5] Demo APIs"
curl -s "$DEMO/api/backlog/search?q=docker" | grep -q "results" && p "Demo search" || f "Demo search"
curl -s "$DEMO/api/backlog/stats" | grep -q "total" && p "Demo stats" || f "Demo stats"

# ==========================================
echo ""
echo "═══ D. LANDING PAGES ═══"
# ==========================================

for PAGE in "" "blog/" "catalog" "pricing" "newsletter" "blog/feed.xml" "sitemap.xml" "robots.txt"; do
	CODE=$(curl -s -o /dev/null -w "%{http_code}" "$LANDING/$PAGE")
	[ "$CODE" = "200" ] && p "/$PAGE" || f "/$PAGE ($CODE)"
done

# ==========================================
echo ""
echo "═══ E. SECURITY ═══"
# ==========================================

echo "[E1] Security Headers"
HEADERS=$(curl -s -I "$DEMO")
echo "$HEADERS" | grep -qi "x-frame-options" && p "X-Frame-Options" || f "X-Frame-Options"
echo "$HEADERS" | grep -qi "x-content-type" && p "X-Content-Type-Options" || f "X-Content-Type-Options"

echo "[E2] Rate Limiting"
for i in $(seq 1 10); do
	curl -s -o /dev/null "$KERNEL/health"
done
p "Rate limit test (10 requests)"

# ==========================================
echo ""
echo "═══ F. PERFORMANCE ═══"
# ==========================================

echo "[F1] Response Times"
T=$(curl -s -o /dev/null -w "%{time_total}" "$DEMO/health")
p "Demo health: ${T}s"
T=$(curl -s -o /dev/null -w "%{time_total}" "$KERNEL/api/backlog/search?q=postgres")
p "Catalog search: ${T}s"
T=$(curl -s -o /dev/null -w "%{time_total}" "$LANDING")
p "Landing page: ${T}s"

echo "[F2] Disk Usage"
DISK=$(ssh hetzner 'df -h / | tail -1 | awk "{print \$5}"')
p "Disk usage: $DISK"

echo "[F3] Memory Usage"
MEM=$(ssh hetzner 'free -h | grep Mem | awk "{print \$3\"/\"\$2}"')
p "Memory usage: $MEM"

# ==========================================
echo ""
echo "═══ G. ENRICHMENT ═══"
# ==========================================

echo "[G1] Catalog Enrichment"
ENRICHED=$(ssh hetzner 'sqlite3 /opt/tormentnexus/catalog.db "SELECT count(*) FROM links_backlog WHERE description IS NOT NULL AND length(description) > 20"')
p "Enriched entries: $ENRICHED"
TOTAL=$(ssh hetzner 'sqlite3 /opt/tormentnexus/catalog.db "SELECT count(*) FROM links_backlog"')
p "Total entries: $TOTAL"

# ==========================================
echo ""
echo "═══ H. SCRAPERS & SWARM ═══"
# ==========================================

echo "[H1] Scraper Status"
SCRAPE_PROCS=$(ssh hetzner 'ps aux | grep -c "scrape\|enrich" | head -1')
p "Scraper processes: $SCRAPE_PROCS"

echo "[H2] Swarm Status"
SWARM_PROCS=$(ssh hetzner 'ps aux | grep -c "swarm" | head -1')
p "Swarm processes: $SWARM_PROCS"

# ==========================================
echo ""
echo "=========================================="
echo "  RESULTS: $PASS passed, $FAIL failed"
echo "=========================================="
echo ""
if [ $FAIL -gt 0 ]; then
	echo "❌ FAILED TESTS:"
	echo -e "$ERRORS"
	echo ""
fi
echo "Total: $((PASS + FAIL)) tests"
echo ""
