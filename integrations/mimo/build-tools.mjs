#!/usr/bin/env node
/**
 * Bundles MiMo Desktop tool sources with zod inlined (same format as built-in tools).
 * Input:  integrations/mimo/tools/src/*.js  (ESM with `import { z } from "zod"`)
 * Output: integrations/mimo/tools/dist/*.js (single-file ESM, zod bundled)
 *
 * Usage: node integrations/mimo/build-tools.mjs
 */
import { readdirSync, mkdirSync, existsSync } from "node:fs";
import { join, dirname, basename } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

const ROOT = dirname(fileURLToPath(import.meta.url));
const SRC = join(ROOT, "tools", "src");
const DIST = join(ROOT, "tools", "dist");

// Locate zod + esbuild from MiMoCode / HyperNexus node_modules
const require = (await import("node:module")).createRequire(import.meta.url);

let zodPath;
try {
  zodPath = dirname(require.resolve("zod/package.json", { paths: [
    "C:/Users/hyper/.config/mimocode",
    join(ROOT, "..", ".."),
    ROOT,
  ]}));
} catch {
  console.error("zod not found — run: npm i zod (in ~/.config/mimocode or HyperNexus root)");
  process.exit(1);
}

let esbuild;
try {
  esbuild = await import(pathToFileURL(require.resolve("esbuild", { paths: [
    join(ROOT, "..", ".."),
    ROOT,
  ]})).href);
} catch {
  console.error("esbuild not found — run: npm i esbuild");
  process.exit(1);
}

console.log("zod:", zodPath);
console.log("esbuild:", esbuild.version || "ok");

mkdirSync(DIST, { recursive: true });

const files = readdirSync(SRC).filter((f) => f.endsWith(".js"));
console.log(`Building ${files.length} tools...`);

for (const f of files) {
  const entry = join(SRC, f);
  const outfile = join(DIST, f);
  try {
    await esbuild.build({
      entryPoints: [entry],
      bundle: true,
      format: "esm",
      platform: "neutral",
      outfile,
      alias: { zod: join(zodPath, "index.js") },
      logLevel: "warning",
    });
    console.log(`  ✓ ${basename(outfile)}`);
  } catch (e) {
    console.error(`  ✗ ${f}:`, e.message);
    process.exitCode = 1;
  }
}

console.log(`\nDone. Output: ${DIST}`);
