// TormentNexus VS Code Extension
// Activates on startup and provides MCP server + custom commands

const vscode = require("vscode");
const http = require("http");
const TN_API = "http://127.0.0.1:7778";

function tnPost(path, body) {
  return new Promise((resolve) => {
    const data = JSON.stringify(body);
    const req = http.request(
      `${TN_API}${path}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json", "Content-Length": data.length },
        timeout: 3000,
      },
      () => resolve()
    );
    req.on("error", () => resolve());
    req.write(data);
    req.end();
  });
}

function tnSearch(query) {
  return new Promise((resolve) => {
    const q = encodeURIComponent(query.slice(0, 200));
    http.get(`${TN_API}/api/memory/search?q=${q}`, { timeout: 3000 }, (res) => {
      let data = "";
      res.on("data", (c) => (data += c));
      res.on("end", () => {
        try {
          const json = JSON.parse(data);
          resolve(json.data || []);
        } catch {
          resolve([]);
        }
      });
    }).on("error", () => resolve([]));
  });
}

function activate(context) {
  console.log("TormentNexus extension activating...");

  // Log startup to TN
  tnPost("/api/memory/add", {
    content: JSON.stringify({
      content: "VS Code session started",
      tags: ["system:session", "agent:vscode"],
      category: "session",
    }),
  });

  // Register commands
  context.subscriptions.push(
    vscode.commands.registerCommand("tormentnexus.tnStore", async () => {
      const content = await vscode.window.showInputBox({
        prompt: "What do you want to store in TormentNexus?",
        placeHolder: "e.g., project uses React 19 with Vite",
      });
      if (!content) return;
      await tnPost("/api/memory/add", {
        content: JSON.stringify({ content, tags: ["agent:vscode"], category: "memory" }),
      });
      vscode.window.showInformationMessage("✅ Stored in TormentNexus");
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand("tormentnexus.tnSearch", async () => {
      const query = await vscode.window.showInputBox({
        prompt: "What do you want to find in TormentNexus?",
        placeHolder: "e.g., React patterns",
      });
      if (!query) return;
      const results = await tnSearch(query);
      if (results.length === 0) {
        vscode.window.showInformationMessage("No results found");
        return;
      }
      const text = results.map((r) => `• ${r.text || r.content || ""}`).join("\n\n");
      const panel = vscode.window.createWebviewPanel("tnSearch", "TormentNexus Results", vscode.ViewColumn.One);
      panel.webview.html = `<!DOCTYPE html><html><body><pre>${text}</pre></body></html>`;
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand("tormentnexus.tnStatus", async () => {
      vscode.window.showInformationMessage("TormentNexus MCP server configured in .vscode/mcp.json");
    })
  );
}

function deactivate() {
  tnPost("/api/memory/add", {
    content: JSON.stringify({
      content: "VS Code session ended",
      tags: ["system:session_end", "agent:vscode"],
      category: "session",
    }),
  });
}

module.exports = { activate, deactivate };
