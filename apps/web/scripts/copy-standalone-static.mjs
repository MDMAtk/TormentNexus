import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const webRoot = path.resolve(__dirname, "..");
const buildDir = path.join(webRoot, ".next-build");
const staticSrc = path.join(buildDir, "static");

const standaloneRoot = path.join(buildDir, "standalone");
const targetNextBuildStatic = path.join(standaloneRoot, "apps", "web", ".next-build", "static");
const targetNextStatic = path.join(standaloneRoot, "apps", "web", ".next", "static");
const rootNextBuildStatic = path.join(standaloneRoot, ".next-build", "static");
const rootNextStatic = path.join(standaloneRoot, ".next", "static");

function copyRecursiveSync(src, dest) {
	const exists = fs.existsSync(src);
	const stats = exists && fs.statSync(src);
	const isDirectory = exists && stats.isDirectory();
	if (isDirectory) {
		if (!fs.existsSync(dest)) {
			fs.mkdirSync(dest, { recursive: true });
		}
		for (const childItemName of fs.readdirSync(src)) {
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

if (fs.existsSync(staticSrc)) {
	console.log(`[copy-static] Copying static assets from ${staticSrc} to standalone directories...`);
	copyRecursiveSync(staticSrc, targetNextBuildStatic);
	copyRecursiveSync(staticSrc, targetNextStatic);
	copyRecursiveSync(staticSrc, rootNextBuildStatic);
	copyRecursiveSync(staticSrc, rootNextStatic);
	console.log("[copy-static] Static assets copied successfully!");
} else {
	console.warn("[copy-static] Warning: Source static directory not found.");
}
