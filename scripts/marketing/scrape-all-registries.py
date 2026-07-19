#!/usr/bin/env python3
"""Comprehensive MCP server scraper — all sources"""

import sqlite3
import json
import urllib.request
import urllib.error
import uuid
import time
import re

DB_PATH = "/opt/tormentnexus/catalog.db"


def fetch(url, timeout=20):
    try:
        req = urllib.request.Request(url, headers={"User-Agent": "TormentNexus/1.0"})
        with urllib.request.urlopen(req, timeout=timeout) as r:
            return r.read().decode("utf-8")
    except Exception as e:
        print(f"  WARN: {url} -> {e}")
        return None


def insert_servers(db, servers, source):
    """Insert servers into links_backlog, deduped by URL"""
    now = int(time.time() * 1000)
    inserted = 0
    seen_urls = set()

    # Get existing URLs to avoid duplicates
    for row in db.execute("SELECT normalized_url FROM links_backlog"):
        seen_urls.add(row[0])

    for s in servers:
        url = s.get("url", "").rstrip("/")
        if not url or url in seen_urls:
            continue
        seen_urls.add(url)

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
                    url,
                    url,
                    s.get("name", "unknown"),
                    s.get("description", f"MCP Server: {s.get('name', 'unknown')}"),
                    json.dumps(["mcp", source]),
                    source,
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
            pass

    return inserted


def scrape_awesome_mcp():
    """Scrape from awesome-mcp-servers README"""
    print("=== awesome-mcp-servers ===")
    raw = fetch(
        "https://raw.githubusercontent.com/punkpeye/awesome-mcp-servers/main/README.md"
    )
    if not raw:
        return []

    # Extract all [name](github-url) patterns
    pattern = re.compile(r"\[([^\]]+)\]\((https://github\.com/[^)]+)\)")
    matches = pattern.findall(raw)

    servers = []
    seen = set()
    for name, url in matches:
        url = url.rstrip("/")
        if url in seen or "awesome-mcp" in url.lower() or "discord" in url.lower():
            continue
        seen.add(url)
        parts = url.replace("https://github.com/", "").split("/")
        if len(parts) >= 2:
            servers.append(
                {"name": name.strip(), "url": url, "repo": f"{parts[0]}/{parts[1]}"}
            )

    print(f"  Found {len(servers)} unique servers")
    return servers


def scrape_github_topics():
    """Search GitHub for MCP server repos"""
    print("=== GitHub Topics ===")
    servers = []
    seen = set()

    queries = [
        "mcp+server+topic:mcp",
        "mcp+server+language:python",
        "mcp+server+language:typescript",
        "model+context+protocol",
        "mcp-server+in:name",
        "mcp+tools+in:description",
    ]

    for q in queries:
        raw = fetch(
            f"https://api.github.com/search/repositories?q={q}&sort=stars&order=desc&per_page=100"
        )
        if not raw:
            time.sleep(2)
            continue

        try:
            data = json.loads(raw)
        except Exception:
            time.sleep(2)
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
                    "description": item.get("description", ""),
                    "stars": item.get("stargazers_count", 0),
                }
            )

        time.sleep(2)  # Rate limit

    print(f"  Found {len(servers)} unique repos")
    return servers


def scrape_npm_mcp():
    """Search npm for MCP server packages"""
    print("=== npm MCP packages ===")
    servers = []

    queries = [
        "mcp-server",
        "mcp+server",
        "model-context-protocol",
        "@modelcontextprotocol",
    ]

    for q in queries:
        raw = fetch(f"https://registry.npmjs.org/-/v1/search?text={q}&size=250")
        if not raw:
            time.sleep(1)
            continue

        try:
            data = json.loads(raw)
        except Exception:
            time.sleep(1)
            continue

        for obj in data.get("objects", []):
            pkg = obj.get("package", {})
            name = pkg.get("name", "")
            desc = pkg.get("description", "")
            homepage = pkg.get("homepage", "")
            repo = pkg.get("links", {}).get("repository", "")

            # Filter: must be MCP-related
            keywords = pkg.get("keywords", [])
            is_mcp = any(
                k in ["mcp", "mcp-server", "model-context-protocol"] for k in keywords
            )
            if not is_mcp and "mcp" not in name.lower() and "mcp" not in desc.lower():
                continue

            url = repo or homepage or f"https://www.npmjs.com/package/{name}"
            servers.append({"name": name, "url": url, "description": desc})

        time.sleep(1)

    print(f"  Found {len(servers)} npm packages")
    return servers


def scrape_pypi_mcp():
    """Search PyPI for MCP server packages"""
    print("=== PyPI MCP packages ===")
    servers = []

    # PyPI simple API doesn't support search, use JSON API
    raw = fetch("https://pypi.org/simple/", timeout=30)
    if not raw:
        return []

    # Extract package names containing "mcp"
    mcp_packages = re.findall(r"<a[^>]*>([^<]*mcp[^<]*)</a>", raw, re.IGNORECASE)

    # Get details for top packages
    count = 0
    for pkg_name in mcp_packages[:500]:  # Limit to 500
        if count >= 200:
            break

        detail = fetch(f"https://pypi.org/pypi/{pkg_name}/json", timeout=5)
        if not detail:
            continue

        try:
            data = json.loads(detail)
            info = data.get("info", {})
            servers.append(
                {
                    "name": info.get("name", pkg_name),
                    "url": info.get("project_url")
                    or info.get("home_page")
                    or f"https://pypi.org/project/{pkg_name}/",
                    "description": info.get("summary", ""),
                }
            )
            count += 1
        except Exception:
            continue

        time.sleep(0.5)

    print(f"  Found {len(servers)} PyPI packages")
    return servers


def scrape_github_curated_lists():
    """Scrape other curated MCP lists"""
    print("=== Curated MCP lists ===")
    servers = []

    lists = [
        "https://raw.githubusercontent.com/punkpeye/awesome-mcp-clients/main/README.md",
        "https://raw.githubusercontent.com/wong2/awesome-mcp-servers/main/README.md",
        "https://raw.githubusercontent.com/jmandel/awesome-mcp-servers/main/README.md",
    ]

    for url in lists:
        raw = fetch(url)
        if not raw:
            continue

        pattern = re.compile(r"\[([^\]]+)\]\((https://github\.com/[^)]+)\)")
        matches = pattern.findall(raw)

        seen = set()
        for name, gh_url in matches:
            gh_url = gh_url.rstrip("/")
            if gh_url in seen:
                continue
            seen.add(gh_url)
            servers.append({"name": name.strip(), "url": gh_url})

        time.sleep(1)

    print(f"  Found {len(servers)} from curated lists")
    return servers


def main():
    db = sqlite3.connect(DB_PATH)
    db.execute("PRAGMA journal_mode=WAL")

    total_inserted = 0

    # Scrape all sources
    all_servers = []
    all_servers.extend(scrape_awesome_mcp())
    all_servers.extend(scrape_github_topics())
    all_servers.extend(scrape_npm_mcp())
    all_servers.extend(scrape_pypi_mcp())
    all_servers.extend(scrape_github_curated_lists())

    # Dedupe by URL
    seen = set()
    unique = []
    for s in all_servers:
        url = s.get("url", "").rstrip("/")
        if url and url not in seen:
            seen.add(url)
            unique.append(s)

    print(f"\nTotal unique servers: {len(unique)}")

    # Insert all
    inserted = insert_servers(db, unique, "multi-source-scrape")
    db.commit()

    # Count totals
    total = db.execute("SELECT count(*) FROM links_backlog").fetchone()[0]
    by_source = db.execute(
        "SELECT source, count(*) FROM links_backlog GROUP BY source ORDER BY count(*) DESC"
    ).fetchall()

    db.close()

    print(f"\nInserted: {inserted}")
    print(f"Total in backlog: {total}")
    print("\nBy source:")
    for src, cnt in by_source:
        print(f"  {src}: {cnt}")


if __name__ == "__main__":
    main()
