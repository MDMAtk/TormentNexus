#!/usr/bin/env python3
"""Mega MCP scraper — exhaustively scrape all available registries"""

import sqlite3
import json
import urllib.request
import urllib.error
import uuid
import time
import re

DB_PATH = "/opt/tormentnexus/catalog.db"


def fetch(url, timeout=25):
    try:
        req = urllib.request.Request(
            url,
            headers={
                "User-Agent": "TormentNexus/1.0",
                "Accept": "application/json, text/html, */*",
            },
        )
        with urllib.request.urlopen(req, timeout=timeout) as r:
            return r.read().decode("utf-8")
    except Exception:
        return None


def get_existing(db):
    return set(r[0] for r in db.execute("SELECT normalized_url FROM links_backlog"))


def insert_batch(db, servers, source, existing):
    now = int(time.time() * 1000)
    inserted = 0
    for s in servers:
        url = s.get("url", "").rstrip("/")
        if not url or url in existing or len(url) > 500:
            continue
        existing.add(url)
        try:
            db.execute(
                """INSERT OR IGNORE INTO links_backlog 
                (uuid,url,normalized_url,title,description,tags,source,is_duplicate,duplicate_of,research_status,http_status,synced_at,created_at,updated_at)
                VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)""",
                (
                    str(uuid.uuid4()),
                    url,
                    url,
                    s.get("name", "")[:200],
                    s.get("description", "")[:500],
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


# ──────────────────────────────────────────────
# SOURCE SCRAPERS
# ──────────────────────────────────────────────


def scrape_awesome_mcp(db, existing):
    print("[1/10] awesome-mcp-servers...")
    raw = fetch(
        "https://raw.githubusercontent.com/punkpeye/awesome-mcp-servers/main/README.md"
    )
    if not raw:
        return 0
    servers = []
    for name, url in re.findall(r"\[([^\]]+)\]\((https://github\.com/[^)]+)\)", raw):
        url = url.rstrip("/")
        if "awesome-mcp" in url.lower() or "discord" in url.lower():
            continue
        servers.append({"name": name.strip(), "url": url})
    n = insert_batch(db, servers, "awesome-mcp-servers", existing)
    print(f"  +{n}")
    return n


def scrape_github_search(db, existing):
    print("[2/10] GitHub search (20 queries)...")
    queries = [
        "mcp+server+topic:mcp",
        "mcp+server+language:python",
        "mcp+server+language:typescript",
        "mcp+server+language:go",
        "mcp+server+language:rust",
        "mcp+server+language:java",
        "model-context-protocol",
        "mcp-server+in:name",
        "mcp+tools+in:description",
        "mcp+database+server",
        "mcp+filesystem+server",
        "mcp+browser+server",
        "mcp+api+server",
        "mcp+web+server",
        "mcp+search+server",
        "mcp+github+server",
        "mcp+slack+server",
        "mcp+discord+server",
        "mcp+docker+server",
        "mcp+kubernetes+server",
    ]
    inserted = 0
    for q in queries:
        raw = fetch(
            f"https://api.github.com/search/repositories?q={q}&sort=stars&order=desc&per_page=100"
        )
        if not raw:
            time.sleep(4)
            continue
        try:
            data = json.loads(raw)
        except Exception:
            time.sleep(4)
            continue
        servers = []
        for item in data.get("items", []):
            servers.append(
                {
                    "name": item["full_name"],
                    "url": item["html_url"],
                    "description": item.get("description", "") or "",
                }
            )
        inserted += insert_batch(db, servers, "github-search", existing)
        time.sleep(3)
    print(f"  +{inserted}")
    return inserted


def scrape_npm(db, existing):
    print("[3/10] npm (10 queries)...")
    queries = [
        "mcp-server",
        "mcp+server",
        "model-context-protocol",
        "@modelcontextprotocol/server",
        "mcp+tool",
        "mcp+client",
        "mcp+proxy",
        "mcp+gateway",
        "mcp+aggregator",
        "mcp+bridge",
    ]
    inserted = 0
    for q in queries:
        raw = fetch(f"https://registry.npmjs.org/-/v1/search?text={q}&size=250")
        if not raw:
            time.sleep(1)
            continue
        try:
            data = json.loads(raw)
        except Exception:
            continue
        servers = []
        for obj in data.get("objects", []):
            pkg = obj.get("package", {})
            name = pkg.get("name", "")
            desc = pkg.get("description", "")
            keywords = pkg.get("keywords", [])
            is_mcp = any(
                k in ["mcp", "mcp-server", "model-context-protocol"] for k in keywords
            )
            if not is_mcp and "mcp" not in name.lower() and "mcp" not in desc.lower():
                continue
            repo = pkg.get("links", {}).get("repository", "")
            homepage = pkg.get("homepage", "")
            url = repo or homepage or f"https://www.npmjs.com/package/{name}"
            servers.append({"name": name, "url": url, "description": desc})
        inserted += insert_batch(db, servers, "npm", existing)
        time.sleep(1)
    print(f"  +{inserted}")
    return inserted


def scrape_pypi(db, existing):
    print("[4/10] PyPI (mcp packages)...")
    # Use PyPI JSON search
    inserted = 0
    for q in ["mcp-server", "mcp-server-", "model-context-protocol"]:
        raw = fetch(f"https://pypi.org/pypi/{q}/json", timeout=10)
        if raw:
            try:
                data = json.loads(raw)
                info = data.get("info", {})
                servers = [
                    {
                        "name": info.get("name", ""),
                        "url": info.get("project_url", "")
                        or f"https://pypi.org/project/{info.get('name', '')}/",
                        "description": info.get("summary", ""),
                    }
                ]
                inserted += insert_batch(db, servers, "pypi", existing)
            except Exception:
                pass
        time.sleep(0.5)

    # Search via Google for PyPI MCP packages
    raw = fetch("https://pypi.org/simple/")
    if raw:
        mcp_pkgs = re.findall(
            r"<a[^>]*>([^<]*(?:mcp|model.context.protocol)[^<]*)</a>",
            raw,
            re.IGNORECASE,
        )
        servers = []
        for pkg in mcp_pkgs[:300]:
            servers.append(
                {
                    "name": pkg,
                    "url": f"https://pypi.org/project/{pkg}/",
                    "description": f"PyPI package: {pkg}",
                }
            )
        inserted += insert_batch(db, servers, "pypi", existing)
    print(f"  +{inserted}")
    return inserted


def scrape_curated_lists(db, existing):
    print("[5/10] Curated awesome lists...")
    urls = [
        (
            "https://raw.githubusercontent.com/wong2/awesome-mcp-servers/main/README.md",
            "wong2-awesome",
        ),
        (
            "https://raw.githubusercontent.com/AI-Engineer-Foundation/awesome-mcp-servers/main/README.md",
            "aief-awesome",
        ),
        (
            "https://raw.githubusercontent.com/michaellatman/mcp-get/main/README.md",
            "mcp-get",
        ),
        (
            "https://raw.githubusercontent.com/taketwo/awesome-mcp-servers/main/README.md",
            "taketwo-awesome",
        ),
        (
            "https://raw.githubusercontent.com/HollyAdequateAI/awesome-mcp-servers/main/README.md",
            "holly-awesome",
        ),
    ]
    inserted = 0
    for url, source in urls:
        raw = fetch(url)
        if not raw:
            continue
        servers = []
        for name, gh_url in re.findall(
            r"\[([^\]]+)\]\((https://github\.com/[^)]+)\)", raw
        ):
            servers.append({"name": name.strip(), "url": gh_url.rstrip("/")})
        inserted += insert_batch(db, servers, source, existing)
        time.sleep(1)
    print(f"  +{inserted}")
    return inserted


def scrape_glama_html(db, existing):
    print("[6/10] Glama.ai web directory...")
    raw = fetch("https://glama.ai/mcp/servers", timeout=30)
    if not raw:
        print("  FAILED")
        return 0
    servers = []
    for m in re.finditer(r'href="(/mcp/servers/[^"]+)"', raw):
        slug = m.group(1).split("/")[-1]
        servers.append({"name": slug, "url": f"https://glama.ai{m.group(1)}"})
    # Also extract any github links
    for name, url in re.findall(r"\[([^\]]+)\]\((https://github\.com/[^)]+)\)", raw):
        servers.append({"name": name.strip(), "url": url.rstrip("/")})
    n = insert_batch(db, servers, "glama-html", existing)
    print(f"  +{n}")
    return n


def scrape_smithery(db, existing):
    print("[7/10] Smithery.ai...")
    inserted = 0
    for page in range(1, 6):
        raw = fetch(f"https://smithery.ai/api/v1/servers?page={page}&limit=100")
        if not raw:
            break
        try:
            data = json.loads(raw)
        except Exception:
            break
        servers = []
        for srv in data.get("servers", data.get("results", [])):
            servers.append(
                {
                    "name": srv.get("name", ""),
                    "url": srv.get("url", "") or srv.get("homepage", ""),
                    "description": srv.get("description", ""),
                }
            )
        inserted += insert_batch(db, servers, "smithery", existing)
        time.sleep(1)
    print(f"  +{inserted}")
    return inserted


def scrape_docker_hub(db, existing):
    print("[8/10] Docker Hub MCP images...")
    raw = fetch(
        "https://hub.docker.com/v2/search/repositories/?query=mcp+server&page_size=100"
    )
    if not raw:
        print("  FAILED")
        return 0
    try:
        data = json.loads(raw)
    except Exception:
        print("  FAILED")
        return 0
    servers = []
    for repo in data.get("results", []):
        name = repo.get("name", "")
        namespace = repo.get("namespace", "")
        desc = repo.get("description", "")
        url = f"https://hub.docker.com/r/{namespace}/{name}"
        servers.append({"name": f"{namespace}/{name}", "url": url, "description": desc})
    n = insert_batch(db, servers, "docker-hub", existing)
    print(f"  +{n}")
    return n


def scrape_crates_io(db, existing):
    print("[9/10] crates.io Rust MCP packages...")
    raw = fetch("https://crates.io/api/v1/crates?q=mcp+server&per_page=100")
    if not raw:
        print("  FAILED")
        return 0
    try:
        data = json.loads(raw)
    except Exception:
        print("  FAILED")
        return 0
    servers = []
    for crate in data.get("crates", []):
        name = crate.get("name", "")
        desc = crate.get("description", "")
        url = f"https://crates.io/crates/{name}"
        servers.append({"name": name, "url": url, "description": desc})
    n = insert_batch(db, servers, "crates-io", existing)
    print(f"  +{n}")
    return n


def scrape_mcp_hub(db, existing):
    print("[10/10] MCP Hub / directories...")
    inserted = 0
    # Try various MCP directories
    directories = [
        "https://mcp-hub.com/api/servers",
        "https://mcpfinder.ai/api/servers",
        "https://mcpservers.org/api/servers",
    ]
    for url in directories:
        raw = fetch(url)
        if not raw:
            continue
        try:
            data = json.loads(raw)
        except Exception:
            continue
        servers = []
        items = (
            data
            if isinstance(data, list)
            else data.get("servers", data.get("results", []))
        )
        for srv in items:
            if isinstance(srv, dict):
                servers.append(
                    {
                        "name": srv.get("name", ""),
                        "url": srv.get("url", "") or srv.get("homepage", ""),
                        "description": srv.get("description", ""),
                    }
                )
        inserted += insert_batch(db, servers, "mcp-hub", existing)
        time.sleep(1)
    print(f"  +{inserted}")
    return inserted


# ──────────────────────────────────────────────
# MAIN
# ──────────────────────────────────────────────


def main():
    db = sqlite3.connect(DB_PATH)
    db.execute("PRAGMA journal_mode=WAL")
    db.execute("PRAGMA busy_timeout=10000")
    existing = get_existing(db)

    print(f"Starting with {len(existing)} existing entries\n")

    total_new = 0
    total_new += scrape_awesome_mcp(db, existing)
    db.commit()
    total_new += scrape_github_search(db, existing)
    db.commit()
    total_new += scrape_npm(db, existing)
    db.commit()
    total_new += scrape_pypi(db, existing)
    db.commit()
    total_new += scrape_curated_lists(db, existing)
    db.commit()
    total_new += scrape_glama_html(db, existing)
    db.commit()
    total_new += scrape_smithery(db, existing)
    db.commit()
    total_new += scrape_docker_hub(db, existing)
    db.commit()
    total_new += scrape_crates_io(db, existing)
    db.commit()
    total_new += scrape_mcp_hub(db, existing)
    db.commit()

    total = db.execute("SELECT count(*) FROM links_backlog").fetchone()[0]
    by_source = db.execute(
        "SELECT source, count(*) FROM links_backlog GROUP BY source ORDER BY count(*) DESC"
    ).fetchall()
    db.close()

    print(f"\n{'=' * 50}")
    print(f"NEW: {total_new}")
    print(f"TOTAL: {total}")
    print("\nBy source:")
    for src, cnt in by_source:
        print(f"  {src}: {cnt}")


if __name__ == "__main__":
    main()
