#!/usr/bin/env node
/**
 * TormentNexus Reddit Post Comments
 * Posts comments from the queue to Reddit
 */

const puppeteer = require("puppeteer");
const fs = require("fs");
const path = require("path");

const QUEUE_FILE = path.join(__dirname, "..", "data", "reddit-queue.json");

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

async function main() {
	console.log("");
	console.log("╔══════════════════════════════════════════════════════════╗");
	console.log("║  TormentNexus Reddit Comment Poster                     ║");
	console.log("║  Posts comments from the queue                          ║");
	console.log("╚══════════════════════════════════════════════════════════╝");
	console.log("");

	// Load queue
	if (!fs.existsSync(QUEUE_FILE)) {
		console.log("  No queue file found. Run the scanner first.");
		return;
	}

	const queue = JSON.parse(fs.readFileSync(QUEUE_FILE, "utf8"));
	const pending = queue.filter((q) => !q.posted);

	if (pending.length === 0) {
		console.log("  No pending comments to post.");
		return;
	}

	console.log(`  Found ${pending.length} pending comments:`);
	console.log("");

	for (let i = 0; i < pending.length; i++) {
		const item = pending[i];
		console.log(
			`  ${i + 1}. [${item.subreddit}] ${item.title.substring(0, 60)}...`,
		);
		console.log(`     URL: ${item.url}`);
		console.log(
			`     Score: ${item.score} | Keywords: ${item.keywords.slice(0, 3).join(", ")}`,
		);
		console.log("");
	}

	const browser = await puppeteer.launch({
		headless: false,
		args: ["--start-maximized"],
		defaultViewport: null,
	});

	const page = await browser.newPage();

	for (let i = 0; i < pending.length; i++) {
		const item = pending[i];

		console.log(`\n${"=".repeat(60)}`);
		console.log(
			`  Posting ${i + 1}/${pending.length}: ${item.title.substring(0, 50)}...`,
		);
		console.log(`${"=".repeat(60)}`);

		// Navigate to the post
		console.log(`  Opening: ${item.url}`);
		await page.goto(item.url, { waitUntil: "networkidle2", timeout: 30000 });
		await sleep(3000);

		// Show the comment
		console.log(`\n  Comment to post:`);
		console.log(`  ${"─".repeat(50)}`);
		console.log(`  ${item.reply}`);
		console.log(`  ${"─".repeat(50)}`);

		// Ask user to post
		console.log(`\n  Press ENTER to post this comment, or 's' to skip:`);
		const answer = await waitForEnter("");

		if (answer && answer.toLowerCase() === "s") {
			console.log("  Skipped.");
			continue;
		}

		try {
			// Find the comment box
			console.log("  Looking for comment box...");

			// Try different selectors for the comment box
			const commentSelectors = [
				'div[data-testid="comment-composer"]',
				'div[contenteditable="true"]',
				'textarea[placeholder*="comment"]',
				'textarea[placeholder*="What are your thoughts"]',
				".public-DraftEditor-content",
				'div[role="textbox"]',
			];

			let commentBox = null;
			for (const selector of commentSelectors) {
				commentBox = await page.$(selector);
				if (commentBox) {
					console.log(`  Found comment box: ${selector}`);
					break;
				}
			}

			if (!commentBox) {
				console.log("  [!] Could not find comment box.");
				console.log("  Please paste the comment manually.");
				await waitForEnter("  Press ENTER when done...");
			} else {
				// Click the comment box
				await commentBox.click();
				await sleep(500);

				// Type the comment
				console.log("  Typing comment...");
				await page.keyboard.type(item.reply, { delay: 5 });
				await sleep(1000);

				// Find and click the submit button
				console.log("  Looking for submit button...");
				const submitSelectors = [
					'button:has-text("Comment")',
					'button[type="submit"]',
					'button:has-text("Reply")',
					'button[data-testid="comment-submit-button"]',
				];

				let submitButton = null;
				for (const selector of submitSelectors) {
					submitButton = await page.$(selector);
					if (submitButton) break;
				}

				if (submitButton) {
					console.log("  Clicking submit...");
					await submitButton.click();
					await sleep(3000);

					// Mark as posted
					item.posted = true;
					console.log("  [+] Comment posted!");
				} else {
					console.log("  [!] Could not find submit button.");
					console.log("  Please click submit manually.");
					await waitForEnter("  Press ENTER when done...");
				}
			}
		} catch (e) {
			console.log(`  [!] Error: ${e.message}`);
			console.log("  Please post the comment manually.");
			await waitForEnter("  Press ENTER when done...");
		}

		// Save queue
		fs.writeFileSync(QUEUE_FILE, JSON.stringify(queue, null, 2));

		// Wait between posts
		if (i < pending.length - 1) {
			console.log(`\n  Waiting 30 seconds before next post...`);
			await sleep(30000);
		}
	}

	console.log(`\n${"=".repeat(60)}`);
	console.log("  Done! All comments processed.");
	console.log(`${"=".repeat(60)}`);

	await browser.close();
}

main().catch(console.error);
