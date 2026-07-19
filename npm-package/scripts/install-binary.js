#!/usr/bin/env node

/**
 * TormentNexus npm installer
 * Downloads the appropriate binary for the current platform
 */

const https = require("https");
const fs = require("fs");
const path = require("path");
const os = require("os");
const { execSync } = require("child_process");

const VERSION = "1.0.0";
const BASE_URL = `https://github.com/MDMAtk/TormentNexus/releases/download/v${VERSION}`;

function getPlatform() {
	const platform = os.platform();
	const arch = os.arch();

	if (platform === "darwin") {
		return arch === "arm64" ? "darwin-arm64" : "darwin-amd64";
	} else if (platform === "linux") {
		return arch === "arm64" ? "linux-arm64" : "linux-amd64";
	} else if (platform === "win32") {
		return "windows-amd64";
	}

	throw new Error(`Unsupported platform: ${platform}-${arch}`);
}

function getBinaryName() {
	const platform = os.platform();
	return platform === "win32" ? "tormentnexus.exe" : "tormentnexus";
}

function download(url, dest) {
	return new Promise((resolve, reject) => {
		const file = fs.createWriteStream(dest);

		https
			.get(url, (response) => {
				// Handle redirects
				if (response.statusCode === 302 || response.statusCode === 301) {
					download(response.headers.location, dest).then(resolve).catch(reject);
					return;
				}

				if (response.statusCode !== 200) {
					reject(new Error(`Download failed: ${response.statusCode}`));
					return;
				}

				response.pipe(file);

				file.on("finish", () => {
					file.close();
					resolve();
				});
			})
			.on("error", (err) => {
				fs.unlink(dest, () => {});
				reject(err);
			});
	});
}

async function main() {
	console.log("Installing TormentNexus...");

	const platform = getPlatform();
	const binaryName = getBinaryName();
	const downloadUrl = `${BASE_URL}/tormentnexus-${platform}.tar.gz`;

	// Create bin directory
	const binDir = path.join(__dirname, "..", "bin");
	if (!fs.existsSync(binDir)) {
		fs.mkdirSync(binDir, { recursive: true });
	}

	const binaryPath = path.join(binDir, binaryName);

	// Check if binary already exists
	if (fs.existsSync(binaryPath)) {
		console.log("Binary already exists, skipping download...");
		return;
	}

	console.log(`Downloading TormentNexus for ${platform}...`);
	console.log(`URL: ${downloadUrl}`);

	// Download tar.gz
	const tarPath = path.join(binDir, "tormentnexus.tar.gz");

	try {
		await download(downloadUrl, tarPath);

		// Extract
		console.log("Extracting...");
		execSync(`tar -xzf "${tarPath}" -C "${binDir}"`, { stdio: "inherit" });

		// Make executable (Unix only)
		if (os.platform() !== "win32") {
			fs.chmodSync(binaryPath, "755");
		}

		// Clean up tar
		fs.unlinkSync(tarPath);

		console.log("TormentNexus installed successfully!");
		console.log("");
		console.log("To get started:");
		console.log("  tormentnexus serve");
		console.log("");
		console.log("Dashboard: http://127.0.0.1:7778");
	} catch (err) {
		console.error("Installation failed:", err.message);
		console.error("");
		console.error("Please download manually from:");
		console.error(`  https://github.com/MDMAtk/TormentNexus/releases`);
		process.exit(1);
	}
}

main();
