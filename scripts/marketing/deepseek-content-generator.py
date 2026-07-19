import os
import json
import urllib.request
import urllib.parse
from datetime import datetime

# Define our 5 blog post topics based on the ROADMAP
TOPICS = [
    {
        "slug": "why-ai-needs-memory-the-end-of-goldfish-agents",
        "title": "Why AI Needs Memory: The End of Goldfish Agents",
        "prompt": "Write a blog post explaining why AI agents need persistent memory. Explain the 'goldfish memory' problem where agents forget everything after a restart. Introduce TormentNexus as the solution, explaining its L1 (Scratchpad) and L2 (Vault) memory architecture."
    },
    {
        "slug": "mcp-the-future-of-ai-tool-integration",
        "title": "MCP: The Future of AI Tool Integration",
        "prompt": "Write a blog post about the Model Context Protocol (MCP). Explain why standardizing tool integrations is crucial for the AI ecosystem. Describe how TormentNexus acts as a universal MCP registry and proxy, managing over 26,000 tools."
    },
    {
        "slug": "orchestrating-the-swarm-multi-agent-collaboration",
        "title": "Orchestrating the Swarm: Multi-Agent Collaboration",
        "prompt": "Write a blog post about multi-agent workflows. Explain the 'Planner -> Implementer -> Tester -> Critic' collaboration cycle. Highlight how TormentNexus orchestrates these swarms using a Go-native event bus for resilient message brokering."
    },
    {
        "slug": "the-llm-waterfall-escaping-api-rate-limits",
        "title": "The LLM Waterfall: Escaping API Rate Limits",
        "prompt": "Write a blog post about the LLM Waterfall pattern. Explain how single-provider architectures fail due to rate limits and outages. Detail how TormentNexus cascades through multiple providers (e.g., NVIDIA NIM -> OpenRouter -> Local Ollama) to ensure zero downtime."
    },
    {
        "slug": "zero-trust-ai-securing-autonomous-agents",
        "title": "Zero-Trust AI: Securing Autonomous Agents",
        "prompt": "Write a blog post on Zero-Trust AI architecture. Discuss the threat model of autonomous agents executing arbitrary code. Explain how TormentNexus secures agents using RBAC policies, authenticated tool calls, and an immutable audit trail."
    }
]

def generate_deepseek_content(prompt):
    api_key = os.environ.get("DEEPSEEK_API_KEY")
    if not api_key:
        print("Warning: DEEPSEEK_API_KEY not found. Using mock content generation.")
        return f"<p>This is a mock generated blog post for: {prompt[:50]}...</p>\n<p>To generate real content, set the DEEPSEEK_API_KEY environment variable.</p>"

    url = "https://api.deepseek.com/chat/completions"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {api_key}"
    }

    data = {
        "model": "deepseek-chat",
        "messages": [
            {"role": "system", "content": "You are a technical content writer for TormentNexus, an open-source AI orchestration platform. Write engaging, technical blog posts formatted in HTML (only output the HTML body content, no markdown blocks, no enclosing html/body tags). Use <h2> for headings, <p> for paragraphs, and <ul>/<li> for lists. Keep it professional but edgy."},
            {"role": "user", "content": prompt}
        ],
        "temperature": 0.7
    }

    req = urllib.request.Request(url, data=json.dumps(data).encode('utf-8'), headers=headers, method='POST')

    try:
        with urllib.request.urlopen(req) as response:
            result = json.loads(response.read().decode('utf-8'))
            return result['choices'][0]['message']['content'].strip()
    except Exception as e:
        print(f"Error calling DeepSeek API: {e}")
        return f"<p>Error generating content: {e}</p>"

def main():
    # Ensure output directory exists
    out_dir = "data/marketing-content/blog"
    os.makedirs(out_dir, exist_ok=True)

    current_date = datetime.now().strftime("%B %d, %Y")

    print(f"Generating {len(TOPICS)} blog posts using DeepSeek API...")

    for topic in TOPICS:
        print(f"Generating: {topic['title']}")
        html_body = generate_deepseek_content(topic['prompt'])

        full_html = f'''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{topic["title"]} — TormentNexus Blog</title>
    <style>
        * {{ margin:0; padding:0; box-sizing:border-box; }}
        body {{ background:#0a0a0f; color:#e0e0e0; font-family:system-ui,sans-serif; line-height:1.8; }}
        .c {{ max-width:740px; margin:0 auto; padding:60px 20px; }}
        h1 {{ font-size:2.2rem; background:linear-gradient(135deg,#7c3aed,#3b82f6); -webkit-background-clip:text; -webkit-text-fill-color:transparent; margin-bottom: 0.5rem; }}
        .m {{ color:#6868a0; font-size:.85rem; margin-bottom:2rem; }}
        h2 {{ font-size:1.4rem; margin:2rem 0 .8rem; color:#c4b5fd; }}
        p {{ margin-bottom:1.2rem; color:#b8b8c8; }}
        ul {{ margin: 1rem 0 1.5rem 1.5rem; color: #b8b8c8; }}
        li {{ margin-bottom: 0.5rem; }}
        a {{ color:#a78bfa; text-decoration:none; }}
        a:hover {{ text-decoration:underline; }}
        .b {{ display:inline-block; margin-bottom:2rem; color:#6868a0; }}
    </style>
</head>
<body>
    <div class="c">
        <a href="/" class="b">← TormentNexus</a>
        <h1>{topic["title"]}</h1>
        <p class="m">{current_date} · TormentNexus Team</p>

        {html_body}

        <p style="margin-top:3rem; border-top:1px solid #1a1a2e; padding-top:2rem; color:#6868a0; font-size:.85rem">
            <a href="https://github.com/MDMAtk/TormentNexus">GitHub</a> · <a href="https://hypernexus.site/docs">Docs</a>
        </p>
    </div>
</body>
</html>'''

        out_path = os.path.join(out_dir, f"{topic['slug']}.html")
        with open(out_path, "w", encoding="utf-8") as f:
            f.write(full_html)

        print(f"  -> Saved to {out_path}")

    print("\nGeneration complete!")

if __name__ == "__main__":
    main()
