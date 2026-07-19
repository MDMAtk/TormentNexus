#!/usr/bin/env node
/**
 * TormentNexus Reddit Commenter
 * Finds relevant posts and adds helpful comments
 * Uses Chrome CDP for browser automation
 */

const puppeteer = require("puppeteer");
const fs = require("fs");
const path = require("path");

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

const KEYWORDS = [
	"memory",
	"persistent memory",
	"context window",
	"forget",
	"mcp",
	"model context protocol",
	"tool use",
	"function calling",
	"local llm",
	"ollama",
	"lm studio",
	"localai",
	"self-hosted",
	"ai agent",
	"multi-agent",
	"agent framework",
	"ai control plane",
	"ai orchestration",
	"vector database",
	"embedding",
	"rag",
	"ai assistant",
	"ai coding",
	"ai developer",
	"open source ai",
	"ai tools",
	// Lower threshold keywords
	"local",
	"self hosted",
	"open source",
	"free",
	"best",
	"recommend",
	"alternative",
	"vs",
	"help",
	"looking for",
	"suggestion",
	"how to",
	"tool",
	"server",
	"database",
	"agent",
	"framework",
	"llm",
	"model",
];

const HIGH_KEYWORDS = [
	"persistent memory",
	"mcp server",
	"model context protocol",
	"ai control plane",
	"ai that remembers",
	"memory for ai",
	"give ai memory",
	"ai memory system",
	"tool catalog",
	"mcp catalog",
	"mcp tools",
];

// Reply templates
const TEMPLATES = {
	memory: `I built TormentNexus which has a 4-tier memory system (L1-L4) for exactly this:

- L1: Session memory (ephemeral)
- L2: Working memory (30 days, semantically indexed)
- L3: Cold archive (1 year, vector storage)
- L4: Limbo (soft delete, can resurrect)

It's open source and works with Ollama, LM Studio, and any OpenAI-compatible API.

GitHub: https://github.com/MDMAtk/TormentNexus`,

	mcp: `I've cataloged 26,000+ MCP servers in TormentNexus. You can search them at https://tormentnexus.site/catalog

It includes tools for:
- Databases (PostgreSQL, SQLite, etc.)
- Filesystems
- Browsers (Playwright, Chrome DevTools)
- APIs (GitHub, Supabase, Vercel)
- And thousands more

GitHub: https://github.com/MDMAtk/TormentNexus`,

	tools: `TormentNexus is an open-source AI control plane that does exactly this:

- 26,000+ MCP tools in one catalog
- Persistent memory (4-tier system)
- Multi-agent orchestration
- Unified dashboard
- Works with Ollama, LM Studio, DeepSeek

It's a Go backend with a Next.js dashboard. Single binary, runs locally.

GitHub: https://github.com/MDMAtk/TormentNexus`,
};

// Queue file
const QUEUE_FILE = path.join(
	__dirname,
	"..",
	"..",
	"data",
	"reddit-queue.json",
);

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
		rl.question(message + " ", () => {
			rl.close();
			resolve();
		});
	});
}

function calculateRelevance(title) {
	const text = title.toLowerCase();
	let score = 0;
	const matched = [];
	const highMatched = [];

	for (const kw of KEYWORDS) {
		if (text.includes(kw.toLowerCase())) {
			score++;
			matched.push(kw);
		}
	}

	for (const kw of HIGH_KEYWORDS) {
		if (text.includes(kw.toLowerCase())) {
			score += 3;
			highMatched.push(kw);
		}
	}

	return { score, matched, highMatched };
}

function getReplyTemplate(title) {
	const text = title.toLowerCase();

	if (
		text.includes("memory") &&
		(text.includes("forget") ||
			text.includes("remember") ||
			text.includes("persistent"))
	) {
		return TEMPLATES.memory;
	}

	if (
		text.includes("mcp") &&
		(text.includes("tool") ||
			text.includes("server") ||
			text.includes("catalog"))
	) {
		return TEMPLATES.mcp;
	}

	return TEMPLATES.tools;
}

// ── Main Function ──
async function main() {
	console.log("");
	console.log("╔══════════════════════════════════════════════════════════╗");
	console.log("║  TormentNexus Reddit Commenter                          ║");
	console.log("║  Finds relevant posts and adds helpful comments         ║");
	console.log("╚══════════════════════════════════════════════════════════╝");
	console.log("");

	const browser = await puppeteer.launch({
		headless: false,
		args: ["--start-maximized", "--no-sandbox"],
		defaultViewport: null,
	});

	const page = await browser.newPage();

	// Load queue
	let queue = [];
	if (fs.existsSync(QUEUE_FILE)) {
		queue = JSON.parse(fs.readFileSync(QUEUE_FILE, "utf8"));
	}

	console.log("Monitoring subreddits:");
	for (const sub of SUBREDDITS) {
		console.log(`  - r/${sub}`);
	}
	console.log("");
	console.log("Press Ctrl+C to stop at any time.");
	console.log("");

	let cycle = 0;

	while (true) {
		cycle++;
		console.log(`\n${"=".repeat(60)}`);
		console.log(`  Cycle ${cycle} — ${new Date().toLocaleTimeString()}`);
		console.log(`${"=".repeat(60)}`);

		for (const subreddit of SUBREDDITS) {
			console.log(`\n--- r/${subreddit} ---`);

			try {
				// Navigate to subreddit's new posts
				await page.goto(`https://www.reddit.com/r/${subreddit}/new/`, {
					waitUntil: "networkidle2",
					timeout: 30000,
				});

				await sleep(3000);

				// Get all post titles and links
				const posts = await page.evaluate(() => {
					const postElements = document.querySelectorAll(
						'h3, [data-testid="post-content"]',
					);
					const results = [];

					// Try to find post links
					const links = document.querySelectorAll(
						'a[href*="/r/"][href*="/comments/"]',
					);
					for (const link of links) {
						const title = link.textContent?.trim();
						const href = link.href;
						if (title && href && !results.find((r) => r.href === href)) {
							results.push({ title: title.substring(0, 200), href });
						}
					}

					return results.slice(0, 15);
				});

				console.log(`  Found ${posts.length} posts`);

				for (const post of posts) {
					// Check if already in queue
					if (queue.find((q) => q.url === post.href)) {
						continue;
					}

					// Calculate relevance
					const { score, matched, highMatched } = calculateRelevance(
						post.title,
					);

					if (score >= 1) {
						console.log(
							`  [+] Relevant (${score}): ${post.title.substring(0, 60)}...`,
						);
						console.log(`      Keywords: ${matched.slice(0, 3).join(", ")}`);

						const includeLink = highMatched.length > 0;
						const reply = getReplyTemplate(post.title);

						// Add to queue
						queue.push({
							title: post.title,
							url: post.href,
							reply: reply,
							includeLink,
							score,
							keywords: matched,
							timestamp: new Date().toISOString(),
							posted: false,
						});

						// Save queue
						fs.mkdirSync(path.dirname(QUEUE_FILE), { recursive: true });
						fs.writeFileSync(QUEUE_FILE, JSON.stringify(queue, null, 2));

						console.log(`      Reply: ${reply.substring(0, 80)}...`);
						console.log(`      Saved to queue!`);
					}
				}

				await sleep(2000);
			} catch (e) {
				console.log(`  [!] Error: ${e.message}`);
			}
		}

		// Show stats
		const pending = queue.filter((q) => !q.posted).length;
		console.log(`\n--- Stats ---`);
		console.log(`  Total in queue: ${queue.length}`);
		console.log(`  Pending: ${pending}`);

		if (pending > 0) {
			console.log(`\n  [!] ${pending} replies ready to post!`);
			console.log(`  View queue: cat ${QUEUE_FILE}`);
		}

		console.log(`\n  Waiting 5 minutes...`);
		await sleep(300000);
	}
}

main().catch(console.error);
