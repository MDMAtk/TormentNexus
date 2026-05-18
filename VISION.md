# Vision: AI Hypervisor (Nexus & HyperCode)

## North Star: The AI Hypervisor
**Nexus** is the underlying coordination kernel and "AI Hypervisor." It acts as the Operating System for AI models, abstracting provider complexity, managing context windows via biological tiered memory, and providing a deterministic execution sandbox.

**HyperCode** is the flagship developer-facing product (CLI and Dashboard) providing an autonomous coding runtime powered by the Nexus Kernel.

## Core Philosophical Pillars
1. **Models as Compute**: Models are ephemeral resources. Nexus manages their allocation, fallback routing, and token budgets.
2. **Tools as Drivers**: MCP servers are "device drivers" for the AI OS. Nexus provides a unified interface for tool discovery, ranking, and progressive disclosure.
3. **Biological Memory**: Intelligence is only as good as its relevance. Nexus utilizes L1 (Active), L2 (Long-Term), and L3 (Cold Archive) tiers with "Heat-based" mechanics (relevance increases heat, time causes decay).
4. **Autonomous Immune System**: The system should heal itself. Every failure is an opportunity for diagnosis, remediation, and verification, with results persisted for collective learning.

## Architectural Layers
- **Nexus Runtime (Go Kernel)**: The authoritative execution kernel (State, Memory, LLM routing, MCP sync). Standardized on Port 4300.
- **Nexus Memory (L1/L2/L3)**: Active memory substrate with SQLite-vec for semantic search and heat-score lifecycle management.
- **Nexus Router**: Progressive tool disclosure and budget-aware provider waterfall.
- **HyperCode Control Plane (TS)**: Next.js dashboard (Port 3000) and tRPC middleware (Port 4100) for observation and high-level agent mission coordination.

## Implementation Milestones

### Phase 4: Deep Orchestration (v1.0.0-alpha.60)
Hardened the multi-agent coordination layer. Implemented the **PairOrchestrator** with a strict `Planner -> Checker -> Implementer -> Critic` state machine and weighted **Consensus Engine**. Integrated **Quota Management** for budget-aware model switching.

### Phase 5: The Immune System (v1.0.0-alpha.61)
Upgraded the **Autonomous Healer** to a full multi-turn loop. Nexus now performs its own `Diagnose -> Fix -> Verify -> Retry` cycles using the native `CodeExecutor` and persists every attempt into the **L2 Vault** for fleet-wide shared intelligence.

### Phase 6: Native Integration & Protocol (Target v1.1.0)
The next evolution focuses on the transition from Electron to a Wails-native runtime and the introduction of the `hypercode://` protocol for seamless browser-to-kernel attachment.
