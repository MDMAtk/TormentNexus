#!/usr/bin/env python3
"""
TormentNexus Reddit Scanner
Scans AI/ML subreddits, finds relevant threads, generates replies
Saves best replies to a queue for manual posting
"""

import os
import json
import time
import sqlite3
import requests
from datetime import datetime

# ── Configuration ──

SUBREDDITS = [
    "MachineLearning",
    "LocalLLaMA",
    "artificial",
    "singularity",
    "LangChain",
    "ChatGPT",
    "OpenAI",
    "LocalAI",
    "selfhosted",
    "opensource",
    "AI_Agents",
    "LLMDevs",
    "AIDevTools",
    "deeplearning",
    "ArtificialIntelligence",
]

RELEVANT_KEYWORDS = [
    "memory",
    "persistent memory",
    "context window",
    "forget",
    "mcp",
    "model context protocol",
    "tool use",
    "function calling",
    "local llm",
    "ollama",
    "lm studio",
    "localai",
    "self-hosted",
    "ai agent",
    "multi-agent",
    "agent framework",
    "ai control plane",
    "ai orchestration",
    "vector database",
    "embedding",
    "rag",
    "ai assistant",
    "ai coding",
    "ai developer",
    "open source ai",
    "ai tools",
]

HIGH_RELEVANCE_KEYWORDS = [
    "persistent memory",
    "mcp server",
    "model context protocol",
    "ai control plane",
    "ai that remembers",
    "memory for ai",
    "give ai memory",
    "ai memory system",
    "tool catalog",
    "mcp catalog",
    "mcp tools",
]

# MiMo API config
MIMO_KEY = os.environ.get("MIMO_API_KEY", "")
MIMO_URL = "https://token-plan-sgp.xiaomimimo.com/v1/chat/completions"
MIMO_MODEL = "mimo-v2.5"

# Paths
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
DATA_DIR = os.path.join(SCRIPT_DIR, "..", "data")
DB_PATH = os.path.join(DATA_DIR, "reddit-scanner.db")
QUEUE_FILE = os.path.join(DATA_DIR, "reddit-queue.json")

# Reply templates
REPLY_TEMPLATES = {
    "memory": """I built TormentNexus which has a 4-tier memory system (L1-L4) for exactly this:

- L1: Session memory (ephemeral)
- L2: Working memory (30 days, semantically indexed)
- L3: Cold archive (1 year, vector storage)
- L4: Limbo (soft delete, can resurrect)

It's open source and works with Ollama, LM Studio, and any OpenAI-compatible API.

GitHub: https://github.com/MDMAtk/TormentNexus""",
    "mcp": """I've cataloged 26,000+ MCP servers in TormentNexus. You can search them at https://tormentnexus.site/catalog

It includes tools for:
- Databases (PostgreSQL, SQLite, etc.)
- Filesystems
- Browsers (Playwright, Chrome DevTools)
- APIs (GitHub, Supabase, Vercel)
- And thousands more

GitHub: https://github.com/MDMAtk/TormentNexus""",
    "tools": """TormentNexus is an open-source AI control plane that does exactly this:

- 26,000+ MCP tools in one catalog
- Persistent memory (4-tier system)
- Multi-agent orchestration
- Unified dashboard
- Works with Ollama, LM Studio, DeepSeek

It's a Go backend with a Next.js dashboard. Single binary, runs locally.

GitHub: https://github.com/MDMAtk/TormentNexus""",
}


def init_db():
    """Initialize database"""
    os.makedirs(DATA_DIR, exist_ok=True)
    db = sqlite3.connect(DB_PATH)
    db.execute("""
        CREATE TABLE IF NOT EXISTS seen_posts (
            post_id TEXT PRIMARY KEY,
            subreddit TEXT,
            title TEXT,
            url TEXT,
            permalink TEXT,
            relevance_score INTEGER,
            replied BOOLEAN DEFAULT FALSE,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    """)
    db.execute("""
        CREATE TABLE IF NOT EXISTS replies (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id TEXT,
            subreddit TEXT,
            title TEXT,
            reply_text TEXT,
            included_link BOOLEAN,
            posted BOOLEAN DEFAULT FALSE,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    """)
    db.commit()
    return db


def call_mimo(prompt, max_tokens=500):
    """Call MiMo API for reply generation"""
    if not MIMO_KEY:
        return None

    try:
        response = requests.post(
            MIMO_URL,
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {MIMO_KEY}",
            },
            json={
                "model": MIMO_MODEL,
                "messages": [{"role": "user", "content": prompt}],
                "max_tokens": max_tokens,
                "temperature": 0.7,
            },
            timeout=60,
        )
        response.raise_for_status()
        result = response.json()
        return result["choices"][0]["message"]["content"].strip()
    except Exception as e:
        print(f"  [!] MiMo error: {e}")
        return None


def calculate_relevance(title, selftext=""):
    """Calculate relevance score"""
    text = (title + " " + selftext).lower()

    score = 0
    matched_keywords = []
    high_matches = []

    for keyword in RELEVANT_KEYWORDS:
        if keyword.lower() in text:
            score += 1
            matched_keywords.append(keyword)

    for keyword in HIGH_RELEVANCE_KEYWORDS:
        if keyword.lower() in text:
            score += 3
            high_matches.append(keyword)

    return score, matched_keywords, high_matches


def generate_reply(title, selftext, subreddit, include_link=False):
    """Generate a helpful reply"""
    text = (title + " " + selftext).lower()

    # Check for template matches
    if "memory" in text and (
        "forget" in text or "remember" in text or "persistent" in text
    ):
        return REPLY_TEMPLATES["memory"]

    if "mcp" in text and ("tool" in text or "server" in text or "catalog" in text):
        return REPLY_TEMPLATES["mcp"]

    # Use MiMo for custom replies
    link_section = ""
    if include_link:
        link_section = """
If genuinely relevant, include at the end:
"Full disclosure: I'm working on TormentNexus (https://github.com/MDMAtk/TormentNexus), an open-source AI control plane with persistent memory and 26K+ MCP tools."
Only include if it genuinely helps."""

    prompt = f"""You are a helpful AI/ML developer on r/{subreddit}.

Someone posted:
Title: {title}
Content: {selftext[:500] if selftext else "(no body)"}

Write a helpful reply:
1. Address their question directly
2. Provide useful technical insight
3. Be conversational, not salesy
4. 2-4 sentences max
5. Reddit formatting (no headers)
{link_section}
Be direct. Don't start with "Great question!" """

    return call_mimo(prompt, max_tokens=300)


def load_queue():
    """Load the reply queue"""
    if os.path.exists(QUEUE_FILE):
        with open(QUEUE_FILE) as f:
            return json.load(f)
    return []


def save_queue(queue):
    """Save the reply queue"""
    os.makedirs(DATA_DIR, exist_ok=True)
    with open(QUEUE_FILE, "w") as f:
        json.dump(queue, f, indent=2)


def main():
    """Main scanning loop"""
    print("=" * 60)
    print("  TormentNexus Reddit Scanner")
    print("=" * 60)
    print()
    print("  Mode: SCAN + GENERATE (manual posting)")
    print(f"  Monitoring {len(SUBREDDITS)} subreddits")
    print(f"  Queue file: {QUEUE_FILE}")
    print()

    db = init_db()

    cycle = 0
    while True:
        cycle += 1
        print(f"\n{'=' * 60}")
        print(f"  Cycle {cycle} — {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"{'=' * 60}")

        queue = load_queue()
        new_replies = 0

        for subreddit_name in SUBREDDITS:
            print(f"\n--- r/{subreddit_name} ---")

            try:
                # Fetch posts using Reddit's JSON API
                import urllib.request

                url = f"https://www.reddit.com/r/{subreddit_name}/new.json?limit=25"
                req = urllib.request.Request(
                    url,
                    headers={
                        "User-Agent": "TormentNexus/1.0 (https://github.com/MDMAtk/TormentNexus)"
                    },
                )

                try:
                    with urllib.request.urlopen(req, timeout=30) as r:
                        data = json.loads(r.read().decode())
                        posts = data.get("data", {}).get("children", [])
                except Exception as e:
                    print(f"  [!] Error fetching: {e}")
                    continue

                for post_data in posts:
                    d = post_data.get("data", {})
                    post_id = d.get("id", "")
                    title = d.get("title", "")
                    selftext = d.get("selftext", "")
                    permalink = d.get("permalink", "")

                    # Skip if already seen
                    existing = db.execute(
                        "SELECT post_id FROM seen_posts WHERE post_id = ?", (post_id,)
                    ).fetchone()
                    if existing:
                        continue

                    # Calculate relevance
                    score, keywords, high_matches = calculate_relevance(title, selftext)

                    # Save to database
                    db.execute(
                        """
                        INSERT OR IGNORE INTO seen_posts 
                        (post_id, subreddit, title, url, permalink, relevance_score)
                        VALUES (?, ?, ?, ?, ?, ?)
                    """,
                        (
                            post_id,
                            subreddit_name,
                            title,
                            d.get("url", ""),
                            permalink,
                            score,
                        ),
                    )
                    db.commit()

                    if score >= 2:
                        include_link = len(high_matches) > 0
                        print(f"  [+] Relevant ({score}): {title[:60]}...")
                        print(f"      Keywords: {', '.join(keywords[:3])}")

                        # Generate reply
                        reply = generate_reply(
                            title, selftext, subreddit_name, include_link
                        )

                        if reply:
                            print(f"      Reply: {reply[:80]}...")

                            # Save to database
                            db.execute(
                                """
                                INSERT INTO replies (post_id, subreddit, title, reply_text, included_link)
                                VALUES (?, ?, ?, ?, ?)
                            """,
                                (post_id, subreddit_name, title, reply, include_link),
                            )
                            db.commit()

                            # Add to queue
                            queue.append(
                                {
                                    "post_id": post_id,
                                    "subreddit": subreddit_name,
                                    "title": title,
                                    "url": f"https://reddit.com{permalink}",
                                    "reply": reply,
                                    "include_link": include_link,
                                    "score": score,
                                    "keywords": keywords,
                                    "timestamp": datetime.now().isoformat(),
                                    "posted": False,
                                }
                            )
                            new_replies += 1

                time.sleep(2)  # Rate limiting

            except Exception as e:
                print(f"  [!] Error with r/{subreddit_name}: {e}")
                continue

        # Save queue
        save_queue(queue)

        # Stats
        total_seen = db.execute("SELECT count(*) FROM seen_posts").fetchone()[0]
        total_relevant = db.execute(
            "SELECT count(*) FROM seen_posts WHERE relevance_score >= 2"
        ).fetchone()[0]
        total_replies = db.execute("SELECT count(*) FROM replies").fetchone()[0]
        pending = len([r for r in queue if not r.get("posted")])

        print("\n--- Stats ---")
        print(f"  Posts seen: {total_seen}")
        print(f"  Relevant: {total_relevant}")
        print(f"  Replies generated: {total_replies}")
        print(f"  New this cycle: {new_replies}")
        print(f"  Pending in queue: {pending}")

        if new_replies > 0:
            print(f"\n  [!] {new_replies} new replies ready to post!")
            print(f"  View queue: cat {QUEUE_FILE}")

        # Wait before next cycle
        print("\n  Waiting 5 minutes...")
        time.sleep(300)


if __name__ == "__main__":
    main()
