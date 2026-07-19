#!/usr/bin/env node
const { execSync } = require("child_process");
const os = require("os");
const path = require("path");
const fs = require("fs");

const HOME = os.homedir();
const SYSTEM = process.platform;
const BINARY = SYSTEM === "win32" ? "tormentnexus.exe" : "tormentnexus";
const PORT = process.env.TORMENTNEXUS_PORT || "7778";
const TN_URL = `http://127.0.0.1:${PORT}`;

const argv = process.argv.slice(2);
const cmd = argv[0];

function request(endpoint, method = "GET", data = null) {
	const url = `${TN_URL}/api/${endpoint}`;
	const body = data ? JSON.stringify(data) : null;
	try {
		const result = execSync(
			`curl -s -X ${method} ${url} ${body ? `-H "Content-Type: application/json" -d '${body}'` : ""}`,
			{ encoding: "utf8", timeout: 10000, stdio: ["pipe", "pipe", "pipe"] },
		);
		return JSON.parse(result || "{}");
	} catch {
		return { error: "Kernel not reachable", url: TN_URL };
	}
}

function showHelp() {
	console.log(`
╔══════════════════════════════════════════════╗
║       TormentNexus CLI — tn v1.0.0          ║
╠══════════════════════════════════════════════╣
║  tn search <query>     Search L2 memory     ║
║  tn store <text>       Store a memory       ║
║  tn status             System health check   ║
║  tn tools              List all MCP tools    ║
║  tn tool <name>        Describe a tool      ║
║  tn sessions           Browse past sessions  ║
║  tn harvest            Pull context          ║
║  tn code <query>       Search codebase       ║
║  tn plan               Project plans         ║
║  tn install            Setup TN + MCP        ║
║  tn help               This message          ║
╚══════════════════════════════════════════════╝
  `);
}

async function main() {
	if (!cmd || cmd === "help") return showHelp();

	switch (cmd) {
		case "status": {
			const s = request("runtime/status");
			if (s.error) return console.log(`❌ ${s.error}`);
			const d = s.data;
			console.log(`\n✅ TormentNexus v${d.version} | Uptime: ${d.uptimeSec}s`);
			console.log(
				`   🛠️  ${d.cli.toolCount} tools | 🧠 ${d.memory.l2.count} memories`,
			);
			console.log(`   🔗 ${TN_URL}\n`);
			break;
		}
		case "search": {
			const query = argv.slice(1).join(" ") || argv[1] || "";
			const r = request("memory/search", "POST", { query, limit: 10 });
			const results = r.data || r.results || [];
			if (!results.length)
				return console.log(`No memories found for: "${query}"`);
			results.forEach((m, i) =>
				console.log(
					`[${i + 1}] ${m.content?.slice(0, 120) || JSON.stringify(m).slice(0, 120)}`,
				),
			);
			break;
		}
		case "store": {
			const content = argv.slice(1).join(" ") || "";
			const r = request("memory/store", "POST", {
				content,
				tags: [],
				category: "cli",
			});
			console.log(`🧠 Stored: ${r.id || "OK"}`);
			break;
		}
		case "tools": {
			const r = request("runtime/status");
			const tools = r.data?.cli?.tools || [];
			console.log(`\n🛠️  ${tools.length} MCP tools available:\n`);
			tools.forEach((t) => console.log(`   ${t}`));
			break;
		}
		case "harvest": {
			const prompt = argv.slice(1).join(" ") || "";
			const r = request("context/harvest", "POST", { prompt });
			console.log(JSON.stringify(r, null, 2));
			break;
		}
		case "install": {
			console.log("Running TormentNexus installer...");
			const installer = path.join(
				__dirname,
				"..",
				"..",
				"..",
				"scripts",
				SYSTEM === "win32" ? "install_codewhale.bat" : "install_codewhale.sh",
			);
			if (fs.existsSync(installer)) execSync(installer, { stdio: "inherit" });
			break;
		}
		default:
			console.log(`Unknown command: ${cmd}`);
			showHelp();
	}
}

main().catch(console.error);
