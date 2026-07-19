
// Twitter thread posting script
// Generated: 2026-07-13T23:13:26.409417

const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ headless: false });
  const page = await browser.newPage();
  
  // Navigate to Twitter
  await page.goto('https://twitter.com/compose/tweet');
  await page.waitForSelector('textarea[aria-label]', { timeout: 30000 });
  
  const tweets = ["\ud83e\uddf5 Thread: Why AI needs persistent memory (1/7)", "Current AI assistants lose context between sessions. You explain your project, close the laptop, and start over. Sound familiar?", "TormentNexus fixes this with a tiered memory system: L1 (session) \u2192 L2 (hot vector) \u2192 L3 (cold archive) \u2192 L4 (limbo). Your AI remembers.", "It also has 26,000+ MCP servers indexed and searchable. Databases, filesystems, browsers, APIs - one click to install.", "Works with local LLMs (LM Studio, Ollama, DeepSeek). No data leaves your machine unless you want it to.", "Built in Go. Single binary. Fast startup. Low resource usage.", "Open source: https://github.com/MDMAtk/TormentNexus\n\nTry it now: https://tormentnexus.site"];
  
  for (let i = 0; i < tweets.length; i++) {
    console.log(`Posting tweet ${i + 1}/${tweets.length}...`);
    
    // Type tweet
    const tweetInput = await page.$('textarea[aria-label]');
    if (tweetInput) {
      await tweetInput.type(tweets[i]);
    }
    
    // Click tweet button
    const tweetButton = await page.$('button[data-testid="tweetButton"]');
    if (tweetButton) {
      await tweetButton.click();
    }
    
    // Wait between tweets
    if (i < tweets.length - 1) {
      await new Promise(r => setTimeout(r, 3000));
      
      // Click reply to continue thread
      const replyButton = await page.$('button[data-testid="reply"]');
      if (replyButton) {
        await replyButton.click();
      }
    }
  }
  
  console.log('Thread posted!');
  await browser.close();
})();
