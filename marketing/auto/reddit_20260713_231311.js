
// Reddit posting script - reddit/opensource
// Generated: 2026-07-13T23:13:11.854914

const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ headless: false });
  const page = await browser.newPage();
  
  // Navigate to Reddit submit page
  await page.goto('https://www.reddit.com/r/opensource/submit');
  await page.waitForSelector('textarea[placeholder]', { timeout: 30000 });
  
  // Fill in title
  const titleInput = await page.$('textarea[placeholder*="Title"]');
  if (titleInput) {
    await titleInput.type("Show r/opensource: TormentNexus - AI Control Plane");
  }
  
  // Fill in body
  const bodyInput = await page.$('textarea[placeholder*="body"]');
  if (bodyInput) {
    await bodyInput.type("\n{\n  \"title\": \"I built an open-source AI control plane with 26K+ MCP servers and persistent memory tiers (TormentNexus)\",\n  \"body\": \"Hey r/opensource,\\n\\nI\u2019ve been working on a side project for the past few months that I think some of you might find interesting \u2014 especially if you\u2019re tired of juggling multiple AI tools, APIs, and memory systems manually.\\n\\n**What is TormentNexus?**\\n\\nIt\u2019s an open-source AI control plane that acts as a single dashboard for managing AI agents, tools, and memory. The core idea is simple: give you one place to orchestrate everything, without locking you into a proprietary service.\\n\\n**Key features that might matter to you:**\\n\\n- **26,000+ MCP servers catalog** \u2014 We\u2019ve indexed a massive library of Model Context Protocol servers, so you can plug in tools for web search, code execution, file access, and more. Live catalog here: https://tormentnexus.site/catalog\\n\\n- **Persistent memory system (L1\u2192L2\u2192L3\u2192L4 tiers)** \u2014 Instead of one-size-fits-all memory, you get tiered storage: short-term context (L1), working memory (L2), long-term (L3), and archival (L4). Helps agents remember things across sessions without blowing context windows.\\n\\n- **Local LLM support** \u2014 Works with LM Studio, Ollama, and DeepSeek out of the box. No cloud dependency if you don\u2019t want it.\\n\\n- **Unified dashboard** \u2014 See all your agents, memory usage, tool calls, and logs in one place. Built with a Go backend, compiles to a single binary, and is pretty fast.\\n\\n- **Open source** \u2014 MIT license, contributions welcome.\\n\\n**Why I built it**\\n\\nI got tired of switching between different agent frameworks, each with their own memory systems and tool sets. I wanted something that felt like a control room \u2014 one place to see what\u2019s happening, add new capabilities, and experiment without rewriting code every time.\\n\\n**Call to action**\\n\\nIf this sounds useful, give it a star, fork it, or open an issue. I\u2019m actively working on docs and more integrations. \\n");
  }
  
  // Wait for manual review
  console.log('Post ready for review. Press Enter to submit...');
  await new Promise(resolve => {
    process.stdin.once('data', resolve);
  });
  
  // Submit
  const submitButton = await page.$('button[type="submit"]');
  if (submitButton) {
    await submitButton.click();
  }
  
  await browser.close();
})();
