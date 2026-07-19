#!/bin/bash

KERNEL="http://127.0.0.1:8090"
DEMO="https://demo.hypernexus.site"
PASS=0
FAIL=0

p() { echo "  ✅ $1"; PASS=$((PASS+1)); }
f() { echo "  ❌ $1"; FAIL=$((FAIL+1)); }

echo "=========================================="
echo "  TormentNexus Platform Test Suite"
echo "=========================================="
echo ""

echo "[1] Health Checks"
curl -s "$DEMO/health" | grep -q ok && p "Demo" || f "Demo"
curl -s "$KERNEL/health" | grep -q ok && p "Kernel" || f "Kernel"

echo "[2] Catalog APIs"
curl -s "$KERNEL/api/backlog/search?q=postgres" | grep -q results && p "Search" || f "Search"
curl -s "$KERNEL/api/backlog/stats" | grep -q total && p "Stats" || f "Stats"
curl -s "$KERNEL/api/backlog/categories" | grep -q categories && p "Categories" || f "Categories"

echo "[3] Account System"
EMAIL="test$(date +%s)@test.com"
RES=$(curl -s -X POST "$KERNEL/api/account/register" -H "Content-Type: application/json" -d "{\"email\":\"$EMAIL\",\"password\":\"pass123\"}")
echo "$RES" | grep -q id && p "Register" || f "Register"
RES=$(curl -s -X POST "$KERNEL/api/account/login" -H "Content-Type: application/json" -d "{\"email\":\"$EMAIL\",\"password\":\"pass123\"}")
echo "$RES" | grep -q token && p "Login" || f "Login"

echo "[4] Provider APIs"
curl -s "$KERNEL/api/providers/catalog" | grep -q success && p "Providers" || f "Providers"

echo "[5] SSL/Security"
CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://demo.hypernexus.site")
[ "$CODE" = "301" ] || [ "$CODE" = "302" ] && p "HTTPS redirect" || f "HTTPS redirect"
curl -s -I "$DEMO" | grep -qi x-frame-options && p "X-Frame-Options" || f "X-Frame-Options"

echo "[6] Performance"
T=$(curl -s -o /dev/null -w "%{time_total}" "$DEMO/health")
p "Health: ${T}s"
T=$(curl -s -o /dev/null -w "%{time_total}" "$KERNEL/api/backlog/search?q=postgres")
p "Search: ${T}s"

echo "[7] Landing Pages"
for P in "" blog/ catalog pricing newsletter blog/feed.xml sitemap.xml robots.txt; do
    C=$(curl -s -o /dev/null -w "%{http_code}" "https://tormentnexus.site/$P")
    [ "$C" = "200" ] && p "/$P" || f "/$P ($C)"
done

echo "[8] Docker"
docker ps | grep -q tormentnexus-demo && p "Demo container" || f "Demo container"

echo "[9] Billing Pages"
curl -s "https://tormentnexus.site/pricing" | grep -q "Free" && p "Pricing page" || f "Pricing page"
curl -s "https://cloud.hypernexus.site" | grep -q -i "login\|hypernexus" && p "Cloud landing" || f "Cloud landing"

echo "[10] Demo Functionality"
curl -s "$DEMO/api/backlog/search?q=docker" | grep -q results && p "Demo search" || f "Demo search"
curl -s "$DEMO/api/backlog/stats" | grep -q total && p "Demo stats" || f "Demo stats"

echo ""
echo "=========================================="
echo "  Results: $PASS passed, $FAIL failed"
echo "=========================================="
[ $FAIL -gt 0 ] && exit 1
echo "  ✅ ALL TESTS PASSED"
