#!/usr/bin/env node
/**
 * TormentNexus NPM Publisher
 * Publishes all @tormentnexus/* packages to npm.
 *
 * PREREQUISITES:
 * 1. Create org on npmjs.com: https://www.npmjs.com/org/create
 *    Name: tormentnexus
 * 2. Create access token: https://www.npmjs.com/settings/<user>/tokens
 *    Type: Automation (no 2FA required for CI/CD)
 *    NO IP restriction (or add your IP)
 * 3. Login: npm login --registry https://registry.npmjs.org
 *    OR set token: export NPM_TOKEN=npm_...
 *
 * USAGE:
 *   node scripts/publish-all-npm.js
 */
const { execSync } = require("child_process");
const path = require("path");
const fs = require("fs");

const BASE = path.join(__dirname, "..", "npm", "@tormentnexus");
const PACKAGES = [
	{ name: "install", dir: "install" },
	{ name: "core", dir: "core" },
	{ name: "cli", dir: "cli" },
	{ name: "openhands", dir: "openhands" },
	{ name: "vscode", dir: "vscode" },
	{ name: "cursor", dir: "cursor" },
];

let published = 0;
let skipped = 0;
let failed = 0;

for (const pkg of PACKAGES) {
	const dir = path.join(BASE, pkg.dir);
	if (!fs.existsSync(dir)) {
		console.log(`⏭️  SKIP @tormentnexus/${pkg.name}: dir not found`);
		skipped++;
		continue;
	}

	const pj = JSON.parse(
		fs.readFileSync(path.join(dir, "package.json"), "utf8"),
	);
	console.log(`\n📦 @tormentnexus/${pkg.name} v${pj.version}`);

	try {
		process.chdir(dir);
		execSync("npm publish --access public --tag alpha", {
			encoding: "utf8",
			timeout: 120000,
			stdio: ["pipe", "pipe", "pipe"],
		});
		console.log(`✅ PUBLISHED`);
		published++;
	} catch (e) {
		const stderr = (e.stderr?.toString() || "") + (e.stdout?.toString() || "");
		if (stderr.includes("previously published")) {
			console.log(`⏭️  Already published (update version to re-publish)`);
			skipped++;
		} else if (stderr.includes("E404") || stderr.includes("404")) {
			console.log(`❌ Org @tormentnexus does not exist on npm!`);
			console.log(`   Create it: https://www.npmjs.com/org/create`);
			failed++;
		} else if (stderr.includes("E401") || stderr.includes("401")) {
			console.log(`❌ Authentication failed. Run: npm login`);
			failed++;
		} else {
			console.log(`❌ FAILED: ${stderr.slice(0, 300)}`);
			failed++;
		}
	}
}

console.log(`\n═══════════════════════════════════`);
console.log(
	`  Published: ${published} | Skipped: ${skipped} | Failed: ${failed}`,
);
console.log(`═══════════════════════════════════`);
