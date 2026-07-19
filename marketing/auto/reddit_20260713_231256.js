
// Reddit posting script - reddit/LocalLLaMA
// Generated: 2026-07-13T23:12:56.898843

const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ headless: false });
  const page = await browser.newPage();
  
  // Navigate to Reddit submit page
  await page.goto('https://www.reddit.com/r/LocalLLaMA/submit');
  await page.waitForSelector('textarea[placeholder]', { timeout: 30000 });
  
  // Fill in title
  const titleInput = await page.$('textarea[placeholder*="Title"]');
  if (titleInput) {
    await titleInput.type("Show r/LocalLLaMA: TormentNexus - AI Control Plane");
  }
  
  // Fill in body
  const bodyInput = await page.$('textarea[placeholder*="body"]');
  if (bodyInput) {
    await bodyInput.type("\n{\n  \"title\": \"Built a stupidly fast AI control plane with 26k+ MCP servers and persistent memory tiers - TormentNexus is now open source\",\n  \"body\": \"Hey r/LocalLLaMA,\\n\\nI've been working on this for a while and finally feel like it's ready to share. It's called TormentNexus - an open-source AI control plane that started as a personal itch and grew into something I actually use daily.\\n\\n**What it does:**\\n- **26,000+ MCP servers catalog** - yeah, you read that right. It's a live, searchable index of Model Context Protocol servers. Found it super useful for discovering tools without digging through scattered repos.\\n- **Persistent memory system** - L1\u2192L2\u2192L3\u2192L4 tiered memory. L1 is short-term session context, L2 is working memory, L3 is long-term storage, L4 is archival. Your LLM actually remembers stuff across sessions.\\n- **Local LLM support** - Works with LM Studio, Ollama, DeepSeek, and any OpenAI-compatible endpoint. No cloud required if you don't want it.\\n- **Unified dashboard** - See all your agents, memory usage, active MCP connections in one place. Real-time monitoring without the bloat.\\n- **Go backend, single binary** - It's fast. Like, really fast. One binary, no dependencies, runs on anything.\\n\\n**Why I built it:**\\nI got tired of juggling 5 different terminals, forgetting which MCP server did what, and having my LLMs forget everything between conversations. The tiered memory thing was a game-changer for me - L1/L2 are ephemeral, L3 persists across sessions, L4 is basically your AI's \\\"life story.\\\"\\n\\n**It's open source:**\\nGitHub: https://github.com/MDMAtk/TormentNexus\\nLive catalog (no install needed to browse): https://tormentnexus.site/catalog\\n\\n**What I'd love:**\\n- Try it out, break it, tell me what sucks\\n- PRs welcome - especially for more MCP server integrations\\n- Ideas for the memory system (I'm thinking vector search for L3/L4)\\n- If you're into");
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
