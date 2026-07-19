
// Reddit posting script - reddit/artificial
// Generated: 2026-07-13T23:13:04.467563

const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ headless: false });
  const page = await browser.newPage();
  
  // Navigate to Reddit submit page
  await page.goto('https://www.reddit.com/r/artificial/submit');
  await page.waitForSelector('textarea[placeholder]', { timeout: 30000 });
  
  // Fill in title
  const titleInput = await page.$('textarea[placeholder*="Title"]');
  if (titleInput) {
    await titleInput.type("I built an open-source AI control plane with 26,000+ MCP servers and persistent memory tiers (L1\u2192L4)");
  }
  
  // Fill in body
  const bodyInput = await page.$('textarea[placeholder*="body"]');
  if (bodyInput) {
    await bodyInput.type("Hey r/artificial,\n\nI\u2019ve been working on a side project for a while called **TormentNexus** \u2014 an open-source AI control plane that I originally built to scratch my own itch for managing multiple agents and models without vendor lock-in.\n\n**What it does:**\n- **26,000+ MCP (Model Context Protocol) servers** are already cataloged and searchable \u2014 you can plug in any compatible model/agent and route tasks through them.\n- **Persistent memory** organized in tiers (L1\u2192L2\u2192L3\u2192L4) so agents can retain context across sessions without blowing up token limits. Think short-term working memory up to long-term storage.\n- **Local LLM support** out of the box \u2014 works with LM Studio, Ollama, and DeepSeek. No cloud required if you don\u2019t want it.\n- **Unified dashboard** for monitoring all your agents, memory usage, and server status in real time.\n- **Go backend** \u2014 compiles to a single binary, stupid fast, easy to deploy on a VPS or a Raspberry Pi.\n\n**Why I built it:**\nI got tired of juggling half a dozen tools just to run local agents with persistent memory. Everything felt either too locked into a specific API, or too complicated to set up. TormentNexus is my attempt at a lightweight, open alternative that doesn\u2019t try to sell you a subscription.\n\n**It\u2019s fully open source** (MIT license) \u2014 check it out here:\n\ud83d\udd17 GitHub: https://github.com/MDMAtk/TormentNexus\n\n**Live catalog** of all MCP servers:\n\ud83d\udd17 https://tormentnexus.site/catalog\n\n**Call to action:**\nIf you\u2019re into self-hosted AI, agent orchestration, or just want to play with persistent memory without the cloud tax \u2014 give it a star, try the binary, and let me know what you think. PRs and issues welcome.\n\nHappy to answer questions in the comments. Cheers!");
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
