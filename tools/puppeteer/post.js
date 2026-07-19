#!/usr/bin/env node
/**
 * TormentNexus — Semi-Automated Poster (Local Version)
 * Opens browser locally, fills in content, waits for user review before posting
 *
 * Usage:
 *   node post.js reddit <subreddit>
 *   node post.js hn
 *   node post.js twitter
 */

const puppeteer = require("puppeteer");
const fs = require("fs");
const path = require("path");

// Marketing content directory (relative to project root)
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

// Post to Reddit
async function postToReddit(subreddit) {
	console.log(`\n=== Posting to r/${subreddit} ===\n`);

	const file = findLatestFile(`reddit_${subreddit}`);
	const content = readContent(file);

	console.log("Title:", content.title);
	console.log("Body preview:", content.body.substring(0, 200) + "...");
	console.log("");

	const browser = await puppeteer.launch({
		headless: false,
		args: ["--start-maximized"],
	});

	const page = await browser.newPage();
	await page.setViewport({ width: 1280, height: 900 });

	// Navigate to Reddit submit page
	console.log("Opening Reddit...");
	await page.goto(`https://www.reddit.com/r/${subreddit}/submit`, {
		waitUntil: "networkidle2",
		timeout: 60000,
	});

	// Wait for user to log in if needed
	await waitForEnter(
		"\n>>> Log in to Reddit if needed, then press ENTER to continue...",
	);

	// Select "Text" post type
	try {
		await page.waitForSelector('button, [role="button"]', { timeout: 5000 });
		const buttons = await page.$$('button, [role="button"]');
		for (const btn of buttons) {
			const text = await page.evaluate((el) => el.textContent, btn);
			if (text && text.includes("Post")) {
				await btn.click();
				break;
			}
		}
	} catch (e) {
		console.log("Note: You may need to manually select 'Text' post type");
	}

	await new Promise((r) => setTimeout(r, 2000));

	// Fill in title
	console.log("Filling in title...");
	try {
		const titleInput = await page.$('textarea, input[type="text"]');
		if (titleInput) {
			await titleInput.click();
			await titleInput.type(content.title.substring(0, 300));
		}
	} catch (e) {
		console.log("Note: You may need to manually fill in the title");
	}

	// Fill in body
	console.log("Filling in body...");
	try {
		const bodyInputs = await page.$$('textarea, [contenteditable="true"]');
		if (bodyInputs.length > 1) {
			await bodyInputs[1].click();
			// Type body in chunks
			const chunks = content.body.match(/.{1,100}/g) || [content.body];
			for (const chunk of chunks) {
				await page.keyboard.type(chunk, { delay: 5 });
			}
		}
	} catch (e) {
		console.log("Note: You may need to manually fill in the body");
	}

	console.log("\n✅ Content filled in!");
	await waitForEnter(
		"\n>>> Review the post in the browser. Press ENTER to submit...",
	);

	// Click submit
	try {
		const submitButtons = await page.$$('button[type="submit"], button');
		for (const btn of submitButtons) {
			const text = await page.evaluate((el) => el.textContent, btn);
			if (text && (text.includes("Post") || text.includes("Submit"))) {
				await btn.click();
				console.log("✅ Post submitted!");
				break;
			}
		}
	} catch (e) {
		console.log("⚠️ Could not find submit button. Please click it manually.");
	}

	await waitForEnter("\n>>> Press ENTER to close the browser...");
	await browser.close();
}

// Post to Hacker News
async function postToHN() {
	console.log("\n=== Posting to Hacker News ===\n");

	const file = findLatestFile("hackernews");
	const content = readContent(file);

	// Parse the nested JSON
	let title, body;
	try {
		const parsed = JSON.parse(content.body);
		title = parsed.title || content.title;
		body = parsed.body || content.body;
	} catch {
		title = content.title;
		body = content.body;
	}

	console.log("Title:", title);
	console.log("Body preview:", body.substring(0, 200) + "...");
	console.log("");

	const browser = await puppeteer.launch({
		headless: false,
		args: ["--start-maximized"],
	});

	const page = await browser.newPage();
	await page.setViewport({ width: 1280, height: 900 });

	// Navigate to HN submit page
	console.log("Opening Hacker News...");
	await page.goto("https://news.ycombinator.com/submit", {
		waitUntil: "networkidle2",
		timeout: 60000,
	});

	// Wait for user to log in if needed
	await waitForEnter(
		"\n>>> Log in to HN if needed, then press ENTER to continue...",
	);

	// Fill in title
	console.log("Filling in title...");
	const titleInput = await page.$('input[name="title"]');
	if (titleInput) {
		await titleInput.click();
		await titleInput.type(title.substring(0, 80));
	}

	// Fill in URL
	console.log("Filling in URL...");
	const urlInput = await page.$('input[name="url"]');
	if (urlInput) {
		await urlInput.click();
		await urlInput.type("https://github.com/MDMAtk/TormentNexus");
	}

	// Fill in text
	console.log("Filling in text...");
	const textInput = await page.$('textarea[name="text"]');
	if (textInput) {
		await textInput.click();
		const chunks = body.match(/.{1,100}/g) || [body];
		for (const chunk of chunks) {
			await page.keyboard.type(chunk, { delay: 5 });
		}
	}

	console.log("\n✅ Content filled in!");
	await waitForEnter(
		"\n>>> Review the post in the browser. Press ENTER to submit...",
	);

	// Click submit
	const submitButton = await page.$('input[type="submit"]');
	if (submitButton) {
		await submitButton.click();
		console.log("✅ Post submitted!");
	}

	await waitForEnter("\n>>> Press ENTER to close the browser...");
	await browser.close();
}

// Post to Twitter
async function postToTwitter() {
	console.log("\n=== Posting Twitter Thread ===\n");

	const file = findLatestFile("twitter");
	const content = readContent(file);

	console.log("Thread:");
	content.tweets.forEach((t, i) => console.log(`  ${i + 1}. ${t}`));
	console.log("");

	const browser = await puppeteer.launch({
		headless: false,
		args: ["--start-maximized"],
	});

	const page = await browser.newPage();
	await page.setViewport({ width: 1280, height: 900 });

	// Navigate to Twitter
	console.log("Opening Twitter...");
	await page.goto("https://twitter.com/compose/tweet", {
		waitUntil: "networkidle2",
		timeout: 60000,
	});

	// Wait for user to log in
	await waitForEnter(
		"\n>>> Log in to Twitter if needed, then press ENTER to continue...",
	);

	for (let i = 0; i < content.tweets.length; i++) {
		console.log(`\nPosting tweet ${i + 1}/${content.tweets.length}...`);
		console.log(`"${content.tweets[i].substring(0, 50)}..."`);

		// Type tweet
		try {
			const tweetInput = await page.$('[contenteditable="true"], textarea');
			if (tweetInput) {
				await tweetInput.click();
				await page.keyboard.type(content.tweets[i], { delay: 20 });
			}
		} catch (e) {
			console.log("Note: You may need to manually type the tweet");
		}

		if (i < content.tweets.length - 1) {
			await waitForEnter(
				`>>> Review tweet ${i + 1}. Press ENTER to post and continue...`,
			);

			// Click tweet button
			try {
				const buttons = await page.$$("button");
				for (const btn of buttons) {
					const text = await page.evaluate((el) => el.textContent, btn);
					if (text && text.includes("Post")) {
						await btn.click();
						break;
					}
				}
			} catch (e) {
				console.log("Note: Please click the Post button manually");
			}

			// Wait for tweet to post
			await new Promise((r) => setTimeout(r, 3000));

			// Start new tweet
			try {
				const newTweetBtn = await page.$('[data-testid="tweetButton"]');
				if (newTweetBtn) await newTweetBtn.click();
			} catch (e) {
				// Continue
			}
		} else {
			await waitForEnter(`>>> Review final tweet. Press ENTER to post...`);
			try {
				const buttons = await page.$$("button");
				for (const btn of buttons) {
					const text = await page.evaluate((el) => el.textContent, btn);
					if (text && text.includes("Post")) {
						await btn.click();
						break;
					}
				}
			} catch (e) {
				console.log("Note: Please click the Post button manually");
			}
		}
	}

	console.log("\n✅ Thread posted!");
	await waitForEnter("\n>>> Press ENTER to close the browser...");
	await browser.close();
}

// Main
const args = process.argv.slice(2);
const command = args[0];

console.log("╔══════════════════════════════════════╗");
console.log("║  TormentNexus Marketing Poster       ║");
console.log("╚══════════════════════════════════════╝");
console.log("");

switch (command) {
	case "reddit":
		postToReddit(args[1] || "MachineLearning");
		break;
	case "hn":
		postToHN();
		break;
	case "twitter":
		postToTwitter();
		break;
	default:
		console.log("Usage:");
		console.log("  node post.js reddit <subreddit>");
		console.log("  node post.js hn");
		console.log("  node post.js twitter");
		console.log("");
		console.log("Examples:");
		console.log("  node post.js reddit MachineLearning");
		console.log("  node post.js reddit LocalLLaMA");
		console.log("  node post.js hn");
		console.log("  node post.js twitter");
}
