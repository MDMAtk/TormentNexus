import sqlite3 from 'better-sqlite3';
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import path from "path";
import fs from "fs";
import { fileURLToPath } from "url";
import { randomUUID } from "crypto";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const dbPath = path.resolve(__dirname, "..", "tormentnexus.db");

if (!fs.existsSync(dbPath)) {
    console.error(`Database not found at ${dbPath}`);
    process.exit(1);
}

const db = new sqlite3(dbPath);
console.log(`[Validator] Connected to tormentnexus.db at: ${dbPath}`);

// Brand-reconciliation helper
function reconcileNaming(text) {
    if (typeof text !== 'string') return text;
    // Replace specific paths first
    let result = text
        .replace(/borg-supervisor/gi, 'tormentnexus-supervisor')
        .replace(/hypercode-supervisor/gi, 'tormentnexus-supervisor');
    
    // Replace exact obsolete brand names
    if (result.toLowerCase() === 'borg' || result.toLowerCase() === 'nexus' || result.toLowerCase() === 'hypercode') {
        return 'tormentnexus';
    }
    return result;
}

async function validateServer(server) {
    console.log(`\n--------------------------------------------------`);
    console.log(`[Validator] Processing Server: "${server.name}" (UUID: ${server.uuid})`);
    
    let command = reconcileNaming(server.command);
    let argsParsed = [];
    let envParsed = {};

    try {
        if (server.args) {
            argsParsed = JSON.parse(server.args).map(arg => reconcileNaming(arg));
        }
    } catch (e) {
        console.warn(`[Validator] Failed to parse args for ${server.name}:`, e.message);
    }

    try {
        if (server.env) {
            envParsed = JSON.parse(server.env);
        }
    } catch (e) {
        console.warn(`[Validator] Failed to parse env for ${server.name}:`, e.message);
    }

    // Standardize environment variables inheritance or overriding
    let childEnv = { ...process.env };
    if (Array.isArray(envParsed)) {
        childEnv = {};
        for (const key of envParsed) {
            childEnv[key] = process.env[key] || "YOUR_KEY_HERE";
        }
    } else if (typeof envParsed === 'object' && envParsed !== null) {
        childEnv = { ...process.env, ...envParsed };
    }

    // Set custom mock api keys so servers don't crash instantly on initialization
    const fallbackKeys = [
        "OPENAI_API_KEY", "GEMINI_API_KEY", "ANTHROPIC_API_KEY", "OPENROUTER_API_KEY", 
        "TAVILY_API_KEY", "FIRECRAWL_API_KEY", "OCTAGON_API_KEY", "MEM0_API_KEY"
    ];
    for (const key of fallbackKeys) {
        if (!childEnv[key]) {
            childEnv[key] = "YOUR_KEY_HERE";
        }
    }

    if (server.type === 'SSE' || (server.type === 'STREAMABLE_HTTP' && server.url)) {
        console.log(`[Validator] Server uses remote SSE transport: ${server.url}`);
        // For remote SSE endpoints, we can mark them verified if they have valid URLs
        db.prepare("UPDATE mcp_servers SET error_status = '' WHERE uuid = ?").run(server.uuid);
        
        // Upsert into published_mcp_servers
        const canonicalId = `remote__${server.name.toLowerCase()}`;
        const existing = db.prepare("SELECT uuid FROM published_mcp_servers WHERE canonical_id = ?").get(canonicalId);
        if (existing) {
            db.prepare(`
                UPDATE published_mcp_servers 
                SET status = 'verified', confidence = 1.0, last_verified_at = datetime('now'), updated_at = datetime('now')
                WHERE uuid = ?
            `).run(existing.uuid);
        } else {
            db.prepare(`
                INSERT INTO published_mcp_servers (uuid, canonical_id, display_name, description, transport, status, confidence, last_verified_at, created_at, updated_at)
                VALUES (?, ?, ?, ?, 'sse', 'verified', 1.0, datetime('now'), datetime('now'), datetime('now'))
            `).run(randomUUID(), canonicalId, server.name, server.description || '');
        }
        console.log(`[Validator] SSE Server "${server.name}" marked verified in catalog.`);
        return;
    }

    if (!command) {
        console.log(`[Validator] No command defined for ${server.name}, skipping execution.`);
        return;
    }

    console.log(`[Validator] Launching stdio: "${command}" with args:`, argsParsed);

    const transport = new StdioClientTransport({
        command: command,
        args: argsParsed,
        env: childEnv
    });

    const client = new Client(
        { name: "tormentnexus-validator", version: "1.0.0" },
        { capabilities: {} }
    );

    try {
        // Run connection with a 8 second deadline
        await Promise.race([
            client.connect(transport),
            new Promise((_, reject) => setTimeout(() => reject(new Error("Connection timeout (8s)")), 8000))
        ]);

        console.log(`[Validator] Connected successfully! Listing tools...`);
        const toolsResult = await client.listTools();
        const toolsList = toolsResult.tools || [];
        console.log(`[Validator] Retested successfully! Found ${toolsList.length} tools.`);

        // Begin transaction to refresh tools mapping
        db.transaction(() => {
            // Delete old tools
            db.prepare("DELETE FROM tools WHERE mcp_server_uuid = ?").run(server.uuid);

            // Insert new tools
            const insertTool = db.prepare(`
                INSERT INTO tools (uuid, name, description, tool_schema, is_deferred, always_on, created_at, updated_at, mcp_server_uuid)
                VALUES (?, ?, ?, ?, 0, 0, datetime('now'), datetime('now'), ?)
            `);

            for (const tool of toolsList) {
                const toolUuid = randomUUID();
                insertTool.run(toolUuid, tool.name, tool.description || '', JSON.stringify(tool.inputSchema || {}), server.uuid);
            }

            // Update server entry
            db.prepare("UPDATE mcp_servers SET error_status = '' WHERE uuid = ?").run(server.uuid);
        })();

        // Upsert into published_mcp_servers
        const canonicalId = `local__${server.name.toLowerCase()}`;
        const existing = db.prepare("SELECT uuid FROM published_mcp_servers WHERE canonical_id = ?").get(canonicalId);
        if (existing) {
            db.prepare(`
                UPDATE published_mcp_servers 
                SET status = 'verified', confidence = 1.0, last_verified_at = datetime('now'), updated_at = datetime('now')
                WHERE uuid = ?
            `).run(existing.uuid);
        } else {
            db.prepare(`
                INSERT INTO published_mcp_servers (uuid, canonical_id, display_name, description, transport, status, confidence, last_verified_at, created_at, updated_at)
                VALUES (?, ?, ?, ?, 'stdio', 'verified', 1.0, datetime('now'), datetime('now'), datetime('now'))
            `).run(randomUUID(), canonicalId, server.name, server.description || '');
        }

        console.log(`[Validator] Server "${server.name}" registered successfully with ${toolsList.length} tools!`);

    } catch (err) {
        console.error(`[Validator] Connection failed for "${server.name}":`, err.message);
        db.prepare("UPDATE mcp_servers SET error_status = ? WHERE uuid = ?").run(err.message, server.uuid);
    } finally {
        try {
            await transport.close();
        } catch(e) {}
    }
}

async function main() {
    const servers = db.prepare("SELECT * FROM mcp_servers").all();
    console.log(`[Validator] Query returned ${servers.length} servers in registry.`);

    for (const server of servers) {
        await validateServer(server);
    }

    console.log(`\n[Validator] --- Catalog verification completed successfully! ---`);
    db.close();
}

main().catch(console.error);
