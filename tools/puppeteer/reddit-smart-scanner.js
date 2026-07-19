#!/usr/bin/env node
/**
 * TormentNexus Smart Reddit Scanner
 * Uses AI to find threads where TormentNexus genuinely fits
 * Generates contextual, organic comments
 */

const puppeteer = require("puppeteer");
const fs = require("fs");
const path = require("path");
const https = require("https");

// ── Configuration ──
const SUBREDDITS = [
  "MachineLearning",
  "LocalLLaMA",
  "artificial",
  "singularity",
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
];

// MiMo API config
const MIMO_KEY = process.env.MIMO_API_KEY || "";
const MIMO_URL = "https://token-plan-sgp.xiaomimimo.com/v1/chat/completions";
const MIMO_MODEL = "mimo-v2.5";

// Paths
const QUEUE_FILE = path.join(__dirname, "..", "data", "reddit-queue.json");

// ── Helper Functions ──
function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function waitForEnter(message) {
  return new Promise((resolve) => {
    const readline = require("readline");
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
    });
    rl.question(message + " ", (answer) => {
      rl.close();
      resolve(answer);
    });
  });
}

async function callMiMo(prompt, maxTokens = 1000) {
  if (!MIMO_KEY) {
    console.log("  [!] No MIMO_API_KEY set");
    return null;
  }

  return new Promise((resolve, reject) => {
    const data = JSON.stringify({
      model: MIMO_MODEL,
      messages: [{ role: "user", content: prompt }],
      max_tokens: maxTokens,
      temperature: 0.7,
    });

    const url = new URL(MIMO_URL);
    const options = {
      hostname: url.hostname,
      port: 443,
      path: url.pathname,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${MIMO_KEY}`,
        "Content-Length": Buffer.byteLength(data),
      },
    };

    const req = https.request(options, (res) => {
      let responseData = "";
      res.on("data", (chunk) => (responseData += chunk));
      res.on("end", () => {
        try {
          const result = JSON.parse(responseData);
          if (result.choices && result.choices[0]) {
            resolve(result.choices[0].message.content.trim());
          } else {
            resolve(null);
          }
        } catch (e) {
          reject(e);
        }
      });
    });

    req.on("error", reject);
    req.write(data);
    req.end();
  });
}

// ── AI Analysis ──
async function analyzeThread(title, selftext, subreddit) {
  const prompt = `You are analyzing a Reddit post to determine if mentioning TormentNexus would be genuinely helpful and organic.

TormentNexus is an open-source AI control plane that provides:
- 4-tier persistent memory system (L1-L4) for AI
- 26,000+ MCP (Model Context Protocol) tools catalog
- Multi-agent orchestration
- Works with Ollama, LM Studio, DeepSeek, any OpenAI-compatible API
- Go backend with Next.js dashboard
- GitHub: https://github.com/MDMAtk/TormentNexus
- Website: https://tormentnexus.site

Reddit Post:
Subreddit: r/${subreddit}
Title: ${title}
Content: ${selftext?.substring(0, 500) || "(no body)"}

Analyze this post and respond with ONLY a JSON object (no markdown):
{
  "isRelevant": true/false,
  "relevanceScore": 0-10,
  "reasonWhy": "brief explanation",
  "talkingPoints": ["specific points from the post that TormentNexus addresses"],
  "suggestedComment": "a natural, helpful comment (2-4 sentences) that mentions TormentNexus where it genuinely fits. Be conversational, not salesy. Only include the GitHub link if it adds value."
}

Rules:
- Only mark as relevant (score >= 6) if TormentNexus genuinely solves a problem mentioned in the post
- The comment should read like a helpful human response, not an ad
- Reference specific details from their post
- Only include links if they're genuinely useful
- If the post is about something TormentNexus doesn't help with, mark as not relevant`;

  try {
    const response = await callMiMo(prompt, 800);
    if (!response) return null;

    // Try to parse JSON from response
    const jsonMatch = response.match(/\{[\s\S]*\}/);
    if (jsonMatch) {
      return JSON.parse(jsonMatch[0]);
    }
    return null;
  } catch (e) {
    console.log(`  [!] AI analysis error: ${e.message}`);
    return null;
  }
}

// ── Main Function ──
async function main() {
  console.log("");
  console.log(
    "╔══════════════════════════════════════════════════════════╗"
  );
  console.log(
    "║  TormentNexus Smart Reddit Scanner                      ║"
  );
  console.log(
    "║  AI-powered thread analysis for organic engagement      ║"
  );
  console.log(
    "╚══════════════════════════════════════════════════════════╝"
  );
  console.log("");

  if (!MIMO_KEY) {
    console.log("  ⚠️  No MIMO_API_KEY set. AI analysis will be skipped.");
    console.log("  Set it with: export MIMO_API_KEY=your_key");
    console.log("");
  }

  console.log("This tool will:");
  console.log("1. Scan AI/ML subreddits for recent posts");
  console.log("2. Use AI to analyze if TormentNexus genuinely fits");
  console.log("3. Generate organic, contextual comments");
  console.log("4. Save to queue for manual review and posting");
  console.log("");
  console.log(
    "  Press ENTER to start scanning, or 'q' to quit:"
  );
  const startAnswer = await waitForEnter("");
  if (startAnswer && startAnswer.toLowerCase() === "q") {
    return;
  }

  const browser = await puppeteer.launch({
    headless: false,
    args: ["--start-maximized", "--no-sandbox"],
    defaultViewport: null,
  });

  const page = await browser.newPage();

  // Load existing queue
  let queue = [];
  if (fs.existsSync(QUEUE_FILE)) {
    try {
      queue = JSON.parse(fs.readFileSync(QUEUE_FILE, "utf8"));
    } catch (e) {
      queue = [];
    }
  }

  console.log("\nMonitoring subreddits:");
  for (const sub of SUBREDDITS) {
    console.log(`  - r/${sub}`);
  }
  console.log("");

  let totalFound = 0;
  let totalRelevant = 0;

  for (const subreddit of SUBREDDITS) {
    console.log(`\n${"─".repeat(60)}`);
    console.log(`  Scanning r/${subreddit}...`);
    console.log(`${"─".repeat(60)}`);

    try {
      // Navigate to subreddit's new posts
      await page.goto(`https://www.reddit.com/r/${subreddit}/new/`, {
        waitUntil: "networkidle2",
        timeout: 30000,
      });

      await sleep(3000);

      // Get posts with titles and links
      const posts = await page.evaluate(() => {
        const results = [];
        const seen = new Set();

        // Find all post links
        const links = document.querySelectorAll(
          'a[href*="/r/"][href*="/comments/"]'
        );
        for (const link of links) {
          const title = link.textContent?.trim();
          const href = link.href;
          if (
            title &&
            href &&
            !seen.has(href) &&
            title.length > 10
          ) {
            seen.add(href);
            results.push({
              title: title.substring(0, 300),
              href,
            });
          }
        }

        return results.slice(0, 10);
      });

      console.log(`  Found ${posts.length} posts`);

      for (const post of posts) {
        // Skip if already in queue
        if (queue.find((q) => q.url === post.href)) {
          continue;
        }

        totalFound++;
        console.log(`\n  Analyzing: ${post.title.substring(0, 60)}...`);

        // Navigate to the post to get content
        try {
          await page.goto(post.href, {
            waitUntil: "networkidle2",
            timeout: 20000,
          });
          await sleep(2000);

          // Get the post content
          const postContent = await page.evaluate(() => {
            // Try to get the post body
            const selectors = [
              '[data-testid="post-content"]',
              ".Post",
              'div[data-click-id="text"]',
              'div[slot="text-body"]',
            ];

            for (const selector of selectors) {
              const el = document.querySelector(selector);
              if (el) return el.textContent?.trim().substring(0, 500);
            }

            // Fallback: get any text content
            const body = document.querySelector("body");
            return body?.textContent?.trim().substring(0, 500);
          });

          // Use AI to analyze
          const analysis = await analyzeThread(
            post.title,
            postContent,
            subreddit
          );

          if (analysis && analysis.isRelevant && analysis.relevanceScore >= 6) {
            totalRelevant++;
            console.log(`  ✅ RELEVANT (score: ${analysis.relevanceScore}/10)`);
            console.log(`  Reason: ${analysis.reasonWhy}`);
            console.log(`  Talking points: ${analysis.talkingPoints?.join(", ")}`);
            console.log(`  Comment: ${analysis.suggestedComment?.substring(0, 80)}...`);

            // Add to queue
            queue.push({
              title: post.title,
              url: post.href,
              subreddit: subreddit,
              reply: analysis.suggestedComment,
              relevanceScore: analysis.relevanceScore,
              reasonWhy: analysis.reasonWhy,
              talkingPoints: analysis.talkingPoints,
              timestamp: new Date().toISOString(),
              posted: false,
            });

            // Save queue
            fs.mkdirSync(path.dirname(QUEUE_FILE), { recursive: true });
            fs.writeFileSync(QUEUE_FILE, JSON.stringify(queue, null, 2));
          } else {
            console.log(
              `  ❌ Not relevant (score: ${analysis?.relevanceScore || 0}/10)`
            );
          }
        } catch (e) {
          console.log(`  [!] Error analyzing post: ${e.message}`);
        }

        // Rate limiting
        await sleep(2000);
      }
    } catch (e) {
      console.log(`  [!] Error scanning subreddit: ${e.message}`);
    }

    // Rate limiting between subreddits
    await sleep(3000);
  }

  // Final summary
  console.log(`\n${"═".repeat(60)}`);
  console.log("  SCAN COMPLETE");
  console.log(`${"═".repeat(60)}`);
  console.log(`  Posts scanned: ${totalFound}`);
  console.log(`  Relevant threads found: ${totalRelevant}`);
  console.log(`  Queue file: ${QUEUE_FILE}`);
  console.log("");

  if (totalRelevant > 0) {
    console.log("  📋 Next steps:");
    console.log("  1. Review the queue file");
    console.log("  2. Manually post each comment on Reddit");
    console.log("  3. Mark as posted in the queue");
    console.log("");
    console.log("  View queue:");
    console.log(`  cat ${QUEUE_FILE}`);
  } else {
    console.log(
      "  No highly relevant threads found this scan. Try again later."
    );
  }

  await browser.close();
}

main().catch(console.error);
