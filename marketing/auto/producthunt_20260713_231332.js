
// Product Hunt posting script
// Generated: 2026-07-13T23:13:32.587934

const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ headless: false });
  const page = await browser.newPage();
  
  // Navigate to Product Hunt
  await page.goto('https://www.producthunt.com/posts/new');
  await page.waitForSelector('input[name="name"]', { timeout: 30000 });
  
  // Fill in product name
  await page.type('input[name="name"]', 'TormentNexus');
  
  // Fill in tagline
  await page.type('input[name="tagline"]', "The open-source control plane for all your AI agents.");
  
  // Fill in description
  const descInput = await page.$('textarea[name="description"]');
  if (descInput) {
    await descInput.type("Meet TormentNexus, the first open-source AI control plane designed to give developers and power users complete command over their AI agents. With a catalog of over 26,000 MCP servers, you can instantly connect your agents to the tools and data they need\u2014from databases to APIs to custom services. No more cobbling together brittle integrations or losing context when your agents restart.\n\nBuilt on a high-performance Go backend and delivered as a single binary, TormentNexus is lightweight, fast, and deploys anywhere. It features persistent memory for your AI agents, so they remember conversations, preferences, and workflows across sessions. Whether you're running local LLMs for privacy or cloud models for scale, the unified dashboard gives you real-time visibility and control over every agent's behavior, memory, and tool usage.");
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
