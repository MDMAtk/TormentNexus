#!/usr/bin/env python3
"""Enrich catalog entries with DeepSeek — auto-categorize, extract capabilities, generate descriptions"""

import sqlite3
import json
import urllib.request
import time
import os

DB_PATH = "/opt/tormentnexus/catalog.db"
DEEPSEEK_KEY = os.environ.get("DEEPSEEK_API_KEY", "")
DEEPSEEK_URL = "https://api.deepseek.com/v1/chat/completions"


def call_deepseek(prompt, max_tokens=200):
    """Call DeepSeek API for enrichment"""
    try:
        data = json.dumps(
            {
                "model": "deepseek-chat",
                "messages": [{"role": "user", "content": prompt}],
                "max_tokens": max_tokens,
                "temperature": 0.3,
            }
        ).encode()

        req = urllib.request.Request(
            DEEPSEEK_URL,
            data=data,
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {DEEPSEEK_KEY}",
            },
        )

        with urllib.request.urlopen(req, timeout=30) as r:
            resp = json.loads(r.read().decode())
            return resp["choices"][0]["message"]["content"].strip()
    except Exception:
        return None


def enrich_entry(title, url, description):
    """Enrich a single catalog entry"""
    prompt = f"""Analyze this MCP server/tool and return a JSON object with:
- "summary": 1-2 sentence description of what it does
- "capabilities": list of 3-5 key capabilities
- "category": one of [database, filesystem, browser, search, api, devops, ai, communication, productivity, security, other]
- "install": npm/pip install command if obvious from the name

Title: {title}
URL: {url}
Description: {description or "none"}

Return ONLY valid JSON, no markdown."""

    result = call_deepseek(prompt)
    if not result:
        return None

    try:
        # Try to parse JSON from response
        result = result.strip()
        if result.startswith("```"):
            result = result.split("```")[1]
            if result.startswith("json"):
                result = result[4:]
        return json.loads(result)
    except:
        return None


def main():
    db = sqlite3.connect(DB_PATH)
    db.execute("PRAGMA journal_mode=WAL")

    # Get entries without enrichment
    rows = db.execute("""
        SELECT rowid, title, url, description 
        FROM links_backlog 
        WHERE description IS NULL OR description = '' OR length(description) < 20
        ORDER BY rowid
        LIMIT 100
    """).fetchall()

    print(f"Found {len(rows)} entries to enrich")

    enriched = 0
    for rowid, title, url, desc in rows:
        result = enrich_entry(title, url, desc)
        if result:
            new_desc = result.get("summary", "")
            capabilities = json.dumps(result.get("capabilities", []))
            category = result.get("category", "other")
            install = result.get("install", "")

            if new_desc:
                db.execute(
                    """
                    UPDATE links_backlog 
                    SET description = ?
                    WHERE rowid = ?
                """,
                    (new_desc, rowid),
                )
                enriched += 1

        time.sleep(0.5)  # Rate limit

        if enriched % 10 == 0 and enriched > 0:
            db.commit()
            print(f"  Enriched {enriched}/{len(rows)}")

    db.commit()
    db.close()

    print(f"\nEnriched {enriched} entries")


if __name__ == "__main__":
    main()
