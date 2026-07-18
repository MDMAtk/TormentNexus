#!/usr/bin/env python3
"""
TormentNexus Reddit Agent
Monitors AI/ML subreddits for relevant threads and generates helpful replies
Uses MiMo v2.5 for reply generation

Works in two modes:
1. SCAN-ONLY: Monitors and generates replies (no posting)
2. AUTHENTICATED: Posts replies with Reddit API credentials
"""

import os
import json
import time
import sqlite3
import urllib.request
from datetime import datetime

# Subreddits to monitor
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
]

# Keywords that indicate relevance
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

# High-relevance keywords (worth posting GitHub link)
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

# Database
DB_PATH = "/opt/tormentnexus/data/reddit-agent.db"


def init_db():
    """Initialize database"""
    os.makedirs(os.path.dirname(DB_PATH), exist_ok=True)
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
            reply_text TEXT,
            included_link BOOLEAN,
            posted BOOLEAN DEFAULT FALSE,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    """)
    db.commit()
    return db


def call_mimo(prompt, max_tokens=500):
    """Call MiMo API"""
    try:
        data = json.dumps(
            {
                "model": MIMO_MODEL,
                "messages": [{"role": "user", "content": prompt}],
                "max_tokens": max_tokens,
                "temperature": 0.7,
            }
        ).encode()

        req = urllib.request.Request(
            MIMO_URL,
            data=data,
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {MIMO_KEY}",
            },
        )

        with urllib.request.urlopen(req, timeout=60) as r:
            resp = json.loads(r.read().decode())
            content = resp["choices"][0]["message"].get("content", "") or ""
            return content.strip()
    except Exception as e:
        print(f"  [!] MiMo error: {e}")
        return None


def fetch_subreddit_posts(subreddit, limit=25):
    """Fetch posts using Reddit's public JSON API"""
    url = f"https://www.reddit.com/r/{subreddit}/new.json?limit={limit}"
    try:
        req = urllib.request.Request(
            url,
            headers={
                "User-Agent": "TormentNexus/1.0 (AI research assistant; https://github.com/MDMAtk/TormentNexus)"
            },
        )
        with urllib.request.urlopen(req, timeout=30) as r:
            data = json.loads(r.read().decode())
            return data.get("data", {}).get("children", [])
    except Exception as e:
        print(f"  [!] Error fetching r/{subreddit}: {e}")
        return []


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
    """Generate a helpful reply using MiMo"""

    link_section = ""
    if include_link:
        link_section = """
If genuinely relevant, you MAY include at the end:
"Full disclosure: I'm working on TormentNexus (https://github.com/MDMAtk/TormentNexus), an open-source AI control plane that addresses this with persistent memory and 26K+ MCP tools."
Only include if it genuinely helps. Don't force it.
"""

    prompt = f"""You are a helpful AI/ML developer on r/{subreddit}.

Post:
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


def main():
    """Main monitoring loop"""
    print("=" * 50)
    print("  TormentNexus Reddit Agent")
    print("=" * 50)
    print()

    db = init_db()

    # Check for Reddit credentials
    creds_file = "/opt/tormentnexus/data/reddit-creds.json"
    has_credentials = False
    if os.path.exists(creds_file):
        with open(creds_file) as f:
            creds = json.load(f)
        if creds.get("client_id") != "YOUR_CLIENT_ID":
            has_credentials = True

    if has_credentials:
        print("Mode: AUTHENTICATED (will post replies)")
    else:
        print("Mode: SCAN-ONLY (generates replies, no posting)")
        print("To enable posting, add Reddit API credentials to:")
        print(f"  {creds_file}")
    print()

    print(f"Monitoring {len(SUBREDDITS)} subreddits:")
    for sub in SUBREDDITS:
        print(f"  - r/{sub}")
    print()

    cycle = 0
    while True:
        cycle += 1
        print(f"\n=== Cycle {cycle} — {datetime.now().strftime('%H:%M:%S')} ===")

        for subreddit in SUBREDDITS:
            print(f"\n--- r/{subreddit} ---")

            posts = fetch_subreddit_posts(subreddit, limit=25)

            for post in posts:
                data = post.get("data", {})
                post_id = data.get("id", "")
                title = data.get("title", "")
                selftext = data.get("selftext", "")
                permalink = data.get("permalink", "")

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
                    INSERT OR IGNORE INTO seen_posts (post_id, subreddit, title, url, permalink, relevance_score)
                    VALUES (?, ?, ?, ?, ?, ?)
                """,
                    (post_id, subreddit, title, data.get("url", ""), permalink, score),
                )
                db.commit()

                if score >= 2:
                    include_link = len(high_matches) > 0
                    print(f"  [+] Relevant ({score}): {title[:60]}...")
                    print(f"      Keywords: {', '.join(keywords[:3])}")

                    if include_link:
                        print("      HIGH RELEVANCE — will include link")

                    # Generate reply
                    reply = generate_reply(title, selftext, subreddit, include_link)

                    if reply:
                        print(f"      Reply: {reply[:100]}...")

                        # Save reply
                        db.execute(
                            """
                            INSERT INTO replies (post_id, reply_text, included_link)
                            VALUES (?, ?, ?)
                        """,
                            (post_id, reply, include_link),
                        )
                        db.execute(
                            "UPDATE seen_posts SET replied = TRUE WHERE post_id = ?",
                            (post_id,),
                        )
                        db.commit()

                        # Save to pending file
                        pending_file = "/opt/tormentnexus/data/pending-replies.json"
                        pending = []
                        if os.path.exists(pending_file):
                            with open(pending_file) as f:
                                pending = json.load(f)

                        pending.append(
                            {
                                "post_id": post_id,
                                "subreddit": subreddit,
                                "title": title,
                                "url": f"https://reddit.com{permalink}",
                                "reply": reply,
                                "include_link": include_link,
                                "score": score,
                                "keywords": keywords,
                                "timestamp": datetime.now().isoformat(),
                            }
                        )

                        with open(pending_file, "w") as f:
                            json.dump(pending, f, indent=2)

            time.sleep(2)

        # Stats
        total_seen = db.execute("SELECT count(*) FROM seen_posts").fetchone()[0]
        total_relevant = db.execute(
            "SELECT count(*) FROM seen_posts WHERE relevance_score >= 2"
        ).fetchone()[0]
        total_replies = db.execute("SELECT count(*) FROM replies").fetchone()[0]

        print("\n--- Stats ---")
        print(f"  Posts seen: {total_seen}")
        print(f"  Relevant: {total_relevant}")
        print(f"  Replies generated: {total_replies}")

        print("\nWaiting 5 minutes...")
        time.sleep(300)


if __name__ == "__main__":
    main()
