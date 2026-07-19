"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = __importStar(require("vscode"));
const axios_1 = __importDefault(require("axios"));
let serverUrl = "http://localhost:7778";
let isConnected = false;
function activate(context) {
    console.log("TormentNexus extension activated");
    // Load configuration
    const config = vscode.workspace.getConfiguration("tormentnexus");
    serverUrl = config.get("serverUrl", "http://localhost:7778");
    // Register commands
    context.subscriptions.push(vscode.commands.registerCommand("tormentnexus.connect", connect), vscode.commands.registerCommand("tormentnexus.disconnect", disconnect), vscode.commands.registerCommand("tormentnexus.searchTools", searchTools), vscode.commands.registerCommand("tormentnexus.addMemory", addMemory), vscode.commands.registerCommand("tormentnexus.searchMemory", searchMemory), vscode.commands.registerCommand("tormentnexus.openDashboard", openDashboard), vscode.commands.registerCommand("tormentnexus.refresh", refresh));
    // Register tree data providers
    const memoryProvider = new MemoryTreeProvider();
    const toolsProvider = new ToolsTreeProvider();
    const statusProvider = new StatusTreeProvider();
    context.subscriptions.push(vscode.window.registerTreeDataProvider("tormentnexus.memory", memoryProvider), vscode.window.registerTreeDataProvider("tormentnexus.tools", toolsProvider), vscode.window.registerTreeDataProvider("tormentnexus.status", statusProvider));
    // Auto-connect if configured
    if (config.get("autoConnect", true)) {
        connect();
    }
}
async function connect() {
    try {
        const response = await axios_1.default.get(`${serverUrl}/health`);
        if (response.data.ok) {
            isConnected = true;
            vscode.window.showInformationMessage("TormentNexus connected!");
            vscode.commands.executeCommand("tormentnexus.refresh");
        }
    }
    catch (error) {
        isConnected = false;
        vscode.window.showErrorMessage("Failed to connect to TormentNexus");
    }
}
function disconnect() {
    isConnected = false;
    vscode.window.showInformationMessage("TormentNexus disconnected");
    vscode.commands.executeCommand("tormentnexus.refresh");
}
async function searchTools() {
    const query = await vscode.window.showInputBox({
        prompt: "Search MCP tools",
        placeHolder: "e.g., postgres, filesystem, browser",
    });
    if (!query)
        return;
    try {
        const response = await axios_1.default.get(`${serverUrl}/api/backlog/search`, {
            params: { q: query, limit: 20 },
        });
        const tools = response.data.results || [];
        if (tools.length === 0) {
            vscode.window.showInformationMessage("No tools found");
            return;
        }
        const items = tools.map((t) => ({
            label: t.name || t.title,
            description: t.description?.substring(0, 100),
            detail: t.url,
        }));
        const selected = await vscode.window.showQuickPick(items, {
            placeHolder: "Select a tool to view details",
        });
        if (selected) {
            const url = selected.detail;
            if (url) {
                vscode.env.openExternal(vscode.Uri.parse(url));
            }
        }
    }
    catch (error) {
        vscode.window.showErrorMessage("Failed to search tools");
    }
}
async function addMemory() {
    const content = await vscode.window.showInputBox({
        prompt: "Add a memory",
        placeHolder: "Enter the memory content...",
    });
    if (!content)
        return;
    try {
        await axios_1.default.post(`${serverUrl}/api/memory/add`, {
            content,
            tags: ["vscode"],
            source: "vscode-extension",
        });
        vscode.window.showInformationMessage("Memory added!");
    }
    catch (error) {
        vscode.window.showErrorMessage("Failed to add memory");
    }
}
async function searchMemory() {
    const query = await vscode.window.showInputBox({
        prompt: "Search memory",
        placeHolder: "Enter search query...",
    });
    if (!query)
        return;
    try {
        const response = await axios_1.default.get(`${serverUrl}/api/memory/search`, {
            params: { q: query, limit: 10 },
        });
        const results = response.data.results || [];
        if (results.length === 0) {
            vscode.window.showInformationMessage("No memories found");
            return;
        }
        const items = results.map((r) => ({
            label: r.content?.substring(0, 50),
            description: r.tags?.join(", "),
            detail: r.content,
        }));
        await vscode.window.showQuickPick(items, {
            placeHolder: "Search results",
        });
    }
    catch (error) {
        vscode.window.showErrorMessage("Failed to search memory");
    }
}
function openDashboard() {
    vscode.env.openExternal(vscode.Uri.parse(serverUrl));
}
function refresh() {
    // Refresh all tree views
    vscode.commands.executeCommand("tormentnexus.memory.refresh");
    vscode.commands.executeCommand("tormentnexus.tools.refresh");
    vscode.commands.executeCommand("tormentnexus.status.refresh");
}
// Tree Data Providers
class MemoryTreeProvider {
    constructor() {
        this._onDidChangeTreeData = new vscode.EventEmitter();
        this.onDidChangeTreeData = this._onDidChangeTreeData.event;
    }
    refresh() {
        this._onDidChangeTreeData.fire(undefined);
    }
    getTreeItem(element) {
        return element;
    }
    async getChildren(element) {
        if (!isConnected) {
            return [
                new vscode.TreeItem("Not connected", vscode.TreeItemCollapsibleState.None),
            ];
        }
        try {
            const response = await axios_1.default.get(`${serverUrl}/api/memory/status`);
            const status = response.data;
            return [
                new vscode.TreeItem(`L1: ${status.l1Count || 0} entries`, vscode.TreeItemCollapsibleState.None),
                new vscode.TreeItem(`L2: ${status.l2Count || 0} entries`, vscode.TreeItemCollapsibleState.None),
                new vscode.TreeItem(`L3: ${status.l3Count || 0} entries`, vscode.TreeItemCollapsibleState.None),
                new vscode.TreeItem(`L4: ${status.l4Count || 0} entries`, vscode.TreeItemCollapsibleState.None),
            ];
        }
        catch {
            return [
                new vscode.TreeItem("Error loading memory", vscode.TreeItemCollapsibleState.None),
            ];
        }
    }
}
class ToolsTreeProvider {
    constructor() {
        this._onDidChangeTreeData = new vscode.EventEmitter();
        this.onDidChangeTreeData = this._onDidChangeTreeData.event;
    }
    refresh() {
        this._onDidChangeTreeData.fire(undefined);
    }
    getTreeItem(element) {
        return element;
    }
    async getChildren(element) {
        if (!isConnected) {
            return [
                new vscode.TreeItem("Not connected", vscode.TreeItemCollapsibleState.None),
            ];
        }
        try {
            const response = await axios_1.default.get(`${serverUrl}/api/backlog/stats`);
            const stats = response.data;
            return [
                new vscode.TreeItem(`Total: ${stats.total || 0} tools`, vscode.TreeItemCollapsibleState.None),
                new vscode.TreeItem(`Enriched: ${stats.enriched || 0}`, vscode.TreeItemCollapsibleState.None),
                new vscode.TreeItem(`Categories: ${Object.keys(stats.byCategory || {}).length}`, vscode.TreeItemCollapsibleState.None),
            ];
        }
        catch {
            return [
                new vscode.TreeItem("Error loading tools", vscode.TreeItemCollapsibleState.None),
            ];
        }
    }
}
class StatusTreeProvider {
    constructor() {
        this._onDidChangeTreeData = new vscode.EventEmitter();
        this.onDidChangeTreeData = this._onDidChangeTreeData.event;
    }
    refresh() {
        this._onDidChangeTreeData.fire(undefined);
    }
    getTreeItem(element) {
        return element;
    }
    async getChildren(element) {
        if (!isConnected) {
            return [
                new vscode.TreeItem("⚪ Disconnected", vscode.TreeItemCollapsibleState.None),
            ];
        }
        try {
            const response = await axios_1.default.get(`${serverUrl}/health`);
            const health = response.data;
            return [
                new vscode.TreeItem("🟢 Connected", vscode.TreeItemCollapsibleState.None),
                new vscode.TreeItem(`Version: ${health.version}`, vscode.TreeItemCollapsibleState.None),
                new vscode.TreeItem(`Uptime: ${Math.floor(health.uptimeSec / 60)}m`, vscode.TreeItemCollapsibleState.None),
            ];
        }
        catch {
            return [
                new vscode.TreeItem("🔴 Error", vscode.TreeItemCollapsibleState.None),
            ];
        }
    }
}
function deactivate() {
    console.log("TormentNexus extension deactivated");
}
//# sourceMappingURL=extension.js.map