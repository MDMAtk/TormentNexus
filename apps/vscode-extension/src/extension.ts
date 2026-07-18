import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
    console.log('TormentNexus VS Code Extension is now active!');

    // Register Webview Provider for the Sidebar
    const provider = new TormentNexusSidebarProvider(context.extensionUri);
    context.subscriptions.push(
        vscode.window.registerWebviewViewProvider('tormentnexus.sidebar', provider)
    );

    // Register Command to Search Memory
    let disposable = vscode.commands.registerCommand('tormentnexus.searchMemory', async () => {
        const query = await vscode.window.showInputBox({ prompt: 'Search TormentNexus Memory' });
        if (query) {
            vscode.window.showInformationMessage(`TormentNexus: Searching for "${query}"...`);
            // Here we'd eventually hit the TN Kernel API directly via http or shell execution
        }
    });
    context.subscriptions.push(disposable);
}

class TormentNexusSidebarProvider implements vscode.WebviewViewProvider {
    constructor(private readonly _extensionUri: vscode.Uri) {}

    public resolveWebviewView(
        webviewView: vscode.WebviewView,
        context: vscode.WebviewViewResolveContext,
        _token: vscode.CancellationToken,
    ) {
        webviewView.webview.options = {
            enableScripts: true,
        };

        webviewView.webview.html = this._getHtmlForWebview();

        // Handle messages from the webview
        webviewView.webview.onDidReceiveMessage(
            async message => {
                switch (message.command) {
                    case 'fetchStatus':
                        const status = await this._checkKernelStatus();
                        webviewView.webview.postMessage({ command: 'updateStatus', status });
                        break;
                    case 'alert':
                        vscode.window.showErrorMessage(message.text);
                        break;
                }
            }
        );
    }

    private async _checkKernelStatus(): Promise<string> {
        try {
            // First check the Go native port 7778
            const response = await fetch('http://127.0.0.1:7778/health');
            if (response.ok) {
                return 'Kernel connected (Port 7778)';
            }
        } catch (e) {
            // Fallback to Next.js or Node wrapper
        }

        try {
            const response = await fetch('http://127.0.0.1:4300/health');
            if (response.ok) {
                return 'Kernel connected (Port 4300)';
            }
        } catch (e) {}

        return 'TN Kernel is unreachable';
    }

    private _getHtmlForWebview() {
        return `<!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>TormentNexus</title>
            <style>
                body {
                    font-family: var(--vscode-font-family);
                    padding: 10px;
                    color: var(--vscode-editor-foreground);
                }
                h1 { font-size: 1.2em; border-bottom: 1px solid var(--vscode-widget-border); padding-bottom: 5px; }
                h2 { font-size: 1em; margin-top: 15px; }
                button {
                    background: var(--vscode-button-background);
                    color: var(--vscode-button-foreground);
                    border: none;
                    padding: 6px 10px;
                    margin-top: 10px;
                    cursor: pointer;
                    width: 100%;
                }
                button:hover {
                    background: var(--vscode-button-hoverBackground);
                }
                .status-box {
                    background: var(--vscode-editor-inactiveSelectionBackground);
                    padding: 8px;
                    border-radius: 4px;
                    margin-top: 10px;
                    font-family: monospace;
                    font-size: 0.9em;
                }
            </style>
        </head>
        <body>
            <h1>⚡ TormentNexus</h1>

            <div class="status-box" id="statusBox">
                Checking Kernel status...
            </div>

            <h2>L1/L2 Memory</h2>
            <button onclick="searchMemory()">Search Memory</button>

            <h2>MCP Catalog</h2>
            <button onclick="openCatalog()">Browse Installed Tools</button>

            <script>
                const vscode = acquireVsCodeApi();

                // On load, ask extension to fetch status
                window.addEventListener('load', () => {
                    vscode.postMessage({ command: 'fetchStatus' });
                });

                window.addEventListener('message', event => {
                    const message = event.data;
                    switch (message.command) {
                        case 'updateStatus':
                            document.getElementById('statusBox').innerText = message.status;
                            break;
                    }
                });

                function searchMemory() {
                    // Trigger vscode command
                    // vscode.postMessage({ command: 'alert', text: 'Not implemented' });
                    // To truly execute a command from webview:
                    // Currently, standard extensions map commands through postMessage logic
                    vscode.postMessage({ command: 'alert', text: 'Memory search triggered.' });
                }

                function openCatalog() {
                    vscode.postMessage({ command: 'alert', text: 'Catalog opening...' });
                }
            </script>
        </body>
        </html>`;
    }
}

export function deactivate() {}
