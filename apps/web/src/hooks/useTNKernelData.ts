/**
 * @file useTNKernelData.ts
 * @module apps/web/src/hooks
 *
 * WHAT: React hooks that provide dashboard data from the TN Kernel API
 * when the TypeScript tRPC core is unreachable.
 *
 * WHY: Full Assimilation — the dashboard must display real data even when
 * the TS core is down. These hooks query the TN Kernel REST endpoints
 * and return data in the same shape as the tRPC queries, so the
 * DashboardHomeView component can render without changes.
 *
 * ADDED: v1.0.0-alpha.52
 */

'use client';

import { useEffect, useState, useCallback, useRef } from 'react';

const TN_KERNEL_BASE = process.env.NEXT_PUBLIC_TN_KERNEL_URL || '';
const GO_PROXY_PREFIX = '/api/go';

/**
 * Build a URL for the TN Kernel.
 *
 * When a direct kernel URL is configured (NEXT_PUBLIC_TN_KERNEL_URL),
 * paths are used as-is.  When going through the Next.js catch-all proxy
 * at /api/go/[...path], the proxy passes the path verbatim to the
 * kernel, so /api/go/health → kernel /health,
 * /api/go/api/mcp/status → kernel /api/mcp/status, etc.
 */
function kernelUrl(path: string): string {
  if (TN_KERNEL_BASE) {
    return `${TN_KERNEL_BASE}${path}`;
  }
  return `${GO_PROXY_PREFIX}${path}`;
}

/** Fetch JSON from the TN Kernel with timeout. */
async function fetchKernel<T>(path: string, init?: RequestInit): Promise<T | null> {
  try {
    const url = kernelUrl(path);
    const res = await fetch(url, {
      headers: { accept: 'application/json' },
      signal: AbortSignal.timeout(5000),
      ...init,
    });
    if (!res.ok) return null;
    const json = await res.json();
    return json.success ? (json.data as T) : null;
  } catch {
    return null;
  }
}

// ────────────────────────────────────────────────────────────────────────────
// Types — mirror the DashboardHomeView interfaces so we can slot data in
// ────────────────────────────────────────────────────────────────────────────

export interface GoMCPStatus {
  initialized: boolean;
  serverCount: number;
  toolCount: number;
  connectedCount: number;
}

export interface GoServerSummary {
  name: string;
  status: string;
  toolCount: number;
  config?: { command: string; args: string[]; env: string[] };
}

export interface GoStartupStatus {
  status: string;
  ready: boolean;
  uptime: number;
  summary?: string;
  blockingReasons?: Array<{ code: string; detail: string }>;
  checks: {
    mcpAggregator: {
      ready: boolean;
      liveReady?: boolean;
      residentReady?: boolean;
      serverCount: number;
      connectedCount?: number;
      residentConnectedCount?: number;
      warmingServerCount?: number;
      failedWarmupServerCount?: number;
      initialization: {
        inProgress: boolean;
        initialized: boolean;
        connectedClientCount: number;
        configuredServerCount: number;
      } | null;
      persistedServerCount: number;
      persistedToolCount: number;
      configuredServerCount?: number;
      advertisedServerCount?: number;
      advertisedToolCount?: number;
      advertisedAlwaysOnServerCount?: number;
      advertisedAlwaysOnToolCount?: number;
      inventoryReady: boolean;
      inventorySource?: 'database' | 'config' | 'empty';
      warmupInProgress?: boolean;
    };
    configSync: {
      ready: boolean;
      status: {
        inProgress: boolean;
        lastServerCount: number;
        lastToolCount: number;
      } | null;
    };
    memory: {
      ready: boolean;
      initialized: boolean;
      agentMemory: boolean;
      claudeMem?: {
        ready?: boolean;
        enabled?: boolean;
        storeExists?: boolean;
        totalEntries?: number;
        sectionCount?: number;
        defaultSectionCount?: number;
        presentDefaultSectionCount?: number;
        missingSections?: string[];
        lastUpdatedAt?: string | null;
      };
      tormentnexus?: {
        ready?: boolean;
        enabled?: boolean;
        storeExists?: boolean;
        totalEntries?: number;
        sectionCount?: number;
        defaultSectionCount?: number;
        presentDefaultSectionCount?: number;
        missingSections?: string[];
        lastUpdatedAt?: string | null;
      };
    };
    browser: { ready: boolean; active: boolean; pageCount: number };
    sessionSupervisor: {
      ready: boolean;
      sessionCount: number;
      restore: {
        restoredSessionCount: number;
        autoResumeCount: number;
      } | null;
    };
    extensionBridge: {
      ready: boolean;
      acceptingConnections?: boolean;
      clientCount: number;
      hasConnectedClients?: boolean;
    };
    executionEnvironment: {
      ready: boolean;
      preferredShellId?: string | null;
      preferredShellLabel?: string | null;
      shellCount: number;
      verifiedShellCount: number;
      toolCount: number;
      verifiedToolCount: number;
      harnessCount: number;
      verifiedHarnessCount: number;
      supportsPowerShell: boolean;
      supportsPosixShell: boolean;
      notes?: string[];
    };
  };
}

export interface GoProviderSummary {
  provider: string;
  name: string;
  configured: boolean;
  authenticated?: boolean;
  tier: string;
  limit: number | null;
  used: number;
  remaining: number | null;
  resetDate?: string | null;
  availability?: string;
  lastError?: string | null;
}

export interface GoFallbackSummary {
  priority: number;
  provider: string;
  model?: string;
  reason: string;
}

export interface GoSessionSummary {
  id: string;
  name: string;
  cliType: string;
  workingDirectory: string;
  status: 'created' | 'starting' | 'running' | 'stopping' | 'stopped' | 'restarting' | 'error';
  restartCount: number;
  maxRestartAttempts: number;
  lastActivityAt: number;
  lastError?: string;
  logs: Array<{ timestamp: number; stream: 'stdout' | 'stderr' | 'system'; message: string }>;
}

// ────────────────────────────────────────────────────────────────────────────
// Hook: useTNKernelDashboard
// ────────────────────────────────────────────────────────────────────────────

export interface TNKernelDashboardData {
  mcpStatus: GoMCPStatus | null;
  startupStatus: GoStartupStatus | null;
  servers: GoServerSummary[];
  providers: GoProviderSummary[];
  fallbackChain: GoFallbackSummary[];
  sessions: GoSessionSummary[];
  goVersion: string | null;
  connected: boolean;
  lastFetchedAt: number | null;
}

const EMPTY_DASHBOARD: TNKernelDashboardData = {
  mcpStatus: null,
  startupStatus: null,
  servers: [],
  providers: [],
  fallbackChain: [],
  sessions: [],
  goVersion: null,
  connected: false,
  lastFetchedAt: null,
};

/**
 * useTNKernelDashboard polls the TN Kernel REST endpoints and returns
 * dashboard data. Intended as a fallback data source when tRPC queries fail.
 */
export function useTNKernelDashboard(pollIntervalMs = 5000): TNKernelDashboardData {
  const [data, setData] = useState<TNKernelDashboardData>(EMPTY_DASHBOARD);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const poll = useCallback(async () => {
    const [
      mcpStatus,
      startupStatus,
      servers,
      billingStatus,
      sessions,
      healthData,
      healthFallback,
    ] = await Promise.all([
      fetchKernel<GoMCPStatus>('/api/mcp/status'),
      fetchKernel<GoStartupStatus>('/api/startup/status'),
      fetchKernel<GoServerSummary[]>('/api/mcp/servers'),
      fetchKernel<any>('/api/billing/status'),
      fetchKernel<{ count: number; sessions: GoSessionSummary[] }>('/api/native/session/list'),
      fetchKernel<any>('/api/health'),
      fetchKernel<any>('/health'),
    ]);

    // Use /api/health result, fall back to /health (older binaries)
    const effectiveHealth = healthData ?? healthFallback;

    // Extract providers from billing status
    let providers: GoProviderSummary[] = [];
    let fallbackChain: GoFallbackSummary[] = [];
    if (billingStatus) {
      providers = billingStatus.providers ?? billingStatus.providerQuotas ?? [];
      fallbackChain = billingStatus.fallbackChain ?? [];
    }

    const connected = mcpStatus !== null || startupStatus !== null || effectiveHealth !== null;

    let parsedSessions: GoSessionSummary[] = [];
    if (sessions) {
      if (Array.isArray(sessions)) {
        parsedSessions = sessions;
      } else if (Array.isArray((sessions as any)?.sessions)) {
        parsedSessions = (sessions as any).sessions;
      }
    }

    setData({
      mcpStatus,
      startupStatus,
      servers: servers ?? [],
      providers,
      fallbackChain,
      sessions: parsedSessions,
      goVersion: effectiveHealth?.version ?? null,
      connected,
      lastFetchedAt: Date.now(),
    });
  }, []);

  useEffect(() => {
    poll();
    intervalRef.current = setInterval(poll, pollIntervalMs);
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [poll, pollIntervalMs]);

  return data;
}

/**
 * useTNKernelConnectivity returns just the connection status and version.
 * Lightweight — only hits /api/health.
 */
export function useTNKernelConnectivity(pollIntervalMs = 10000): {
  connected: boolean;
  version: string | null;
  uptime: number | null;
} {
  const [state, setState] = useState({ connected: false, version: null as string | null, uptime: null as number | null });

  useEffect(() => {
    let cancelled = false;

    async function check() {
      try {
        // Try /api/health first (newer binaries), fall back to /health (older)
        let health = await fetchKernel<any>('/api/health');
        if (!health) {
          health = await fetchKernel<any>('/health');
        }
        if (!cancelled) {
          setState({
            connected: health !== null,
            version: health?.version ?? null,
            uptime: health?.uptimeSec ?? health?.uptime ?? null,
          });
        }
      } catch {
        if (!cancelled) setState({ connected: false, version: null, uptime: null });
      }
    }

    check();
    const id = setInterval(check, pollIntervalMs);
    return () => { cancelled = true; clearInterval(id); };
  }, [pollIntervalMs]);

  return state;
}
