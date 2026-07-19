#!/usr/bin/env node

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");
const os = require("os");

const VERSION = "1.0.0-b13";
const REPO = "MDMAtk/TormentNexus";

function getPlatform() {
	const platform = os.platform();
	const arch = os.arch();

	if (platform === "win32") {
		return { name: "windows", ext: ".zip", binary: "tormentnexus.exe" };
	} else if (platform === "darwin") {
		const suffix = arch === "arm64" ? "arm64" : "amd64";
		return {
			name: `darwin-${suffix}`,
			ext: ".tar.gz",
			binary: "tormentnexus",
		};
	} else if (platform === "linux") {
		const suffix = arch === "arm64" ? "arm64" : "amd64";
		return {
			name: `linux-${suffix}`,
			ext: ".tar.gz",
			binary: "tormentnexus",
		};
	}

	throw new Error(`Unsupported platform: ${platform}-${arch}`);
}

function downloadFile(url, dest) {
	return new Promise((resolve, reject) => {
		const file = fs.createWriteStream(dest);

		const request = https.get(url, (response) => {
			if (response.statusCode === 302 || response.statusCode === 301) {
				downloadFile(response.headers.location, dest)
					.then(resolve)
					.catch(reject);
				return;
			}

			if (response.statusCode !== 200) {
				reject(new Error(`Download failed: ${response.statusCode}`));
				return;
			}

			const totalBytes = parseInt(response.headers["content-length"], 10);
			let downloadedBytes = 0;

			response.on("data", (chunk) => {
				downloadedBytes += chunk.length;
				const progress = ((downloadedBytes / totalBytes) * 100).toFixed(1);
				process.stdout.write(`\r  Downloading: ${progress}%`);
			});

			response.pipe(file);
		});

		file.on("finish", () => {
			file.close();
			console.log("");
			resolve();
		});

		file.on("error", (err) => {
			fs.unlink(dest, () => {});
			reject(err);
		});

		request.on("error", (err) => {
			reject(err);
		});
	});
}

async function extractArchive(archivePath, destDir, platform) {
	try {
		if (platform.ext === ".zip") {
			execSync(
				`powershell -command "Expand-Archive -Path '${archivePath}' -DestinationPath '${destDir}' -Force"`,
				{ stdio: "inherit" },
			);
		} else {
			execSync(`tar -xzf "${archivePath}" -C "${destDir}"`, {
				stdio: "inherit",
			});
		}
	} catch (error) {
		throw new Error(`Failed to extract archive: ${error.message}`);
	}
}

async function install() {
	console.log("");
	console.log("  ⚡ TormentNexus Installer");
	console.log("  =========================");
	console.log("");

	try {
		const platform = getPlatform();
		console.log(`  Platform: ${platform.name}`);

		const packageName = `tormentnexus-${platform.name}${platform.ext}`;
		const downloadUrl = `https://github.com/${REPO}/releases/download/${VERSION}/${packageName}`;

		console.log(`  Downloading from GitHub releases...`);
		console.log("");

		const binDir = path.join(__dirname, "..", "bin");
		if (!fs.existsSync(binDir)) {
			fs.mkdirSync(binDir, { recursive: true });
		}

		const archivePath = path.join(binDir, packageName);
		await downloadFile(downloadUrl, archivePath);

		console.log("  Extracting...");
		await extractArchive(archivePath, binDir, platform);

		// Make binary executable on Unix
		if (platform.ext !== ".zip") {
			const binaryPath = path.join(binDir, platform.binary);
			fs.chmodSync(binaryPath, 0o755);
		}

		// Clean up archive
		fs.unlinkSync(archivePath);

		console.log("");
		console.log("  ✅ Installation complete!");
		console.log("");
		console.log("  Run: npx tormentnexus serve");
		console.log("  Or:  tormentnexus serve");
		console.log("");
		console.log("  Dashboard: http://localhost:7778");
		console.log("");
	} catch (error) {
		console.error("");
		console.error("  ❌ Installation failed:", error.message);
		console.error("");
		console.error("  Please download manually from:");
		console.error(`  https://github.com/${REPO}/releases`);
		console.error("");
		process.exit(1);
	}
}

install();
