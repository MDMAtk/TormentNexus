
// Reddit posting script - reddit/MachineLearning
// Generated: 2026-07-13T23:12:48.327201

const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ headless: false });
  const page = await browser.newPage();
  
  // Navigate to Reddit submit page
  await page.goto('https://www.reddit.com/r/MachineLearning/submit');
  await page.waitForSelector('textarea[placeholder]', { timeout: 30000 });
  
  // Fill in title
  const titleInput = await page.$('textarea[placeholder*="Title"]');
  if (titleInput) {
    await titleInput.type("I built an open-source AI control plane with 26K+ MCP servers and persistent memory tiers (TormentNexus)");
  }
  
  // Fill in body
  const bodyInput = await page.$('textarea[placeholder*="body"]');
  if (bodyInput) {
    await bodyInput.type("Hey r/MachineLearning,\n\nI've been working on a side project for a while now called **TormentNexus** \u2014 an open-source AI control plane that's grown way bigger than I initially expected. Thought I'd share it here since it might be useful for some of you.\n\n**What it does:**\n- **26,000+ MCP servers** already cataloged and queryable (live at https://tormentnexus.site/catalog)\n- **4-tier persistent memory** (L1\u2192L2\u2192L3\u2192L4) so your AI agents can actually remember things across sessions without context blowing up\n- **Local LLM support** \u2014 works with LM Studio, Ollama, DeepSeek, and others out of the box\n- **Unified dashboard** for monitoring all your agents, memory usage, and server health\n- **Go backend, single binary** \u2014 it's fast, lightweight, and stupid easy to deploy\n\n**Why I built it:**\nI got tired of juggling a dozen different tools and having agents forget everything the moment I closed a chat. Wanted something that felt like a real control plane \u2014 not just another wrapper. It's still early days but I use it daily for my own projects.\n\n**It's open source:** https://github.com/MDMAtk/TormentNexus\n\nWould love feedback, issues, or PRs. Specifically curious if anyone's tried it with larger agentic workflows or custom MCP servers beyond what's in the catalog.\n\n**TL;DR:** Open-source AI control plane with 26K+ MCP servers, persistent memory tiers, local LLM support, and a Go backend. GitHub: https://github.com/MDMAtk/TormentNexus \u2014 live catalog: https://tormentnexus.site/catalog");
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
