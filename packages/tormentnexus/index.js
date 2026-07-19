#!/usr/bin/env node

const { execSync, spawn } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");

// Get the binary path
const binDir = path.join(__dirname, "bin");
const isWindows = os.platform() === "win32";
const binaryName = isWindows ? "tormentnexus.exe" : "tormentnexus";
const binaryPath = path.join(binDir, binaryName);

// Check if binary exists
if (!fs.existsSync(binaryPath)) {
	console.error("");
	console.error("  ❌ TormentNexus binary not found!");
	console.error("  Please reinstall: npm install -g tormentnexus");
	console.error("");
	process.exit(1);
}

// Get arguments (skip 'node' and script path)
const args = process.argv.slice(2);

// If no arguments, show help
if (args.length === 0) {
	console.log("");
	console.log("  ⚡ TormentNexus");
	console.log("  ===============");
	console.log("");
	console.log("  Usage:");
	console.log("    tormentnexus serve    Start the server");
	console.log("    tormentnexus mcp      Start MCP server");
	console.log("    tormentnexus --help   Show help");
	console.log("");
	console.log("  Dashboard: http://localhost:7778");
	console.log("");
	process.exit(0);
}

// Run the binary
try {
	const result = spawn(binaryPath, args, {
		stdio: "inherit",
		cwd: process.cwd(),
	});

	result.on("close", (code) => {
		process.exit(code || 0);
	});

	result.on("error", (err) => {
		console.error("");
		console.error("  ❌ Error running TormentNexus:", err.message);
		console.error("");
		process.exit(1);
	});
} catch (err) {
	console.error("");
	console.error("  ❌ Error running TormentNexus:", err.message);
	console.error("");
	process.exit(1);
}
