// @borg/ai — type-safe stub package
// Provides type declarations for consumers while deferring real implementations
// to the TypeScript core runtime.

// ── Type declarations ──────────────────────────────────────────────────────

export interface ChatMessage {
  role: 'system' | 'user' | 'assistant' | 'tool';
  content: string;
  name?: string;
  toolCallId?: string;
  toolCalls?: any[];
}

export interface ModelSelectionRequest {
  taskComplexity?: 'low' | 'medium' | 'high';
  taskType?: string;
  routingTaskType?: string;
  provider?: string;
}

export interface SelectedModel {
  provider: string;
  modelId: string;
  reason?: string;
}

export interface QuotaConfig {
  providers?: Record<string, any>;
  dailyBudgetUsd?: number;
  monthlyBudgetUsd?: number;
  providerLimits?: Record<string, any>;
}

export interface IAgent {
  id: string;
  name: string;
  role: string;
  start(): Promise<void>;
  stop(): Promise<void>;
}

// ── Class stubs ────────────────────────────────────────────────────────────

export class ModelSelector {
  async selectModel(_req?: Partial<ModelSelectionRequest & Record<string, any>>): Promise<SelectedModel> {
    return { provider: 'stub', modelId: 'stub-model', reason: '@borg/ai stub' };
  }
}

export class LLMService {
  modelSelector: ModelSelector = new ModelSelector();
  async generateText(
    _provider: string,
    _modelId: string,
    _systemPrompt: string,
    _userPrompt: string,
    _opts?: any,
  ): Promise<{ content: string; usage?: { inputTokens: number; outputTokens: number } }> {
    return { content: '[stub response from @borg/ai]' };
  }
}

export class QuotaService {
  protected configState: QuotaConfig = { providers: {}, dailyBudgetUsd: 5, monthlyBudgetUsd: 100, providerLimits: {} };
  setConfig(config: QuotaConfig) { Object.assign(this.configState, config); }
  getConfig(): QuotaConfig { return this.configState; }
  markAuthRevoked(_provider: string): void {}
}

export const DEFAULT_OPENROUTER_FREE_MODEL = 'google/gemma-3-27b-it:free';

// ── Additional stubs that @borg/core references ────────────────────────────

export class SearchService {
  async search(_query: string, _opts?: any): Promise<any[]> { return []; }
}

export interface SearchResult {
  title: string;
  url: string;
  snippet: string;
  score?: number;
}

export class InputTools {
  async process(_input: string): Promise<any> { return {}; }
}

export class SystemStatusTool {
  async getStatus(): Promise<any> { return { status: 'stub' }; }
}

export class ProcessRegistry {
  list(): any[] { return []; }
  get(_name: string): any { return null; }
}

export class ChainExecutor {
  async executeChain(_req: any): Promise<any> { return { result: '[ChainExecutor stub]' }; }
}

export class ChainRequest {
  tools?: string[];
  input?: string;
  [key: string]: any;
}

export class TerminalService {
  async execute(_cmd: string): Promise<any> { return { output: '', exitCode: 0 }; }
}

export class BrowserTool {
  async browse(_url: string): Promise<any> { return { content: '[BrowserTool stub]' }; }
}

export class McpServerRegistry {
  list(): any[] { return []; }
  get(_name: string): any { return null; }
}

export class IMCPServer {
  async callTool(_name: string, _args: any): Promise<any> { return null; }
  listTools(): any[] { return []; }
}
