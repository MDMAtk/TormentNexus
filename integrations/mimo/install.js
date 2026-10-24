#!/usr/bin/env node
/**
 * HyperNexus MiMo Addon Installer
 * Installs skills, tools, commands, hooks, and agents into:
 *   - MiMoCode CLI  (~/.config/mimocode/)
 *   - MiMo Desktop  (AppData/Roaming/Xiaomi MiMo AI/engine-config/)
 *   - Agents compat (~/.agents/skills/)
 *
 * Usage: node integrations/mimo/install.js [--dry-run]
 */
import { cpSync, mkdirSync, existsSync, readFileSync, writeFileSync, readdirSync } from "node:fs";
import { join, dirname, basename } from "node:path";
import { fileURLToPath } from "node:url";
import { homedir, platform, arch } from "node:os";

const ROOT = dirname(fileURLToPath(import.meta.url));
const DRY = process.argv.includes("--dry-run");
const HOME = homedir();
const WIN = platform() === "win32";

function log(msg) {
  console.log("  " + msg);
}
function copyDir(src, dst) {
  if (DRY) {
    log(`[dry] copy ${src} -> ${dst}`);
    return;
  }
  mkdirSync(dst, { recursive: true });
  cpSync(src, dst, { recursive: true, force: true });
  log(`copied ${basename(src)}/ -> ${dst}`);
}
function copyFile(src, dst) {
  if (DRY) {
    log(`[dry] copy ${src} -> ${dst}`);
    return;
  }
  mkdirSync(dirname(dst), { recursive: true });
  cpSync(src, dst, { force: true });
  log(`copied ${basename(src)} -> ${dst}`);
}
function writeJson(path, data) {
  if (DRY) {
    log(`[dry] write ${path}`);
    return;
  }
  mkdirSync(dirname(path), { recursive: true });
  writeFileSync(path, JSON.stringify(data, null, 2) + "\n", "utf8");
  log(`wrote ${path}`);
}
function mergeJsoncMcp(path, entry) {
  let raw = existsSync(path) ? readFileSync(path, "utf8") : "{}";
  // Strip // comments for parsing (best-effort JSONC -> JSON)
  const stripped = raw
    .replace(/^\s*\/\/.*$/gm, "")
    .replace(/\/\*[\s\S]*?\*\//g, "")
    .replace(/,\s*([}\]])/g, "$1");
  let obj = {};
  try {
    obj = JSON.parse(stripped);
  } catch {
    obj = {};
  }
  obj.mcp = obj.mcp || {};
  obj.mcp.hypernexus = entry;
  // Write back as plain JSON (readable, valid JSONC superset)
  writeJson(path, obj);
}

// ─── Paths ────────────────────────────────────────────────────────────────
const MIMOCODE_HOME = join(HOME, ".config", "mimocode");
const DESKTOP_ROOT = WIN
  ? join(HOME, "AppData", "Roaming", "Xiaomi MiMo AI", "engine-config")
  : join(HOME, "Library", "Application Support", "Xiaomi MiMo AI", "engine-config");
const AGENTS_ROOT = join(HOME, ".agents", "skills");

const SRC_SKILL = join(ROOT, "skills", "hypernexus");
const SRC_TOOLS = join(ROOT, "tools", "dist");
const SRC_COMMANDS = join(ROOT, "mimocode", "commands");
const SRC_HOOKS = join(ROOT, "mimocode", "hooks");
const SRC_AGENTS = join(ROOT, "mimocode", "agents");

// ─── MCP entry for MiMoCode ───────────────────────────────────────────────
function findBinary() {
  const candidates = WIN
    ? [
        join(HOME, "AppData", "Local", "TormentNexus", "tormentnexus.exe"),
        join(HOME, ".tormentnexus", "tormentnexus.exe"),
        join(HOME, ".tormentnexus", "bin", "tormentnexus.exe"),
        join(HOME, "AppData", "Local", "HyperNexus", "hypernexus.exe"),
        join(HOME, ".hypernexus", "hypernexus.exe"),
        join(HOME, ".hypernexus", "bin", "hypernexus.exe"),
        "tormentnexus.exe",
        "hypernexus.exe",
      ]
    : [
        join(HOME, ".tormentnexus", "tormentnexus"),
        join(HOME, ".tormentnexus", "bin", "tormentnexus"),
        "/usr/local/bin/tormentnexus",
        join(HOME, ".hypernexus", "hypernexus"),
        join(HOME, ".hypernexus", "bin", "hypernexus"),
        "/usr/local/bin/hypernexus",
        "tormentnexus",
        "hypernexus",
      ];
  for (const c of candidates) {
    if (c === "tormentnexus.exe" || c === "tormentnexus" || c === "hypernexus.exe" || c === "hypernexus") return c;
    if (existsSync(c)) return c;
  }
  return candidates[candidates.length - 1];
}

const BINARY = findBinary();
const WORKSPACE_ROOT = process.env.HYPERNEXUS_WORKSPACE_ROOT || join(HOME, "workspace");
const MCP_ENTRY = {
  type: "local",
  command: [BINARY, "mcp"],
  env: {
    HYPERNEXUS_WORKSPACE_ROOT: WORKSPACE_ROOT,
    HN_EDITION: process.env.HN_EDITION || "corporate",
  },
};

// ─── Install ──────────────────────────────────────────────────────────────
console.log(`\nHyperNexus MiMo Addon Installer${DRY ? " (dry run)" : ""}`);
console.log(`  Platform: ${platform()} (${arch()})`);
console.log(`  Binary:   ${BINARY}`);
console.log(`  Edition:  ${MCP_ENTRY.env.HN_EDITION}\n`);

let count = 0;

// 1. MiMoCode skill + commands + hooks + agents
console.log("[1/5] MiMoCode (CLI)");
copyDir(SRC_SKILL, join(MIMOCODE_HOME, "skills", "hypernexus"));
count++;
if (existsSync(SRC_COMMANDS)) {
  copyDir(SRC_COMMANDS, join(MIMOCODE_HOME, "commands"));
  count++;
}
if (existsSync(SRC_HOOKS)) {
  copyDir(SRC_HOOKS, join(MIMOCODE_HOME, "hooks"));
  count++;
}
if (existsSync(SRC_AGENTS)) {
  copyDir(SRC_AGENTS, join(MIMOCODE_HOME, "agents"));
  count++;
}
mergeJsoncMcp(join(MIMOCODE_HOME, "mimocode.jsonc"), MCP_ENTRY);
count++;

// 2. MiMo Desktop skill + tools
console.log("[2/5] MiMo Desktop");
if (existsSync(DESKTOP_ROOT) || true) {
  copyDir(SRC_SKILL, join(DESKTOP_ROOT, "skills", "hypernexus"));
  count++;
  if (existsSync(SRC_TOOLS)) {
    const tools = readdirSync(SRC_TOOLS).filter((f) => f.endsWith(".js"));
    for (const t of tools) {
      copyFile(join(SRC_TOOLS, t), join(DESKTOP_ROOT, "tools", t));
    }
    log(`${tools.length} tools installed to Desktop engine-config`);
    count++;
  } else {
    log("WARN: tools/dist missing — run build-tools.mjs first");
  }
}

// 3. Agents compat skill
console.log("[3/5] Agents skill");
copyDir(SRC_SKILL, join(AGENTS_ROOT, "hypernexus"));
count++;

// 4. Desktop skill already includes locales (Plugins page metadata)
console.log("[4/5] Skill locales (Desktop Plugins page)");
log("en-US.json + zh-CN.json included in skill dir");

// 5. Summary
console.log("[5/5] Done\n");
console.log(`  ${count} surfaces updated`);
console.log(`  MiMoCode skill:   ${join(MIMOCODE_HOME, "skills", "hypernexus")}`);
console.log(`  MiMoCode MCP:     ${join(MIMOCODE_HOME, "mimocode.jsonc")}`);
console.log(`  Desktop skill:    ${join(DESKTOP_ROOT, "skills", "hypernexus")}`);
console.log(`  Desktop tools:    ${join(DESKTOP_ROOT, "tools")}`);
console.log(`  Agents skill:     ${join(AGENTS_ROOT, "hypernexus")}`);
console.log(`\nRestart MiMoCode / MiMo Desktop (new conversations) to pick up changes.\n`);
