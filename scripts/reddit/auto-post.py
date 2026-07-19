#!/usr/bin/env python3
"""
TormentNexus Marketing Automation
Uses CDP (Chrome DevTools Protocol) to auto-post to Reddit, HN, Product Hunt, Twitter
Uses DeepSeek to generate platform-specific content
"""

import json
import time
import urllib.request
import os
from datetime import datetime

# DeepSeek API (set via environment variable)
DEEPSEEK_KEY = os.environ.get("DEEPSEEK_API_KEY", "")
DEEPSEEK_URL = "https://api.deepseek.com/v1/chat/completions"

# Platform configs
PLATFORMS = {
    "reddit": {
        "name": "Reddit",
        "subreddits": ["MachineLearning", "LocalLLaMA", "artificial", "opensource"],
        "url": "https://www.reddit.com/submit",
    },
    "hackernews": {
        "name": "Hacker News",
        "url": "https://news.ycombinator.com/submit",
    },
    "producthunt": {
        "name": "Product Hunt",
        "url": "https://www.producthunt.com/posts/new",
    },
    "twitter": {
        "name": "Twitter/X",
        "url": "https://twitter.com/compose/tweet",
    },
}


def call_deepseek(prompt, max_tokens=500):
    """Call DeepSeek API for content generation"""
    try:
        data = json.dumps(
            {
                "model": "deepseek-chat",
                "messages": [{"role": "user", "content": prompt}],
                "max_tokens": max_tokens,
                "temperature": 0.7,
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
    except Exception as e:
        print(f"  [!] DeepSeek error: {e}")
        return None


def generate_reddit_post(subreddit):
    """Generate a Reddit post for a specific subreddit"""
    prompt = f"""Write a Reddit post for r/{subreddit} about TormentNexus, an open-source AI control plane.

Key features to highlight:
- 26,000+ MCP (Model Context Protocol) servers catalog
- Persistent memory system (L1→L2→L3→L4 tiers)
- Local LLM support (LM Studio, Ollama, DeepSeek)
- Unified dashboard for monitoring
- Go backend, single binary, fast

The post should:
- Be authentic and not too salesy
- Include a compelling title
- Have a clear call-to-action
- Mention it's open source
- Include GitHub link: https://github.com/MDMAtk/TormentNexus
- Include live catalog: https://tormentnexus.site/catalog

Format as JSON with "title" and "body" fields."""

    result = call_deepseek(prompt)
    if result:
        try:
            # Try to parse JSON
            result = result.strip()
            if result.startswith("```"):
                result = result.split("```")[1]
                if result.startswith("json"):
                    result = result[4:]
            return json.loads(result)
        except:
            # Fallback: use raw text
            return {
                "title": f"Show r/{subreddit}: TormentNexus - AI Control Plane",
                "body": result,
            }
    return None


def generate_hn_post():
    """Generate a Hacker News post"""
    prompt = """Write a Show HN post for TormentNexus, an open-source AI control plane.

Key features:
- 26,000+ MCP servers catalog (searchable database)
- Persistent memory system for AI agents
- Local LLM support (no data leaves your machine)
- Go backend, single binary, fast startup
- Unified dashboard

The post should:
- Be technical and substantive
- Explain the problem it solves
- Include architecture details
- Mention it's open source
- GitHub: https://github.com/MDMAtk/TormentNexus
- Live: https://tormentnexus.site

Format as JSON with "title" and "body" fields."""

    result = call_deepseek(prompt)
    if result:
        try:
            result = result.strip()
            if result.startswith("```"):
                result = result.split("```")[1]
                if result.startswith("json"):
                    result = result[4:]
            return json.loads(result)
        except:
            return {
                "title": "Show HN: TormentNexus – Open-source AI control plane with 26K+ MCP tools",
                "body": result,
            }
    return None


def generate_twitter_thread():
    """Generate a Twitter thread"""
    prompt = """Write a Twitter thread (5-7 tweets) about TormentNexus, an open-source AI control plane.

Key points to cover:
1. The problem: AI assistants lose context between sessions
2. The solution: Persistent memory system (L1→L2→L3→L4)
3. 26,000+ MCP servers catalog - searchable database
4. Local LLM support - no data leaves your machine
5. Go backend, single binary, fast
6. Open source - anyone can contribute
7. Call to action: Try it now

Rules:
- Each tweet should be under 280 characters
- Use emojis sparingly but effectively
- Include GitHub link in last tweet
- Make it a compelling narrative
- Thread should build excitement

Format as JSON with "tweets" array."""

    result = call_deepseek(prompt)
    if result:
        try:
            result = result.strip()
            if result.startswith("```"):
                result = result.split("```")[1]
                if result.startswith("json"):
                    result = result[4:]
            return json.loads(result)
        except:
            return {
                "tweets": [
                    "🧵 Thread: Why AI needs persistent memory (1/7)",
                    "Current AI assistants lose context between sessions. You explain your project, close the laptop, and start over. Sound familiar?",
                    "TormentNexus fixes this with a tiered memory system: L1 (session) → L2 (hot vector) → L3 (cold archive) → L4 (limbo). Your AI remembers.",
                    "It also has 26,000+ MCP servers indexed and searchable. Databases, filesystems, browsers, APIs - one click to install.",
                    "Works with local LLMs (LM Studio, Ollama, DeepSeek). No data leaves your machine unless you want it to.",
                    "Built in Go. Single binary. Fast startup. Low resource usage.",
                    "Open source: https://github.com/MDMAtk/TormentNexus\n\nTry it now: https://tormentnexus.site",
                ]
            }
    return None


def generate_producthunt_post():
    """Generate a Product Hunt post"""
    prompt = """Write a Product Hunt launch post for TormentNexus, an open-source AI control plane.

Key features:
- 26,000+ MCP servers catalog
- Persistent memory for AI agents
- Local LLM support
- Unified dashboard
- Go backend, single binary

The post should have:
- Tagline (under 60 characters)
- Description (2-3 paragraphs)
- 3 key features
- Maker's comment

Format as JSON with "tagline", "description", "features", and "maker_comment" fields."""

    result = call_deepseek(prompt)
    if result:
        try:
            result = result.strip()
            if result.startswith("```"):
                result = result.split("```")[1]
                if result.startswith("json"):
                    result = result[4:]
            return json.loads(result)
        except:
            return {
                "tagline": "Give your AI persistent memory and 26K+ tools",
                "description": result,
                "features": [
                    "Persistent Memory System",
                    "26,000+ MCP Servers",
                    "Local LLM Support",
                ],
                "maker_comment": "Built this to solve my own problem - AI assistants losing context between sessions.",
            }
    return None


def save_post(platform, content, subreddit=None):
    """Save generated post to file"""
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")

    if subreddit:
        filename = f"marketing/auto/{platform}_{subreddit}_{timestamp}.json"
    else:
        filename = f"marketing/auto/{platform}_{timestamp}.json"

    os.makedirs(os.path.dirname(filename), exist_ok=True)

    with open(filename, "w") as f:
        json.dump(content, f, indent=2)

    print(f"  [OK] Saved to {filename}")
    return filename


def create_chrome_script(platform, content, subreddit=None):
    """Create a Chrome CDP script for posting"""
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    script_file = f"marketing/auto/{platform}_{timestamp}.js"

    if platform == "reddit":
        title = content.get("title", "")
        body = content.get("body", "")
        url = (
            f"https://www.reddit.com/r/{subreddit}/submit"
            if subreddit
            else "https://www.reddit.com/submit"
        )

        script = f"""
// Reddit posting script - {platform}/{subreddit}
// Generated: {datetime.now().isoformat()}

const puppeteer = require('puppeteer');

(async () => {{
  const browser = await puppeteer.launch({{ headless: false }});
  const page = await browser.newPage();
  
  // Navigate to Reddit submit page
  await page.goto('{url}');
  await page.waitForSelector('textarea[placeholder]', {{ timeout: 30000 }});
  
  // Fill in title
  const titleInput = await page.$('textarea[placeholder*="Title"]');
  if (titleInput) {{
    await titleInput.type({json.dumps(title)});
  }}
  
  // Fill in body
  const bodyInput = await page.$('textarea[placeholder*="body"]');
  if (bodyInput) {{
    await bodyInput.type({json.dumps(body)});
  }}
  
  // Wait for manual review
  console.log('Post ready for review. Press Enter to submit...');
  await new Promise(resolve => {{
    process.stdin.once('data', resolve);
  }});
  
  // Submit
  const submitButton = await page.$('button[type="submit"]');
  if (submitButton) {{
    await submitButton.click();
  }}
  
  await browser.close();
}})();
"""
    elif platform == "hackernews":
        title = content.get("title", "")
        body = content.get("body", "")

        script = f"""
// Hacker News posting script
// Generated: {datetime.now().isoformat()}

const puppeteer = require('puppeteer');

(async () => {{
  const browser = await puppeteer.launch({{ headless: false }});
  const page = await browser.newPage();
  
  // Navigate to HN submit page
  await page.goto('https://news.ycombinator.com/submit');
  await page.waitForSelector('input[name="title"]', {{ timeout: 30000 }});
  
  // Fill in title
  await page.type('input[name="title"]', {json.dumps(title)});
  
  // Fill in URL (if applicable)
  await page.type('input[name="url"]', 'https://github.com/MDMAtk/TormentNexus');
  
  // Fill in text
  const textInput = await page.$('textarea[name="text"]');
  if (textInput) {{
    await textInput.type({json.dumps(body)});
  }}
  
  // Wait for manual review
  console.log('Post ready for review. Press Enter to submit...');
  await new Promise(resolve => {{
    process.stdin.once('data', resolve);
  }});
  
  // Submit
  const submitButton = await page.$('input[type="submit"]');
  if (submitButton) {{
    await submitButton.click();
  }}
  
  await browser.close();
}})();
"""
    elif platform == "twitter":
        tweets = content.get("tweets", [])

        script = f"""
// Twitter thread posting script
// Generated: {datetime.now().isoformat()}

const puppeteer = require('puppeteer');

(async () => {{
  const browser = await puppeteer.launch({{ headless: false }});
  const page = await browser.newPage();
  
  // Navigate to Twitter
  await page.goto('https://twitter.com/compose/tweet');
  await page.waitForSelector('textarea[aria-label]', {{ timeout: 30000 }});
  
  const tweets = {json.dumps(tweets)};
  
  for (let i = 0; i < tweets.length; i++) {{
    console.log(`Posting tweet ${{i + 1}}/${{tweets.length}}...`);
    
    // Type tweet
    const tweetInput = await page.$('textarea[aria-label]');
    if (tweetInput) {{
      await tweetInput.type(tweets[i]);
    }}
    
    // Click tweet button
    const tweetButton = await page.$('button[data-testid="tweetButton"]');
    if (tweetButton) {{
      await tweetButton.click();
    }}
    
    // Wait between tweets
    if (i < tweets.length - 1) {{
      await new Promise(r => setTimeout(r, 3000));
      
      // Click reply to continue thread
      const replyButton = await page.$('button[data-testid="reply"]');
      if (replyButton) {{
        await replyButton.click();
      }}
    }}
  }}
  
  console.log('Thread posted!');
  await browser.close();
}})();
"""
    elif platform == "producthunt":
        tagline = content.get("tagline", "")
        description = content.get("description", "")
        features = content.get("features", [])

        script = f"""
// Product Hunt posting script
// Generated: {datetime.now().isoformat()}

const puppeteer = require('puppeteer');

(async () => {{
  const browser = await puppeteer.launch({{ headless: false }});
  const page = await browser.newPage();
  
  // Navigate to Product Hunt
  await page.goto('https://www.producthunt.com/posts/new');
  await page.waitForSelector('input[name="name"]', {{ timeout: 30000 }});
  
  // Fill in product name
  await page.type('input[name="name"]', 'TormentNexus');
  
  // Fill in tagline
  await page.type('input[name="tagline"]', {json.dumps(tagline)});
  
  // Fill in description
  const descInput = await page.$('textarea[name="description"]');
  if (descInput) {{
    await descInput.type({json.dumps(description)});
  }}
  
  // Wait for manual review
  console.log('Post ready for review. Press Enter to submit...');
  await new Promise(resolve => {{
    process.stdin.once('data', resolve);
  }});
  
  // Submit
  const submitButton = await page.$('button[type="submit"]');
  if (submitButton) {{
    await submitButton.click();
  }}
  
  await browser.close();
}})();
"""
    else:
        return None

    with open(script_file, "w") as f:
        f.write(script)

    print(f"  [OK] CDP script saved to {script_file}")
    return script_file


def main():
    print("=" * 50)
    print("  TormentNexus Marketing Automation")
    print("=" * 50)
    print()

    # Generate content for all platforms
    print("[1/4] Generating Reddit posts...")
    for subreddit in PLATFORMS["reddit"]["subreddits"]:
        print(f"  Generating for r/{subreddit}...")
        content = generate_reddit_post(subreddit)
        if content:
            save_post("reddit", content, subreddit)
            create_chrome_script("reddit", content, subreddit)
        time.sleep(1)  # Rate limit

    print()
    print("[2/4] Generating Hacker News post...")
    hn_content = generate_hn_post()
    if hn_content:
        save_post("hackernews", hn_content)
        create_chrome_script("hackernews", hn_content)

    print()
    print("[3/4] Generating Twitter thread...")
    twitter_content = generate_twitter_thread()
    if twitter_content:
        save_post("twitter", twitter_content)
        create_chrome_script("twitter", twitter_content)

    print()
    print("[4/4] Generating Product Hunt post...")
    ph_content = generate_producthunt_post()
    if ph_content:
        save_post("producthunt", ph_content)
        create_chrome_script("producthunt", ph_content)

    print()
    print("=" * 50)
    print("  Content generation complete!")
    print("=" * 50)
    print()
    print("Generated files are in: marketing/auto/")
    print()
    print("To post:")
    print("  1. Review the generated content in marketing/auto/*.json")
    print("  2. Run the CDP scripts: node marketing/auto/*.js")
    print("  3. Each script opens a browser for manual review before posting")
    print()
    print("Note: You need to be logged into each platform in the browser")


if __name__ == "__main__":
    main()
