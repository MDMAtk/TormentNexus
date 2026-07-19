const puppeteer = require("puppeteer");

(async () => {
  console.log("Launching browser...");
  const browser = await puppeteer.launch({
    headless: false,
    args: ["--start-maximized"]
  });
  
  const page = await browser.newPage();
  console.log("Browser launched!");
  
  await page.goto("https://news.ycombinator.com");
  const title = await page.title();
  console.log("Page title:", title);
  
  console.log("\nBrowser is open! Close it manually or press Ctrl+C.");
})();
