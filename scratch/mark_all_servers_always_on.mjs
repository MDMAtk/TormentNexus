#!/usr/bin/env node
/**
 * Mark all configured MCP servers as always_on=true via the Go sidecar update API.
 * Uses uuid for lookup since the update endpoint is uuid-keyed.
 */

const GO_BASE = "http://localhost:4300";

async function main() {
	console.log("[mark-always-on] Fetching configured servers...");
	const listResp = await fetch(`${GO_BASE}/api/mcp/servers/configured`);
	const listJson = await listResp.json();
	const data = listJson.data ?? listJson;
	const servers = Array.isArray(data) ? data : (data.servers ?? []);
	console.log(`[mark-always-on] Found ${servers.length} configured servers`);

	let success = 0;
	let failed = 0;
	let alreadyOn = 0;

	for (const server of servers) {
		const uuid = server.uuid;
		const name = server.name;
		if (!uuid) {
			failed++;
			continue;
		}
		if (server.always_on === true) {
			alreadyOn++;
			continue;
		}
		try {
			const updatePayload = { uuid, always_on: true };
			const resp = await fetch(`${GO_BASE}/api/mcp/servers/update`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(updatePayload),
			});
			const json = await resp.json();
			if (json.success) {
				success++;
				if (success % 50 === 0) {
					process.stdout.write(`\n  Progress: ${success} updated\n`);
				} else {
					process.stdout.write(".");
				}
			} else {
				failed++;
				process.stdout.write(
					`\n✗ ${name}: ${JSON.stringify(json).slice(0, 100)}\n`,
				);
			}
		} catch (err) {
			failed++;
			process.stdout.write(`\n✗ ${name}: ${err.message}\n`);
		}
	}

	console.log(`\n\n[mark-always-on] Done.`);
	console.log(`  Updated: ${success}`);
	console.log(`  Already on: ${alreadyOn}`);
	console.log(`  Failed: ${failed}`);
	console.log(`  Total: ${servers.length}`);
}

main().catch((err) => {
	console.error("[mark-always-on] Fatal:", err);
	process.exit(1);
});
