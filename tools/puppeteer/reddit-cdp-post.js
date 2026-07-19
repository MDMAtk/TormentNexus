#!/usr/bin/env node
/**
 * TormentNexus Reddit CDP Poster
 * Automates Reddit posting with proper CDP control
 */

const puppeteer = require("puppeteer");
const fs = require("fs");
const path = require("path");

// Marketing content directory
const AUTO_DIR = path.join(__dirname, "..", "..", "marketing", "auto");

// Find the latest JSON file matching pattern
function findLatestFile(pattern) {
	if (!fs.existsSync(AUTO_DIR)) {
		console.error("Marketing directory not found:", AUTO_DIR);
		process.exit(1);
	}
	const files = fs
		.readdirSync(AUTO_DIR)
		.filter((f) => f.endsWith(".json") && f.includes(pattern))
		.sort()
		.reverse();
	if (files.length === 0) {
		console.error(`No files found matching: ${pattern}`);
		process.exit(1);
	}
	return path.join(AUTO_DIR, files[0]);
}

// Read JSON content
function readContent(file) {
	try {
		return JSON.parse(fs.readFileSync(file, "utf8"));
	} catch (e) {
		console.error("Error reading file:", e.message);
		process.exit(1);
	}
}

// Wait helper
function sleep(ms) {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

// Wait for user input
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

// Post to Reddit with full CDP control
async function postToReddit(subreddit, content) {
	console.log(`\n=== Posting to r/${subreddit} ===\n`);
	console.log("Title:", content.title);
	console.log("Body preview:", content.body.substring(0, 200) + "...");
	console.log("");

	const browser = await puppeteer.launch({
		headless: false,
		args: ["--start-maximized", "--no-sandbox"],
		defaultViewport: null,
	});

	const page = await browser.newPage();

	// Navigate to subreddit
	console.log("Opening r/" + subreddit + "...");
	await page.goto(`https://www.reddit.com/r/${subreddit}/submit`, {
		waitUntil: "networkidle2",
		timeout: 60000,
	});

	await sleep(3000);

	// Wait for user to login if needed
	console.log(
		"\n>>> Login to Reddit if needed, then press ENTER to continue...",
	);
	await waitForEnter("");

	// Click on "Post" tab if not already there
	try {
		const postTab = await page.$(
			'button:has-text("Post"), [aria-label="Post"]',
		);
		if (postTab) {
			await postTab.click();
			await sleep(1000);
		}
	} catch (e) {
		console.log("Post tab not found, continuing...");
	}

	// Inject title
	console.log("Injecting title...");
	try {
		await page.waitForSelector(
			'input[name="title"], textarea[placeholder*="title"], [data-testid="post-title"]',
			{ timeout: 10000 },
		);
		const titleInput = await page.$(
			'input[name="title"], textarea[placeholder*="title"], [data-testid="post-title"]',
		);
		if (titleInput) {
			await titleInput.click({ clickCount: 3 });
			await titleInput.type(content.title, { delay: 10 });
			console.log("Title injected!");
		}
	} catch (e) {
		console.log("Could not find title input, trying alternative...");
		// Try clicking and typing directly
		await page.keyboard.press("Tab");
		await sleep(500);
		await page.keyboard.type(content.title, { delay: 10 });
	}

	await sleep(1000);

	// Switch to markdown mode if available
	console.log("Switching to markdown mode...");
	try {
		const markdownButton = await page.$(
			'button:has-text("Markdown"), [aria-label="Markdown"], button:has-text(" markdown")',
		);
		if (markdownButton) {
			await markdownButton.click();
			await sleep(1000);
			console.log("Switched to markdown mode!");
		}
	} catch (e) {
		console.log("Markdown button not found, continuing...");
	}

	// Inject body (without URLs)
	console.log("Injecting body...");
	try {
		await page.waitForSelector(
			'textarea[placeholder*="body"], [data-testid="post-body"], div[contenteditable="true"]',
			{ timeout: 10000 },
		);
		const bodyInput = await page.$(
			'textarea[placeholder*="body"], [data-testid="post-body"], div[contenteditable="true"]',
		);
		if (bodyInput) {
			await bodyInput.click();
			await sleep(500);

			// Remove URLs from body for the main post
			const bodyText = content.body
				.replace(/https?:\/\/[^\s)]+/g, "") // Remove URLs
				.replace(/\n{3,}/g, "\n\n") // Clean up extra newlines
				.trim();

			await page.keyboard.type(bodyText, { delay: 5 });
			console.log("Body injected!");
		}
	} catch (e) {
		console.log("Could not find body input, trying alternative...");
		await page.keyboard.press("Tab");
		await sleep(500);
		await page.keyboard.type(content.body.replace(/https?:\/\/[^\s)]+/g, ""), {
			delay: 5,
		});
	}

	await sleep(1000);

	// Add flair if available
	console.log("Checking for flair...");
	try {
		const flairButton = await page.$(
			'button:has-text("Flair"), [aria-label="Flair"], button:has-text("Add flair")',
		);
		if (flairButton) {
			await flairButton.click();
			await sleep(1000);

			// Try to select "Show and Tell" or similar flair
			const flairOption = await page.$(
				'button:has-text("Show and Tell"), button:has-text("Project"), button:has-text("Show")',
			);
			if (flairOption) {
				await flairOption.click();
				await sleep(500);
				console.log("Flair added!");
			}
		}
	} catch (e) {
		console.log("Flair not found or not required, continuing...");
	}

	await sleep(1000);

	// Click Post button
	console.log("\n=== READY TO POST ===");
	console.log("Review the post in the browser.");
	console.log("When ready, press ENTER to click Post...");
	await waitForEnter("");

	try {
		const postButton = await page.$(
			'button:has-text("Post"), button[type="submit"], [data-testid="post-button"]',
		);
		if (postButton) {
			await postButton.click();
			console.log("Post button clicked!");
		} else {
			console.log("Post button not found. Please click manually.");
		}
	} catch (e) {
		console.log("Could not click Post button. Please click manually.");
	}

	await sleep(5000);

	// Get the post URL
	const postUrl = page.url();
	console.log("\nPost URL:", postUrl);

	// Add reply with GitHub URL
	console.log("\nAdding reply with GitHub URL...");
	try {
		// Wait for the post page to load
		await page.waitForSelector(
			'textarea[placeholder*="comment"], [data-testid="comment-input"], div[contenteditable="true"]',
			{ timeout: 10000 },
		);
		const commentInput = await page.$(
			'textarea[placeholder*="comment"], [data-testid="comment-input"], div[contenteditable="true"]',
		);
		if (commentInput) {
			await commentInput.click();
			await sleep(500);

			const replyText = `GitHub: https://github.com/MDMAtk/TormentNexus\n\nFeel free to star, fork, or open an issue if you have questions!`;
			await page.keyboard.type(replyText, { delay: 10 });

			await sleep(1000);

			// Click reply button
			const replyButton = await page.$(
				'button:has-text("Reply"), button:has-text("Comment"), button[type="submit"]',
			);
			if (replyButton) {
				await replyButton.click();
				console.log("Reply added!");
			}
		}
	} catch (e) {
		console.log("Could not add reply automatically. Please add manually.");
	}

	console.log("\n=== POST COMPLETE ===");
	console.log("Post URL:", postUrl);
	console.log("Press ENTER to close browser...");
	await waitForEnter("");

	await browser.close();
}

// Main
async function main() {
	const args = process.argv.slice(2);
	const subreddit = args[0] || "MachineLearning";

	console.log("╔══════════════════════════════════════╗");
	console.log("║  TormentNexus Reddit CDP Poster      ║");
	console.log("╚══════════════════════════════════════╝");

	const file = findLatestFile(`reddit_${subreddit}`);
	const content = readContent(file);

	// Clean up content - parse nested JSON if needed
	let title = content.title;
	let body = content.body;

	// If body contains nested JSON, parse it
	try {
		const parsed = JSON.parse(body);
		if (parsed.title) title = parsed.title;
		if (parsed.body) body = parsed.body;
	} catch (e) {
		// Use as-is
	}

	await postToReddit(subreddit, { title, body });
}

main().catch(console.error);
