#!/usr/bin/env python3
"""
TormentNexus Reddit Thread Finder
Finds relevant threads where posting TN link would be helpful
"""

import os
import json
import urllib.request

# Subreddits to monitor
SUBREDDITS = [
    "MachineLearning",
    "LocalLLaMA",
    "artificial",
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
    "singularity",
    "Automate",
    "programming",
    "golang",
]

# High-relevance keywords (posting link is genuinely helpful)
HIGH_RELEVANCE_KEYWORDS = [
    # MCP related
    "mcp server",
    "model context protocol",
    "mcp tool",
    "mcp catalog",
    "mcp registry",
    "tool use",
    "function calling",
    # Memory related
    "persistent memory",
    "ai memory",
    "memory for ai",
    "give ai memory",
    "ai that remembers",
    "context window",
    "forget",
    "memory system",
    "memory tier",
    # AI control plane
    "ai control plane",
    "ai orchestration",
    "ai agent framework",
    "multi-agent",
    "agent orchestration",
    # Local LLM
    "local llm",
    "ollama",
    "lm studio",
    "localai",
    "self-hosted ai",
    # Tools
    "ai tools",
    "ai dev tools",
    "developer tools",
    "tool catalog",
    "tool registry",
]


def fetch_subreddit_posts(subreddit, limit=50):
    """Fetch posts using Reddit's public JSON API"""
    url = f"https://www.reddit.com/r/{subreddit}/hot.json?limit={limit}"
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
    """Calculate how relevant a post is to TormentNexus"""
    text = (title + " " + selftext).lower()

    score = 0
    matched_keywords = []

    for keyword in HIGH_RELEVANCE_KEYWORDS:
        if keyword.lower() in text:
            score += 1
            matched_keywords.append(keyword)

    return score, matched_keywords


def main():
    """Find relevant threads"""
    print("=" * 60)
    print("  TormentNexus Reddit Thread Finder")
    print("=" * 60)
    print()
    print("Finding threads where posting TN link would be helpful...")
    print()

    all_relevant = []

    for subreddit in SUBREDDITS:
        print(f"--- r/{subreddit} ---")

        posts = fetch_subreddit_posts(subreddit, limit=50)

        for post in posts:
            data = post.get("data", {})
            title = data.get("title", "")
            selftext = data.get("selftext", "")
            permalink = data.get("permalink", "")
            url = data.get("url", "")

            score, keywords = calculate_relevance(title, selftext)

            if score >= 2:  # At least 2 keyword matches
                all_relevant.append(
                    {
                        "subreddit": subreddit,
                        "title": title,
                        "url": f"https://reddit.com{permalink}",
                        "score": score,
                        "keywords": keywords,
                    }
                )

    # Sort by relevance score
    all_relevant.sort(key=lambda x: x["score"], reverse=True)

    print("\n" + "=" * 60)
    print(f"  Found {len(all_relevant)} relevant threads")
    print("=" * 60)
    print()

    # Show top results
    for i, thread in enumerate(all_relevant[:30], 1):
        print(f"{i}. [Score: {thread['score']}] r/{thread['subreddit']}")
        print(f"   Title: {thread['title'][:80]}...")
        print(f"   Keywords: {', '.join(thread['keywords'][:5])}")
        print(f"   URL: {thread['url']}")
        print()

    # Save to file
    output_file = "/opt/tormentnexus/data/relevant-threads.json"
    os.makedirs(os.path.dirname(output_file), exist_ok=True)
    with open(output_file, "w") as f:
        json.dump(all_relevant, f, indent=2)

    print(f"Saved {len(all_relevant)} threads to {output_file}")


if __name__ == "__main__":
    main()
