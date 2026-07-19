#!/usr/bin/env python3
"""Continue scraping more sources"""

import sqlite3
import json
import urllib.request
import uuid
import time
import re

DB_PATH = "/opt/tormentnexus/catalog.db"


def fetch(url, timeout=20):
    try:
        req = urllib.request.Request(url, headers={"User-Agent": "TormentNexus/1.0"})
        with urllib.request.urlopen(req, timeout=timeout) as r:
            return r.read().decode("utf-8")
    except:
        return None


db = sqlite3.connect(DB_PATH)
existing = set(r[0] for r in db.execute("SELECT normalized_url FROM links_backlog"))
now = int(time.time() * 1000)
total_inserted = 0

# 1. GitHub: mcp by language
print("[1] GitHub: mcp by language...")
inserted = 0
for lang in [
    "python",
    "typescript",
    "go",
    "rust",
    "java",
    "ruby",
    "swift",
    "c++",
    "c#",
]:
    raw = fetch(
        f"https://api.github.com/search/repositories?q=mcp+language:{lang}&sort=stars&order=desc&per_page=100"
    )
    if not raw:
        time.sleep(4)
        continue
    data = json.loads(raw)
    for item in data.get("items", []):
        url = item["html_url"].rstrip("/")
        if url in existing:
            continue
        existing.add(url)
        src = f"github-lang-{lang}"
        try:
            db.execute(
                "INSERT OR IGNORE INTO links_backlog (uuid,url,normalized_url,title,description,tags,source,is_duplicate,duplicate_of,research_status,http_status,synced_at,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
                (
                    str(uuid.uuid4()),
                    url,
                    url,
                    item["full_name"],
                    item.get("description", "") or "",
                    json.dumps(["mcp", src]),
                    src,
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
        except:
            pass
    time.sleep(3)
total_inserted += inserted
print(f"  +{inserted}")

# 2. GitHub: MCP categories
print("[2] GitHub: MCP categories...")
inserted = 0
cats = [
    "database",
    "filesystem",
    "browser",
    "search",
    "github",
    "slack",
    "discord",
    "notion",
    "docker",
    "aws",
    "stripe",
    "email",
    "calendar",
    "monitoring",
    "testing",
    "security",
    "payment",
    "crm",
    "cms",
    "ecommerce",
    "social",
    "analytics",
    "logging",
    "auth",
]
for cat in cats:
    raw = fetch(
        f"https://api.github.com/search/repositories?q=mcp+{cat}&sort=stars&order=desc&per_page=50"
    )
    if not raw:
        time.sleep(4)
        continue
    data = json.loads(raw)
    for item in data.get("items", []):
        url = item["html_url"].rstrip("/")
        if url in existing:
            continue
        existing.add(url)
        src = f"github-cat-{cat}"
        try:
            db.execute(
                "INSERT OR IGNORE INTO links_backlog (uuid,url,normalized_url,title,description,tags,source,is_duplicate,duplicate_of,research_status,http_status,synced_at,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
                (
                    str(uuid.uuid4()),
                    url,
                    url,
                    item["full_name"],
                    item.get("description", "") or "",
                    json.dumps(["mcp", src]),
                    src,
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
        except:
            pass
    time.sleep(3)
total_inserted += inserted
print(f"  +{inserted}")

# 3. npm: AI SDK packages
print("[3] npm: AI SDK packages...")
inserted = 0
for q in [
    "@anthropic-ai",
    "@openai",
    "@google-ai",
    "@mistralai",
    "@huggingface",
    "openai+sdk",
    "anthropic+sdk",
    "cohere+ai",
]:
    raw = fetch(f"https://registry.npmjs.org/-/v1/search?text={q}&size=250")
    if not raw:
        continue
    data = json.loads(raw)
    for obj in data.get("objects", []):
        pkg = obj.get("package", {})
        name = pkg.get("name", "")
        desc = pkg.get("description", "")
        url = (
            pkg.get("links", {})
            .get("repository", f"https://www.npmjs.com/package/{name}")
            .rstrip("/")
        )
        if url in existing:
            continue
        existing.add(url)
        try:
            db.execute(
                "INSERT OR IGNORE INTO links_backlog (uuid,url,normalized_url,title,description,tags,source,is_duplicate,duplicate_of,research_status,http_status,synced_at,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
                (
                    str(uuid.uuid4()),
                    url,
                    url,
                    name,
                    desc,
                    json.dumps(["tool", "npm-ai-sdk"]),
                    "npm-ai-sdk",
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
        except:
            pass
    time.sleep(1)
total_inserted += inserted
print(f"  +{inserted}")

# 4. GitHub: awesome-mcp forks
print("[4] GitHub: awesome-mcp forks...")
inserted = 0
raw = fetch(
    "https://api.github.com/search/repositories?q=awesome+mcp+in:name+fork:true&sort=updated&order=desc&per_page=100"
)
if raw:
    data = json.loads(raw)
    for item in data.get("items", []):
        repo = item["full_name"]
        readme = fetch(f"https://raw.githubusercontent.com/{repo}/main/README.md")
        if not readme:
            continue
        for name, url in re.findall(
            r"\[([^\]]+)\]\((https://github\.com/[^)]+)\)", readme
        ):
            url = url.rstrip("/")
            if url in existing:
                continue
            existing.add(url)
            try:
                db.execute(
                    "INSERT OR IGNORE INTO links_backlog (uuid,url,normalized_url,title,description,tags,source,is_duplicate,duplicate_of,research_status,http_status,synced_at,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
                    (
                        str(uuid.uuid4()),
                        url,
                        url,
                        name.strip(),
                        "",
                        json.dumps(["mcp", "awesome-fork-v2"]),
                        "awesome-fork-v2",
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
            except:
                pass
        time.sleep(2)
total_inserted += inserted
print(f"  +{inserted}")

# 5. GitHub: more awesome lists
print("[5] GitHub: more awesome lists...")
inserted = 0
repos = [
    "f/awesome-chatgpt-prompts",
    "PraisonAI/PraisonAI",
    "langgenius/dify",
    "lobehub/lobe-chat",
    "BerriAI/litellm",
    "mlabonne/llm-course",
    "openai/openai-cookbook",
    "anthropics/anthropic-cookbook",
    "microsoft/promptflow",
    "langchain-ai/langchain",
    "run-llama/llama_index",
    "hwchase17/langchain",
    "jerryjliu/llama_index",
]
for repo in repos:
    readme = fetch(f"https://raw.githubusercontent.com/{repo}/main/README.md")
    if not readme:
        readme = fetch(f"https://raw.githubusercontent.com/{repo}/master/README.md")
    if not readme:
        continue
    for name, url in re.findall(r"\[([^\]]+)\]\((https://github\.com/[^)]+)\)", readme):
        url = url.rstrip("/")
        if url in existing:
            continue
        existing.add(url)
        try:
            db.execute(
                "INSERT OR IGNORE INTO links_backlog (uuid,url,normalized_url,title,description,tags,source,is_duplicate,duplicate_of,research_status,http_status,synced_at,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
                (
                    str(uuid.uuid4()),
                    url,
                    url,
                    name.strip(),
                    "",
                    json.dumps(["tool", "awesome-ai-v2"]),
                    "awesome-ai-v2",
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
        except:
            pass
    time.sleep(1)
total_inserted += inserted
print(f"  +{inserted}")

db.commit()
total = db.execute("SELECT count(*) FROM links_backlog").fetchone()[0]
db.close()
print(f"\nTotal new: {total_inserted}")
print(f"Grand total: {total}")
