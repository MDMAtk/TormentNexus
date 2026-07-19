#!/usr/bin/env node
const { execSync } = require("child_process");
const path = require("path");

console.log("╔══════════════════════════════════════════╗");
console.log("║   TormentNexus Universal Installer      ║");
console.log("║   38 AI Clients • One Command           ║");
console.log("╚══════════════════════════════════════════╝\n");

const installer = path.join(
	__dirname,
	"..",
	"..",
	"..",
	"scripts",
	"install-client-support.py",
);
try {
	execSync(`python3 "${installer}"`, { stdio: "inherit" });
} catch {
	execSync(`python "${installer}"`, { stdio: "inherit" });
}

console.log(
	"\n✅ Done! TormentNexus installed for all AI clients on your system.",
);
console.log("   Run 'tn status' to check your memory system.");
console.log("   Run 'tn search \"your topic\"' to search memories.\n");
