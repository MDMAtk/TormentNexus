#!/usr/bin/env node
/**
 * TormentNexus Reddit Thread Finder
 * Finds active threads where TormentNexus is relevant to mention
 */

const puppeteer = require("puppeteer");
const fs = require("fs");
const path = require("path");
const https = require("https");

const SUBREDDITS = [
	"LocalLLaMA",
	"MachineLearning",
	"artificial",
	"selfhosted",
	"opensource",
	"AI_Agents",
	"ChatGPT",
	"LocalAI",
	"LangChain",
	"OpenAI",
];

const MIMO_KEY = process.env.MIMO_API_KEY || "";
const MIMO_URL = "https://token-plan-sgp.xiaomimimo.com/v1/chat/completions";
const MIMO_MODEL = "mimo-v2.5";

function sleep(ms) {
	return new Promise((r) => setTimeout(r, ms));
}

function callMiMo(prompt) {
	return new Promise((resolve, reject) => {
		const data = JSON.stringify({
			model: MIMO_MODEL,
			messages: [{ role: "user", content: prompt }],
			max_tokens: 600,
			temperature: 0.5,
		});

		const url = new URL(MIMO_URL);
		const req = https.request(
			{
				hostname: url.hostname,
				port: 443,
				path: url.pathname,
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					Authorization: `Bearer ${MIMO_KEY}`,
					"Content-Length": Buffer.byteLength(data),
				},
			},
			(res) => {
				let body = "";
				res.on("data", (c) => (body += c));
				res.on("end", () => {
					try {
						const r = JSON.parse(body);
						resolve(r.choices?.[0]?.message?.content?.trim() || null);
					} catch (e) {
						reject(e);
					}
				});
			},
		);
		req.on("error", reject);
		req.write(data);
		req.end();
	});
}

async function analyzePost(title, subreddit) {
	const prompt = `Analyze this Reddit post. Would mentioning an open-source project called "TormentNexus" be relevant and helpful?

TormentNexus provides:
- Persistent memory for AI (4-tier system)
- 26,000+ MCP tool catalog
- Multi-agent orchestration
- Works with Ollama, LM Studio, local LLMs
- GitHub: https://github.com/MDMAtk/TormentNexus

Post from r/${subreddit}:
"${title}"

Reply with ONLY this JSON (no markdown):
{"relevant":true/false,"score":1-10,"why":"brief reason","comment":"2-3 sentence organic comment mentioning TormentNexus where it fits naturally"}`;

	try {
		const response = await callMiMo(prompt);
		const match = response?.match(/\{[\s\S]*\}/);
		return match ? JSON.parse(match[0]) : null;
	} catch {
		return null;
	}
}

async function main() {
	console.log("\n=== TormentNexus Reddit Thread Finder ===\n");
	console.log("Scanning subreddits for relevant threads...\n");

	const browser = await puppeteer.launch({ headless: false });
	const page = await browser.newPage();
	const results = [];

	for (const sub of SUBREDDITS) {
		console.log(`Scanning r/${sub}...`);

		try {
			await page.goto(`https://old.reddit.com/r/${sub}/new/`, {
				waitUntil: "domcontentloaded",
				timeout: 20000,
			});
			await sleep(2000);

			const posts = await page.evaluate(() => {
				const items = [];
				document.querySelectorAll(".thing.link").forEach((el) => {
					const title = el.querySelector("a.title")?.textContent?.trim();
					const href = el.querySelector("a.title")?.href;
					const comments = el.querySelector("a.comments")?.textContent?.trim();
					if (title && href) {
						items.push({ title, href, comments: comments || "0 comments" });
					}
				});
				return items.slice(0, 8);
			});

			console.log(`  Found ${posts.length} posts`);

			for (const post of posts) {
				console.log(`  Analyzing: ${post.title.substring(0, 50)}...`);
				const analysis = await analyzePost(post.title, sub);

				if (analysis?.relevant && analysis.score >= 5) {
					console.log(`  ✅ RELEVANT (${analysis.score}/10): ${analysis.why}`);
					results.push({
						subreddit: sub,
						title: post.title,
						url: post.href,
						comments: post.comments,
						score: analysis.score,
						reason: analysis.why,
						comment: analysis.comment,
					});
				} else {
					console.log(`  ❌ Not relevant (${analysis?.score || 0}/10)`);
				}

				await sleep(1000);
			}
		} catch (e) {
			console.log(`  Error: ${e.message}`);
		}

		await sleep(2000);
	}

	await browser.close();

	// Sort by relevance score
	results.sort((a, b) => b.score - a.score);

	// Output results
	console.log("\n" + "═".repeat(70));
	console.log("  RELEVANT THREADS FOR TORMENTNEXUS");
	console.log("═".repeat(70) + "\n");

	if (results.length === 0) {
		console.log("No highly relevant threads found. Try again later.\n");
		return;
	}

	console.log(`Found ${results.length} relevant threads:\n`);

	for (let i = 0; i < results.length; i++) {
		const r = results[i];
		console.log(`${i + 1}. [${r.score}/10] r/${r.subreddit}`);
		console.log(`   Title: ${r.title}`);
		console.log(`   URL: ${r.url}`);
		console.log(`   Activity: ${r.comments}`);
		console.log(`   Why: ${r.reason}`);
		console.log(`\n   Suggested comment:\n   "${r.comment}"\n`);
		console.log("   " + "─".repeat(60) + "\n");
	}

	// Save to file
	const outputFile = path.join(
		__dirname,
		"..",
		"data",
		"reddit-threads-found.json",
	);
	fs.mkdirSync(path.dirname(outputFile), { recursive: true });
	fs.writeFileSync(outputFile, JSON.stringify(results, null, 2));
	console.log(`\nSaved to: ${outputFile}`);
}

main().catch(console.error);
