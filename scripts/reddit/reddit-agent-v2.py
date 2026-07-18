#!/usr/bin/env python3
"""
TormentNexus Reddit Agent v2
Monitors AI/ML subreddits, finds relevant threads, generates and posts helpful replies
Uses PRAW for Reddit API and MiMo v2.5 for reply generation
"""

import os
import json
import time
import sqlite3
import sys
from datetime import datetime

try:
    import praw
except ImportError:
    print("Installing PRAW...")
    os.system(f"{sys.executable} -m pip install praw --quiet")
    import praw

# Try to import requests for MiMo API
try:
    import requests
except ImportError:
    os.system(f"{sys.executable} -m pip install requests --quiet")
    import requests

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

# Database
DB_PATH = os.path.join(os.path.dirname(__file__), "..", "data", "reddit-agent.db")

# Reply templates for common scenarios
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
    db.execute("""
        CREATE TABLE IF NOT EXISTS post_log (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id TEXT,
            subreddit TEXT,
            action TEXT,
            success BOOLEAN,
            error TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    """)
    db.commit()
    return db


def get_reddit_instance():
    """Get authenticated Reddit instance using PRAW"""
    creds_file = os.path.join(
        os.path.dirname(__file__), "..", "data", "reddit-creds.json"
    )

    if not os.path.exists(creds_file):
        print(f"  [!] No credentials found at {creds_file}")
        print("  Create the file with:")
        print(
            '  {"client_id": "...", "client_secret": "...", "username": "...", "password": "..."}'
        )
        return None

    with open(creds_file) as f:
        creds = json.load(f)

    if creds.get("client_id") == "YOUR_CLIENT_ID":
        print("  [!] Credentials are placeholder. Please update reddit-creds.json")
        return None

    try:
        reddit = praw.Reddit(
            client_id=creds["client_id"],
            client_secret=creds["client_secret"],
            username=creds["username"],
            password=creds["password"],
            user_agent="TormentNexus/1.0 (AI research assistant; https://github.com/MDMAtk/TormentNexus)",
        )
        # Test authentication
        print(f"  [+] Logged in as: {reddit.user.me()}")
        return reddit
    except Exception as e:
        print(f"  [!] Authentication failed: {e}")
        return None


def call_mimo(prompt, max_tokens=500):
    """Call MiMo API for reply generation"""
    if not MIMO_KEY:
        print("  [!] No MIMO_API_KEY set")
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
    # Check for template matches first
    text = (title + " " + selftext).lower()

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


def post_reply(reddit, post_id, reply_text, db):
    """Post a reply to Reddit"""
    try:
        submission = reddit.submission(id=post_id)
        comment = submission.reply(reply_text)

        # Log success
        db.execute(
            """
            INSERT INTO post_log (post_id, action, success)
            VALUES (?, 'reply', TRUE)
        """,
            (post_id,),
        )
        db.execute(
            """
            UPDATE replies SET posted = TRUE WHERE post_id = ?
        """,
            (post_id,),
        )
        db.commit()

        print(f"  [+] Posted reply: {comment.permalink}")
        return True
    except Exception as e:
        print(f"  [!] Failed to post: {e}")
        db.execute(
            """
            INSERT INTO post_log (post_id, action, success, error)
            VALUES (?, 'reply', FALSE, ?)
        """,
            (post_id, str(e)),
        )
        db.commit()
        return False


def main():
    """Main monitoring loop"""
    print("=" * 60)
    print("  TormentNexus Reddit Agent v2")
    print("=" * 60)
    print()

    db = init_db()
    reddit = get_reddit_instance()

    if not reddit:
        print("  [!] Running in SCAN-ONLY mode (no posting)")
        print("  Add Reddit credentials to enable posting")
    else:
        print("  [+] Running in AUTHENTICATED mode (will post)")

    print()
    print(f"  Monitoring {len(SUBREDDITS)} subreddits:")
    for sub in SUBREDDITS:
        print(f"    - r/{sub}")
    print()

    cycle = 0
    while True:
        cycle += 1
        print(f"\n{'=' * 60}")
        print(f"  Cycle {cycle} — {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"{'=' * 60}")

        for subreddit_name in SUBREDDITS:
            print(f"\n--- r/{subreddit_name} ---")

            try:
                subreddit = reddit.subreddit(subreddit_name) if reddit else None
                if subreddit:
                    posts = list(subreddit.new(limit=25))
                else:
                    # Fallback to JSON API for scanning
                    posts = []
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
                            for p in data.get("data", {}).get("children", []):
                                d = p.get("data", {})

                                # Create a simple object to mimic PRAW
                                class Post:
                                    def __init__(self, data):
                                        self.id = data.get("id")
                                        self.title = data.get("title", "")
                                        self.selftext = data.get("selftext", "")
                                        self.permalink = data.get("permalink", "")
                                        self.url = data.get("url", "")

                                posts.append(Post(d))
                    except Exception as e:
                        print(f"  [!] Error: {e}")
                        continue

                for post in posts:
                    post_id = post.id

                    # Skip if already seen
                    existing = db.execute(
                        "SELECT post_id FROM seen_posts WHERE post_id = ?", (post_id,)
                    ).fetchone()
                    if existing:
                        continue

                    # Calculate relevance
                    score, keywords, high_matches = calculate_relevance(
                        post.title, post.selftext
                    )

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
                            post.title,
                            post.url,
                            post.permalink,
                            score,
                        ),
                    )
                    db.commit()

                    if score >= 2:
                        include_link = len(high_matches) > 0
                        print(f"  [+] Relevant ({score}): {post.title[:60]}...")
                        print(f"      Keywords: {', '.join(keywords[:3])}")

                        if include_link:
                            print("      HIGH RELEVANCE — will include link")

                        # Generate reply
                        reply = generate_reply(
                            post.title, post.selftext, subreddit_name, include_link
                        )

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
                                """
                                UPDATE seen_posts SET replied = TRUE WHERE post_id = ?
                            """,
                                (post_id,),
                            )
                            db.commit()

                            # Post if authenticated
                            if reddit:
                                success = post_reply(reddit, post_id, reply, db)
                                if success:
                                    print("      [+] Posted successfully!")
                                time.sleep(10)  # Rate limiting
                            else:
                                # Save to pending file
                                pending_file = os.path.join(
                                    os.path.dirname(__file__),
                                    "..",
                                    "data",
                                    "pending-replies.json",
                                )
                                pending = []
                                if os.path.exists(pending_file):
                                    with open(pending_file) as f:
                                        pending = json.load(f)

                                pending.append(
                                    {
                                        "post_id": post_id,
                                        "subreddit": subreddit_name,
                                        "title": post.title,
                                        "url": f"https://reddit.com{post.permalink}",
                                        "reply": reply,
                                        "include_link": include_link,
                                        "score": score,
                                        "keywords": keywords,
                                        "timestamp": datetime.now().isoformat(),
                                    }
                                )

                                with open(pending_file, "w") as f:
                                    json.dump(pending, f, indent=2)

                time.sleep(2)  # Rate limiting between subreddits

            except Exception as e:
                print(f"  [!] Error with r/{subreddit_name}: {e}")
                continue

        # Stats
        total_seen = db.execute("SELECT count(*) FROM seen_posts").fetchone()[0]
        total_relevant = db.execute(
            "SELECT count(*) FROM seen_posts WHERE relevance_score >= 2"
        ).fetchone()[0]
        total_replies = db.execute("SELECT count(*) FROM replies").fetchone()[0]
        total_posted = db.execute(
            "SELECT count(*) FROM replies WHERE posted = TRUE"
        ).fetchone()[0]

        print("\n--- Stats ---")
        print(f"  Posts seen: {total_seen}")
        print(f"  Relevant: {total_relevant}")
        print(f"  Replies generated: {total_replies}")
        print(f"  Replies posted: {total_posted}")

        # Wait before next cycle (5 minutes)
        print("\n  Waiting 5 minutes before next cycle...")
        time.sleep(300)


if __name__ == "__main__":
    main()
