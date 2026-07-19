#!/usr/bin/env node
const http = require("http");

const PORT = process.env.TORMENTNEXUS_PORT || "7778";
const TN_URL = `http://127.0.0.1:${PORT}`;

const argv = process.argv.slice(2);
const cmd = argv[0];

function request(endpoint, method = "GET", data = null) {
	return new Promise((resolve) => {
		const url = new URL(`/api/${endpoint}`, TN_URL);
		const body = data ? JSON.stringify(data) : null;
		const opts = {
			hostname: url.hostname,
			port: url.port,
			path: url.pathname,
			method,
			headers: { "Content-Type": "application/json" },
			timeout: 10000,
		};
		if (body) opts.headers["Content-Length"] = Buffer.byteLength(body);

		const req = http.request(opts, (res) => {
			let chunks = "";
			res.on("data", (c) => (chunks += c));
			res.on("end", () => {
				try {
					resolve(JSON.parse(chunks || "{}"));
				} catch {
					resolve({ raw: chunks });
				}
			});
		});

		req.on("error", () => resolve({ error: "Kernel not reachable" }));
		req.on("timeout", () => {
			req.destroy();
			resolve({ error: "Timeout" });
		});

		if (body) req.write(body);
		req.end();
	});
}

function showHelp() {
	console.log(`
╔══════════════════════════════════════════════╗
║       TormentNexus CLI — tn v1.0.0-b1      ║
╠══════════════════════════════════════════════╣
║  tn search <query>     Search L2 memory     ║
║  tn store <text>       Store a memory       ║
║  tn status             System health check   ║
║  tn tools              List all MCP tools    ║
║  tn tool <name>        Describe a tool      ║
║  tn install            Setup TN for clients  ║
║  tn help               This message          ║
╚══════════════════════════════════════════════╝
`);
}

async function main() {
	if (!cmd || cmd === "help") return showHelp();

	switch (cmd) {
		case "status": {
			const s = await request("runtime/status");
			if (s.error) return console.log(`❌ ${s.error}`);
			const d = s.data;
			console.log(`\n✅ TormentNexus v${d.version} | Uptime: ${d.uptimeSec}s`);
			console.log(
				`   🛠️  ${d.cli?.toolCount || 0} tools | 🧠 ${d.memory?.l2?.count || 0} memories`,
			);
			console.log(`   🔗 ${TN_URL}\n`);
			break;
		}
		case "search": {
			const query = argv.slice(1).join(" ");
			const r = await request("memory/search", "POST", { query, limit: 10 });
			const results = r.data || r.results || [];
			if (!results.length)
				return console.log(`No memories found for: "${query}"`);
			results.forEach((m, i) =>
				console.log(
					`[${i + 1}] ${(m.content || JSON.stringify(m)).slice(0, 120)}`,
				),
			);
			break;
		}
		case "store": {
			const content = argv.slice(1).join(" ");
			const r = await request("memory/store", "POST", {
				content,
				tags: [],
				category: "cli",
			});
			console.log(`🧠 Stored: ${r.id || "OK"}`);
			break;
		}
		case "tools": {
			const r = await request("runtime/status");
			const tools = r.data?.cli?.tools || [];
			console.log(`\n🛠️  ${tools.length} MCP tools available:\n`);
			tools.forEach((t) => console.log(`   ${t}`));
			break;
		}
		default:
			console.log(`Unknown command: ${cmd}`);
			showHelp();
	}
}

main().catch(console.error);
