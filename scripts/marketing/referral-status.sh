#!/bin/bash
# TormentNexus — Referral System
# Track referrals and give credits

DB="/opt/tormentnexus/catalog.db"
REFDB="/opt/tormentnexus/data/referrals.db"

mkdir -p /opt/tormentnexus/data

# Create referrals table if not exists
sqlite3 "$REFDB" <<"SQL"
CREATE TABLE IF NOT EXISTS referrals (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  referrer_id TEXT NOT NULL,
  referred_id TEXT,
  referral_code TEXT UNIQUE NOT NULL,
  status TEXT DEFAULT 'pending',
  credits INTEGER DEFAULT 10,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  completed_at DATETIME
);

CREATE TABLE IF NOT EXISTS credits (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id TEXT NOT NULL,
  amount INTEGER NOT NULL,
  reason TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
SQL

echo "=== Referral System Status ==="
echo ""

echo "Total referrals:"
sqlite3 "$REFDB" "SELECT count(*) FROM referrals;"

echo ""
echo "Pending referrals:"
sqlite3 "$REFDB" "SELECT count(*) FROM referrals WHERE status='pending';"

echo ""
echo "Completed referrals:"
sqlite3 "$REFDB" "SELECT count(*) FROM referrals WHERE status='completed';"

echo ""
echo "Total credits given:"
sqlite3 "$REFDB" "SELECT COALESCE(sum(amount), 0) FROM credits;"

echo ""
echo "Recent referrals:"
sqlite3 "$REFDB" "SELECT referral_code, status, created_at FROM referrals ORDER BY created_at DESC LIMIT 5;"
