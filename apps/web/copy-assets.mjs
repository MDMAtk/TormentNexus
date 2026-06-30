#!/usr/bin/env node
// Copies Next.js standalone build output to the Wails frontend dist directory.
// Run AFTER `pnpm build` completes.
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const buildDir = path.resolve(__dirname, ".next-build");
const publicDir = path.resolve(__dirname, "public");
const destDir = path.resolve(
	__dirname,
	"../../go/cmd/tormentnexus-gui/frontend/dist",
);

function copyRecursiveSync(src, dest) {
	const exists = fs.existsSync(src);
	const stats = exists && fs.statSync(src);
	const isDirectory = exists && stats.isDirectory();
	if (isDirectory) {
		if (!fs.existsSync(dest)) {
			fs.mkdirSync(dest, { recursive: true });
		}
		for (const childItemName of fs.readdirSync(src)) {
			if (childItemName === "cache" || childItemName.startsWith("trace"))
				continue;
			copyRecursiveSync(
				path.join(src, childItemName),
				path.join(dest, childItemName),
			);
		}
	} else {
		fs.mkdirSync(path.dirname(dest), { recursive: true });
		fs.copyFileSync(src, dest);
	}
}

if (!fs.existsSync(buildDir)) {
	console.error(`ERROR: ${buildDir} does not exist. Run 'pnpm build' first.`);
	process.exit(1);
}

console.log(`Copying assets from ${buildDir} to ${destDir}...`);
if (fs.existsSync(destDir)) {
	fs.rmSync(destDir, { recursive: true, force: true });
}

// Copy .next-build/static (client JS/CSS chunks)
const staticSrc = path.join(buildDir, "static");
if (fs.existsSync(staticSrc)) {
	copyRecursiveSync(staticSrc, path.join(destDir, "_next", "static"));
}

// Copy standalone .html pages from server/app
const serverAppDir = path.join(buildDir, "server", "app");
if (fs.existsSync(serverAppDir)) {
	function walkHtml(dir, basePath) {
		if (!fs.existsSync(dir)) return;
		for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
			const full = path.join(dir, entry.name);
			if (entry.isDirectory()) {
				walkHtml(full, path.join(basePath, entry.name));
			} else if (entry.name.endsWith(".html")) {
				const relDest = path.join(destDir, basePath, entry.name);
				fs.mkdirSync(path.dirname(relDest), { recursive: true });
				fs.copyFileSync(full, relDest);
			}
		}
	}
	walkHtml(serverAppDir, "");
}

// Copy public/ assets
if (fs.existsSync(publicDir)) {
	copyRecursiveSync(publicDir, destDir);
}

// Copy the root HTML page
const rootHtml = path.join(serverAppDir, "index.html");
if (fs.existsSync(rootHtml)) {
	fs.copyFileSync(rootHtml, path.join(destDir, "index.html"));
}

// Ensure placeholder exists
fs.mkdirSync(destDir, { recursive: true });
fs.writeFileSync(path.join(destDir, "placeholder.txt"), "placeholder");

console.log("Frontend assets copy complete!");
