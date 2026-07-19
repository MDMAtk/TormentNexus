#!/usr/bin/env python3
"""Bulk scrape MCP servers from awesome-mcp-servers and GitHub topics into catalog.db"""

import sqlite3
import json
import urllib.request
import urllib.error
import uuid
import time
import re

DB_PATH = "/opt/tormentnexus/catalog.db"


def fetch(url, timeout=15):
    try:
        req = urllib.request.Request(url, headers={"User-Agent": "TormentNexus/1.0"})
        with urllib.request.urlopen(req, timeout=timeout) as r:
            return r.read().decode("utf-8")
    except Exception as e:
        print(f"  WARN: {url} -> {e}")
        return None


def scrape_awesome_mcp():
    """Scrape github URLs from awesome-mcp-servers README"""
    print("=== Scraping awesome-mcp-servers ===")
    raw = fetch(
        "https://raw.githubusercontent.com/punkpeye/awesome-mcp-servers/main/README.md"
    )
    if not raw:
        return []

    # Extract [name](github-url) patterns
    pattern = re.compile(r"\[([^\]]+)\]\((https://github\.com/[^)]+)\)")
    matches = pattern.findall(raw)

    servers = []
    seen = set()
    for name, url in matches:
        url = url.rstrip("/")
        if url in seen or "awesome-mcp" in url.lower():
            continue
        seen.add(url)
        # Extract repo name from URL
        parts = url.replace("https://github.com/", "").split("/")
        if len(parts) >= 2:
            repo = f"{parts[0]}/{parts[1]}"
            servers.append(
                {
                    "name": name.strip(),
                    "url": url,
                    "repo": repo,
                    "source": "awesome-mcp-servers",
                }
            )

    print(f"  Found {len(servers)} unique servers")
    return servers


def scrape_github_topics():
    """Search GitHub for MCP server repos"""
    print("=== Scraping GitHub MCP topics ===")
    servers = []
    seen = set()

    for q in ["mcp+server", "mcp+server+topic:mcp", "model+context+protocol"]:
        raw = fetch(
            f"https://api.github.com/search/repositories?q={q}&sort=stars&order=desc&per_page=100"
        )
        if not raw:
            continue

        try:
            data = json.loads(raw)
        except:
            continue

        for item in data.get("items", []):
            url = item["html_url"]
            if url in seen:
                continue
            seen.add(url)
            servers.append(
                {
                    "name": item["full_name"],
                    "url": url,
                    "repo": item["full_name"],
                    "source": "github-search",
                    "description": item.get("description", ""),
                    "stars": item.get("stargazers_count", 0),
                }
            )

        time.sleep(1)  # Rate limit

    print(f"  Found {len(servers)} unique repos")
    return servers


def insert_into_db(servers):
    """Insert servers into catalog.db"""
    print(f"=== Inserting {len(servers)} servers into catalog.db ===")
    db = sqlite3.connect(DB_PATH)
    db.execute("PRAGMA journal_mode=WAL")

    now = int(time.time() * 1000)
    inserted = 0

    for s in servers:
        uid = str(uuid.uuid4())
        try:
            db.execute(
                """
                INSERT OR IGNORE INTO links_backlog (
                    uuid, url, normalized_url, title, description, tags, source,
                    is_duplicate, duplicate_of, research_status, http_status,
                    synced_at, created_at, updated_at
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
                (
                    uid,
                    s["url"],
                    s["url"],
                    s["name"],
                    s.get("description", f"MCP Server: {s['name']}"),
                    json.dumps(["mcp", "registry", s.get("source", "unknown")]),
                    s.get("source", "unknown"),
                    False,
                    "",
                    "pending",
                    200,
                    now,
                    now,
                    now,
                ),
            )
            inserted += 1
        except Exception:
            pass  # Duplicate or schema issue

    db.commit()

    # Count totals
    total = db.execute("SELECT count(*) FROM links_backlog").fetchone()[0]
    db.close()
    print(f"  Inserted: {inserted}, Total in backlog: {total}")


def main():
    servers = []
    servers.extend(scrape_awesome_mcp())
    servers.extend(scrape_github_topics())

    # Dedupe by URL
    seen = set()
    unique = []
    for s in servers:
        if s["url"] not in seen:
            seen.add(s["url"])
            unique.append(s)

    print(f"\nTotal unique servers: {len(unique)}")
    insert_into_db(unique)
    print("Done!")


if __name__ == "__main__":
    main()
