"use client";

import { useState, useEffect, useCallback } from "react";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
	CardDescription,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { CatalogBrowser } from "@/components/CatalogBrowser";
import {
	Activity,
	Brain,
	Shield,
	Terminal,
	Box,
	RefreshCw,
	CheckCircle2,
	AlertTriangle,
	Search,
	Database,
	Server,
	Globe,
	Users,
	Clock,
	Plus,
	Play,
	ChevronRight,
	FileCode,
	GitBranch,
	Eye,
	Network,
	Code,
	Wrench,
	Gauge,
	HardDrive,
	Wifi,
	MonitorSmartphone,
	KeyRound,
	Fingerprint,
	Building2,
	CreditCard,
	Store,
	Send,
	Copy,
	Check,
	Lock,
} from "lucide-react";

interface SystemHealth {
	service: string;
	version: string;
	uptimeSec: number;
	baseUrl: string;
}

interface MemoryEntry {
	id: string;
	content: string;
	category: string;
	tags: string[];
	importance: number;
	heatScore: number;
	createdAt: string;
}

export function UnifiedDashboard() {
	const [activeTab, setActiveTab] = useState("overview");
	const [health, setHealth] = useState<SystemHealth | null>(null);
	const [memorySearch, setMemorySearch] = useState("");
	const [memoryResults, setMemoryResults] = useState<MemoryEntry[]>([]);
	const [memoryStore, setMemoryStore] = useState("");
	const [memoryStoreTags, setMemoryStoreTags] = useState("");
	const [tools, setTools] = useState<string[]>([]);
	const [toolSearch, setToolSearch] = useState("");
	const [loading, setLoading] = useState<Record<string, boolean>>({});
	const [copied, setCopied] = useState("");

	const setLoadingState = useCallback((key: string, val: boolean) => {
		setLoading((prev) => ({ ...prev, [key]: val }));
	}, []);

	const fetchHealth = useCallback(async () => {
		setLoadingState("health", true);
		try {
			const res = await fetch("/api/system");
			const data = await res.json();
			setHealth(data);
			if (data.tools) setTools(data.tools.map((t: any) => t.name || t));
		} catch {
			/* offline */
		}
		setLoadingState("health", false);
	}, [setLoadingState]);

	const searchMemory = useCallback(async () => {
		if (!memorySearch.trim()) return;
		setLoadingState("memorySearch", true);
		try {
			const res = await fetch("/api/memory/search", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ query: memorySearch, limit: 20 }),
			});
			const data = await res.json();
			setMemoryResults(data.results || data.memories || []);
		} catch {
			setMemoryResults([]);
		}
		setLoadingState("memorySearch", false);
	}, [memorySearch, setLoadingState]);

	const storeMemory = useCallback(async () => {
		if (!memoryStore.trim()) return;
		setLoadingState("memoryStore", true);
		try {
			await fetch("/api/memory/remember", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					content: memoryStore,
					tags: memoryStoreTags
						.split(",")
						.map((t) => t.trim())
						.filter(Boolean),
				}),
			});
			setMemoryStore("");
			setMemoryStoreTags("");
		} catch {
			/* error */
		}
		setLoadingState("memoryStore", false);
	}, [memoryStore, memoryStoreTags, setLoadingState]);

	useEffect(() => {
		fetchHealth();
	}, [fetchHealth]);

	const copyCmd = (text: string) => {
		navigator.clipboard.writeText(text);
		setCopied(text);
		setTimeout(() => setCopied(""), 2000);
	};

	const uptime = (sec: number) => {
		if (sec < 60) return `${sec}s`;
		if (sec < 3600) return `${Math.floor(sec / 60)}m`;
		if (sec < 86400)
			return `${Math.floor(sec / 3600)}h ${Math.floor((sec % 3600) / 60)}m`;
		return `${Math.floor(sec / 86400)}d ${Math.floor((sec % 86400) / 3600)}h`;
	};

	const filteredTools = tools.filter(
		(t) => !toolSearch || t.toLowerCase().includes(toolSearch.toLowerCase()),
	);

	const StatusDot = ({ color = "green" }: { color?: string }) => (
		<span className={`inline-block w-2 h-2 rounded-full bg-${color}-400`} />
	);

	return (
		<div className="min-h-screen bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950 text-gray-100">
			{/* Header */}
			<header className="sticky top-0 z-50 border-b border-gray-800 bg-gray-950/80 backdrop-blur-xl">
				<div className="flex items-center justify-between px-6 py-3">
					<div className="flex items-center gap-3">
						<div className="w-8 h-8 rounded-lg bg-gradient-to-br from-purple-500 to-blue-500 flex items-center justify-center font-bold text-sm">
							TN
						</div>
						<div>
							<h1 className="text-lg font-bold tracking-tight">TormentNexus</h1>
							<p className="text-xs text-gray-500">
								Universal AI Control Plane
							</p>
						</div>
						{health && (
							<Badge
								variant="outline"
								className="ml-4 bg-green-950/30 text-green-400 border-green-800"
							>
								<CheckCircle2 className="h-3 w-3 mr-1" /> v{health.version} · UP{" "}
								{uptime(health.uptimeSec)}
							</Badge>
						)}
					</div>
					<div className="flex items-center gap-2">
						<Button
							variant="ghost"
							size="sm"
							onClick={fetchHealth}
							disabled={loading.health}
						>
							<RefreshCw
								className={`h-4 w-4 ${loading.health ? "animate-spin" : ""}`}
							/>
						</Button>
					</div>
				</div>
			</header>

			{/* Tabs */}
			<div className="border-b border-gray-800 bg-gray-950/50">
				<div className="px-6">
					<Tabs value={activeTab} onValueChange={setActiveTab}>
						<TabsList className="bg-transparent border-none h-12 gap-1">
							{[
								{ v: "overview", i: MonitorSmartphone, l: "Overview" },
								{ v: "memory", i: Brain, l: "Memory" },
								{ v: "tools", i: Wrench, l: "Tools" },
								{ v: "catalog", i: Database, l: "Catalog" },
								{ v: "agents", i: Users, l: "Agents" },
								{ v: "code", i: FileCode, l: "Code" },
								{ v: "security", i: Shield, l: "Security" },
								{ v: "infra", i: Server, l: "Infrastructure" },
								{ v: "commercial", i: Building2, l: "Commercial" },
							].map((t) => (
								<TabsTrigger
									key={t.v}
									value={t.v}
									className="data-[state=active]:bg-gray-800 data-[state=active]:text-white text-gray-400 hover:text-gray-200 px-3 py-2 rounded-md transition-colors"
								>
									<t.i className="h-4 w-4 mr-1.5" />
									{t.l}
								</TabsTrigger>
							))}
						</TabsList>

						{/* ═══ OVERVIEW ═══ */}
						<TabsContent value="overview" className="mt-0 p-6 space-y-6">
							<div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
								{[
									{
										l: "Status",
										v: health ? "Online" : "Offline",
										i: Activity,
										c: "green",
									},
									{
										l: "Version",
										v: health?.version || "—",
										i: Globe,
										c: "blue",
									},
									{
										l: "Uptime",
										v: health ? uptime(health.uptimeSec) : "—",
										i: Clock,
										c: "purple",
									},
									{
										l: "Tools",
										v: String(tools.length),
										i: Wrench,
										c: "amber",
									},
									{ l: "Memory", v: "L2/L3/L4", i: Database, c: "cyan" },
									{ l: "License", v: "Corporate", i: KeyRound, c: "emerald" },
								].map((s) => (
									<Card key={s.l} className="bg-gray-900/50 border-gray-800">
										<CardContent className="p-4">
											<div className="flex items-center gap-2 mb-1">
												<s.i className={`h-4 w-4 text-${s.c}-400`} />
												<span className="text-xs text-gray-500 uppercase tracking-wider">
													{s.l}
												</span>
											</div>
											<p className="text-lg font-semibold">{s.v}</p>
										</CardContent>
									</Card>
								))}
							</div>

							<Card className="bg-gray-900/50 border-gray-800">
								<CardHeader className="pb-3">
									<CardTitle className="text-sm text-gray-400">
										Quick Actions
									</CardTitle>
								</CardHeader>
								<CardContent className="flex flex-wrap gap-2">
									{[
										{ l: "Search Memory", i: Search, t: "memory" },
										{ l: "Browse Tools", i: Wrench, t: "tools" },
										{ l: "Run Code", i: Play, t: "code" },
										{ l: "View Agents", i: Users, t: "agents" },
										{ l: "Security", i: Shield, t: "security" },
										{ l: "Billing", i: CreditCard, t: "commercial" },
									].map((a) => (
										<Button
											key={a.l}
											variant="outline"
											size="sm"
											onClick={() => setActiveTab(a.t)}
											className="border-gray-700 hover:bg-gray-800"
										>
											<a.i className="h-4 w-4 mr-1.5" />
											{a.l}
										</Button>
									))}
								</CardContent>
							</Card>

							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
								<Card className="bg-gray-900/50 border-gray-800">
									<CardHeader className="pb-2">
										<CardTitle className="text-sm flex items-center gap-2">
											<Brain className="h-4 w-4 text-purple-400" /> Memory
											System
										</CardTitle>
									</CardHeader>
									<CardContent>
										<div className="space-y-2 text-sm">
											{[
												"L2 Vault",
												"L3 Cold Archive",
												"L4 Limbo",
												"GraphRAG",
												"Spaced Repetition",
											].map((m) => (
												<div key={m} className="flex justify-between">
													<span className="text-gray-400">{m}</span>
													<span className="text-green-400">Active</span>
												</div>
											))}
										</div>
									</CardContent>
								</Card>
								<Card className="bg-gray-900/50 border-gray-800">
									<CardHeader className="pb-2">
										<CardTitle className="text-sm flex items-center gap-2">
											<Server className="h-4 w-4 text-blue-400" /> Services
										</CardTitle>
									</CardHeader>
									<CardContent>
										<div className="space-y-2 text-sm">
											{[
												{ l: "Go Kernel", v: `v${health?.version || "?"}` },
												{ l: "Next.js Dashboard", v: "Running" },
												{ l: "Gossip P2P", v: "Port 8190" },
												{ l: "A2A Broker", v: "Active" },
												{ l: "Session Import", v: "9 sessions" },
											].map((s) => (
												<div key={s.l} className="flex justify-between">
													<span className="text-gray-400">{s.l}</span>
													<span className="text-green-400">{s.v}</span>
												</div>
											))}
										</div>
									</CardContent>
								</Card>
								<Card className="bg-gray-900/50 border-gray-800">
									<CardHeader className="pb-2">
										<CardTitle className="text-sm flex items-center gap-2">
											<Shield className="h-4 w-4 text-emerald-400" /> Security
										</CardTitle>
									</CardHeader>
									<CardContent>
										<div className="space-y-2 text-sm">
											{[
												{ l: "License", v: "Corporate" },
												{ l: "RBAC", v: "Enforced" },
												{ l: "Audit Logging", v: "Active" },
												{ l: "SSO/OIDC", v: "Configured" },
												{ l: "Stripe Billing", v: "Live" },
											].map((s) => (
												<div key={s.l} className="flex justify-between">
													<span className="text-gray-400">{s.l}</span>
													<span className="text-green-400">{s.v}</span>
												</div>
											))}
										</div>
									</CardContent>
								</Card>
							</div>

							<Card className="bg-gray-900/50 border-gray-800">
								<CardHeader className="pb-2">
									<CardTitle className="text-sm flex items-center gap-2">
										<Terminal className="h-4 w-4 text-amber-400" /> Quick
										Install
									</CardTitle>
								</CardHeader>
								<CardContent>
									<div className="space-y-2">
										{[
											{
												cmd: "npx @tormentnexus/install",
												desc: "One-command setup for 38+ AI clients",
											},
											{
												cmd: "npm install -g @tormentnexus/cli",
												desc: "CLI tools for memory and tool search",
											},
											{ cmd: "tn status", desc: "Check system health" },
											{
												cmd: 'tn search "your topic"',
												desc: "Search L2 memory vault",
											},
										].map((item) => (
											<div
												key={item.cmd}
												className="flex items-center gap-3 p-2 rounded-md bg-gray-800/50 hover:bg-gray-800 transition-colors group"
											>
												<code className="text-sm text-purple-300 font-mono flex-1">
													{item.cmd}
												</code>
												<span className="text-xs text-gray-500">
													{item.desc}
												</span>
												<Button
													variant="ghost"
													size="sm"
													className="h-6 w-6 p-0 opacity-0 group-hover:opacity-100"
													onClick={() => copyCmd(item.cmd)}
												>
													{copied === item.cmd ? (
														<Check className="h-3 w-3 text-green-400" />
													) : (
														<Copy className="h-3 w-3" />
													)}
												</Button>
											</div>
										))}
									</div>
								</CardContent>
							</Card>
						</TabsContent>

						{/* ═══ MEMORY ═══ */}
						<TabsContent value="memory" className="mt-0 p-6 space-y-6">
							<Card className="bg-gray-900/50 border-gray-800">
								<CardHeader className="pb-3">
									<CardTitle className="text-sm flex items-center gap-2">
										<Search className="h-4 w-4 text-purple-400" /> Search Memory
										(L2 Vault)
									</CardTitle>
									<CardDescription>
										Search your persistent memory store using semantic vector
										search
									</CardDescription>
								</CardHeader>
								<CardContent>
									<div className="flex gap-2">
										<Input
											placeholder="Search memories..."
											value={memorySearch}
											onChange={(e) => setMemorySearch(e.target.value)}
											onKeyDown={(e) => e.key === "Enter" && searchMemory()}
											className="bg-gray-800 border-gray-700"
										/>
										<Button
											onClick={searchMemory}
											disabled={loading.memorySearch}
										>
											{loading.memorySearch ? (
												<RefreshCw className="h-4 w-4 animate-spin" />
											) : (
												<Search className="h-4 w-4" />
											)}
										</Button>
									</div>
									{memoryResults.length > 0 && (
										<div className="mt-4 space-y-2 max-h-96 overflow-y-auto">
											{memoryResults.map((mem, i) => (
												<div
													key={i}
													className="p-3 rounded-lg bg-gray-800/50 border border-gray-700"
												>
													<p className="text-sm">{mem.content}</p>
													<div className="mt-2 flex items-center gap-2 flex-wrap">
														{mem.tags?.map((tag) => (
															<Badge
																key={tag}
																variant="secondary"
																className="text-xs"
															>
																{tag}
															</Badge>
														))}
														<span className="text-xs text-gray-500">
															{mem.category}
														</span>
													</div>
												</div>
											))}
										</div>
									)}
								</CardContent>
							</Card>
							<Card className="bg-gray-900/50 border-gray-800">
								<CardHeader className="pb-3">
									<CardTitle className="text-sm flex items-center gap-2">
										<Database className="h-4 w-4 text-blue-400" /> Store Memory
									</CardTitle>
									<CardDescription>
										Save information to the L2 vault with tags for future
										retrieval
									</CardDescription>
								</CardHeader>
								<CardContent className="space-y-3">
									<Textarea
										placeholder="Enter memory content..."
										value={memoryStore}
										onChange={(e) => setMemoryStore(e.target.value)}
										className="bg-gray-800 border-gray-700 min-h-24"
									/>
									<div className="flex gap-2">
										<Input
											placeholder="Tags (comma-separated)"
											value={memoryStoreTags}
											onChange={(e) => setMemoryStoreTags(e.target.value)}
											className="bg-gray-800 border-gray-700"
										/>
										<Button
											onClick={storeMemory}
											disabled={loading.memoryStore || !memoryStore.trim()}
										>
											{loading.memoryStore ? (
												<RefreshCw className="h-4 w-4 animate-spin" />
											) : (
												<>
													<Plus className="h-4 w-4 mr-1" /> Store
												</>
											)}
										</Button>
									</div>
								</CardContent>
							</Card>
							<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
								{[
									{
										n: "L2 Vault",
										d: "Active persistent memory with vector embeddings",
										s: "Active",
										c: "purple",
										i: Brain,
									},
									{
										n: "L3 Cold Archive",
										d: "Compressed long-term storage for older memories",
										s: "Active",
										c: "blue",
										i: HardDrive,
									},
									{
										n: "L4 Limbo",
										d: "Quarantined memories pending review or deletion",
										s: "Active",
										c: "amber",
										i: AlertTriangle,
									},
								].map((tier) => (
									<Card key={tier.n} className="bg-gray-900/50 border-gray-800">
										<CardHeader className="pb-2">
											<CardTitle className="text-sm flex items-center gap-2">
												<tier.i className={`h-4 w-4 text-${tier.c}-400`} />{" "}
												{tier.n}
											</CardTitle>
										</CardHeader>
										<CardContent>
											<p className="text-xs text-gray-400 mb-2">{tier.d}</p>
											<Badge variant="outline">{tier.s}</Badge>
										</CardContent>
									</Card>
								))}
							</div>
						</TabsContent>

						{/* ═══ TOOLS ═══ */}
						<TabsContent value="tools" className="mt-0 p-6">
							<Card className="bg-gray-900/50 border-gray-800">
								<CardHeader className="pb-3">
									<div className="flex items-center justify-between">
										<div>
											<CardTitle className="text-sm flex items-center gap-2">
												<Wrench className="h-4 w-4 text-amber-400" /> MCP Tool
												Registry
											</CardTitle>
											<CardDescription>
												{tools.length} registered tools across all MCP servers
											</CardDescription>
										</div>
										<Badge
											variant="outline"
											className="bg-amber-950/30 text-amber-400 border-amber-800"
										>
											{tools.length} Tools
										</Badge>
									</div>
								</CardHeader>
								<CardContent>
									<Input
										placeholder="Filter tools..."
										value={toolSearch}
										onChange={(e) => setToolSearch(e.target.value)}
										className="bg-gray-800 border-gray-700 mb-4"
									/>
									<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2 max-h-[500px] overflow-y-auto">
										{filteredTools.map((tool) => (
											<div
												key={tool}
												className="flex items-center gap-2 p-2 rounded-md bg-gray-800/50 hover:bg-gray-800 transition-colors"
											>
												<Wrench className="h-3 w-3 text-amber-400 flex-shrink-0" />
												<span className="text-sm truncate">{tool}</span>
											</div>
										))}
									</div>
								</CardContent>
							</Card>
						</TabsContent>

						{/* ═══ CATALOG ═══ */}
						<TabsContent value="catalog" className="mt-0 p-6">
							<CatalogBrowser />
						</TabsContent>

						{/* ═══ AGENTS ═══ */}
						<TabsContent value="agents" className="mt-0 p-6">
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
								{[
									{
										n: "Go Coder Agent",
										d: "Native Go code generation and refactoring",
										s: "Idle",
										c: ["coding", "refactoring", "filesystem"],
									},
									{
										n: "A2A Broker",
										d: "Agent-to-Agent message routing and negotiation",
										s: "Active",
										c: ["routing", "bidding", "heartbeat"],
									},
									{
										n: "Swarm Director",
										d: "Task orchestration and agent coordination",
										s: "Active",
										c: ["planning", "delegation", "consensus"],
									},
									{
										n: "Session Keeper",
										d: "Session persistence and context management",
										s: "Active",
										c: ["persistence", "context", "recovery"],
									},
									{
										n: "Memory Harvester",
										d: "Automatic memory extraction and consolidation",
										s: "Active",
										c: ["extraction", "consolidation", "decay"],
									},
									{
										n: "Catalog Sync",
										d: "MCP tool registry synchronization from Glama.ai",
										s: "Active",
										c: ["sync", "index", "update"],
									},
								].map((a) => (
									<Card key={a.n} className="bg-gray-900/50 border-gray-800">
										<CardHeader className="pb-2">
											<div className="flex items-center justify-between">
												<CardTitle className="text-sm">{a.n}</CardTitle>
												<Badge
													variant={a.s === "Active" ? "default" : "secondary"}
													className="text-xs"
												>
													{a.s}
												</Badge>
											</div>
										</CardHeader>
										<CardContent>
											<p className="text-xs text-gray-400 mb-3">{a.d}</p>
											<div className="flex flex-wrap gap-1">
												{a.c.map((cap) => (
													<Badge
														key={cap}
														variant="outline"
														className="text-xs border-gray-700"
													>
														{cap}
													</Badge>
												))}
											</div>
										</CardContent>
									</Card>
								))}
							</div>
						</TabsContent>

						{/* ═══ CODE ═══ */}
						<TabsContent value="code" className="mt-0 p-6">
							<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
								<Card className="bg-gray-900/50 border-gray-800">
									<CardHeader className="pb-2">
										<CardTitle className="text-sm flex items-center gap-2">
											<Play className="h-4 w-4 text-green-400" /> Code Execution
										</CardTitle>
										<CardDescription>
											Execute code in sandboxed environments
										</CardDescription>
									</CardHeader>
									<CardContent>
										<Textarea
											placeholder="Enter code..."
											className="bg-gray-800 border-gray-700 font-mono text-sm min-h-32"
										/>
										<Button className="mt-2" size="sm">
											<Play className="h-4 w-4 mr-1" /> Run
										</Button>
									</CardContent>
								</Card>
								<Card className="bg-gray-900/50 border-gray-800">
									<CardHeader className="pb-2">
										<CardTitle className="text-sm flex items-center gap-2">
											<GitBranch className="h-4 w-4 text-blue-400" /> Code
											Intelligence
										</CardTitle>
										<CardDescription>
											AST analysis, dependency graphs, symbol navigation
										</CardDescription>
									</CardHeader>
									<CardContent className="space-y-2">
										{[
											{
												l: "Dependency Graph",
												i: Network,
												d: "Visualize package dependencies",
											},
											{
												l: "Symbol Search",
												i: Search,
												d: "Find definitions and references",
											},
											{
												l: "Code Fix",
												i: Wrench,
												d: "Auto-diagnose and fix issues",
											},
											{
												l: "AST Explorer",
												i: Code,
												d: "Parse and analyze syntax trees",
											},
										].map((f) => (
											<div
												key={f.l}
												className="flex items-center gap-3 p-2 rounded-md bg-gray-800/50 hover:bg-gray-800 cursor-pointer"
											>
												<f.i className="h-4 w-4 text-gray-400" />
												<div>
													<p className="text-sm font-medium">{f.l}</p>
													<p className="text-xs text-gray-500">{f.d}</p>
												</div>
												<ChevronRight className="h-4 w-4 text-gray-600 ml-auto" />
											</div>
										))}
									</CardContent>
								</Card>
							</div>
						</TabsContent>

						{/* ═══ SECURITY ═══ */}
						<TabsContent value="security" className="mt-0 p-6">
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
								{[
									{
										n: "RBAC Policies",
										d: "Role-based access control for tool calls and API endpoints",
										i: Shield,
										s: "Enforced",
									},
									{
										n: "Audit Logging",
										d: "Complete activity trail for compliance and forensics",
										i: Eye,
										s: "Active",
									},
									{
										n: "SSO/OIDC",
										d: "Single sign-on with SAML and OAuth2 providers",
										i: Fingerprint,
										s: "Configured",
									},
									{
										n: "License Manager",
										d: "Ed25519-verified commercial license with seat limits",
										i: KeyRound,
										s: "Valid",
									},
									{
										n: "Tenant Isolation",
										d: "Docker container isolation per organization",
										i: Lock,
										s: "3 Tenants",
									},
									{
										n: "Op Blocker",
										d: "Blocks destructive shell operations automatically",
										i: AlertTriangle,
										s: "Active",
									},
								].map((item) => (
									<Card key={item.n} className="bg-gray-900/50 border-gray-800">
										<CardHeader className="pb-2">
											<div className="flex items-center gap-2">
												<item.i className="h-4 w-4 text-emerald-400" />
												<CardTitle className="text-sm">{item.n}</CardTitle>
											</div>
										</CardHeader>
										<CardContent>
											<p className="text-xs text-gray-400 mb-2">{item.d}</p>
											<Badge
												variant="outline"
												className="bg-emerald-950/30 text-emerald-400 border-emerald-800"
											>
												{item.s}
											</Badge>
										</CardContent>
									</Card>
								))}
							</div>
						</TabsContent>

						{/* ═══ INFRASTRUCTURE ═══ */}
						<TabsContent value="infra" className="mt-0 p-6">
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
								{[
									{
										n: "Go Kernel",
										d: "Primary sidecar process on port 8090",
										i: Server,
										s: `v${health?.version || "?"} UP ${uptime(health?.uptimeSec || 0)}`,
									},
									{
										n: "Docker Tenants",
										d: "Isolated tenant containers with sidecar + web",
										i: Box,
										s: "3 Active",
									},
									{
										n: "Gossip P2P",
										d: "UDP peer-to-peer memory synchronization",
										i: Wifi,
										s: "Port 8190",
									},
									{
										n: "Marketing Agent",
										d: "Automated outreach and lead management",
										i: Send,
										s: "Running",
									},
									{
										n: "Nginx Proxy",
										d: "Reverse proxy with TLS termination",
										i: Globe,
										s: "4 Sites",
									},
									{
										n: "PM2 Manager",
										d: "Auto-restart and log management",
										i: Gauge,
										s: "2 Services",
									},
								].map((item) => (
									<Card key={item.n} className="bg-gray-900/50 border-gray-800">
										<CardHeader className="pb-2">
											<div className="flex items-center gap-2">
												<item.i className="h-4 w-4 text-blue-400" />
												<CardTitle className="text-sm">{item.n}</CardTitle>
											</div>
										</CardHeader>
										<CardContent>
											<p className="text-xs text-gray-400 mb-2">{item.d}</p>
											<Badge
												variant="outline"
												className="bg-blue-950/30 text-blue-400 border-blue-800"
											>
												{item.s}
											</Badge>
										</CardContent>
									</Card>
								))}
							</div>
						</TabsContent>

						{/* ═══ COMMERCIAL ═══ */}
						<TabsContent value="commercial" className="mt-0 p-6 space-y-4">
							<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
								<Card className="bg-gray-900/50 border-gray-800">
									<CardHeader className="pb-2">
										<CardTitle className="text-sm flex items-center gap-2">
											<CreditCard className="h-4 w-4 text-purple-400" /> Stripe
											Billing
										</CardTitle>
									</CardHeader>
									<CardContent className="space-y-2 text-sm">
										<div className="flex justify-between">
											<span className="text-gray-400">Plan</span>
											<span>Professional</span>
										</div>
										<div className="flex justify-between">
											<span className="text-gray-400">Seats</span>
											<span>50</span>
										</div>
										<div className="flex justify-between">
											<span className="text-gray-400">Status</span>
											<span className="text-green-400">Active</span>
										</div>
										<div className="flex justify-between">
											<span className="text-gray-400">License</span>
											<span>Corporate (Ed25519)</span>
										</div>
									</CardContent>
								</Card>
								<Card className="bg-gray-900/50 border-gray-800">
									<CardHeader className="pb-2">
										<CardTitle className="text-sm flex items-center gap-2">
											<Store className="h-4 w-4 text-amber-400" /> Marketplace
										</CardTitle>
									</CardHeader>
									<CardContent>
										{[
											{ n: "AI Code Review Agent", i: 1247 },
											{ n: "Terraform Planner", i: 892 },
											{ n: "Database Migration Tool", i: 634 },
										].map((item) => (
											<div
												key={item.n}
												className="flex items-center justify-between p-2 rounded-md bg-gray-800/50 mb-2"
											>
												<span className="text-sm">{item.n}</span>
												<Badge variant="secondary" className="text-xs">
													{item.i} installs
												</Badge>
											</div>
										))}
									</CardContent>
								</Card>
							</div>
							<Card className="bg-gray-900/50 border-gray-800">
								<CardHeader className="pb-2">
									<CardTitle className="text-sm flex items-center gap-2">
										<Fingerprint className="h-4 w-4 text-cyan-400" /> SSO / OIDC
										Configuration
									</CardTitle>
								</CardHeader>
								<CardContent>
									<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
										<div>
											<label className="text-xs text-gray-400">
												OIDC Issuer URL
											</label>
											<Input
												placeholder="https://identity.example.com"
												className="bg-gray-800 border-gray-700 mt-1"
											/>
										</div>
										<div>
											<label className="text-xs text-gray-400">Client ID</label>
											<Input
												placeholder="tormentnexus-client"
												className="bg-gray-800 border-gray-700 mt-1"
											/>
										</div>
									</div>
									<Button className="mt-3" size="sm" variant="outline">
										<Lock className="h-4 w-4 mr-1" /> Save Configuration
									</Button>
								</CardContent>
							</Card>
						</TabsContent>
					</Tabs>
				</div>
			</div>

			{/* Footer */}
			<footer className="border-t border-gray-800 bg-gray-950/50 px-6 py-3">
				<div className="flex items-center justify-between text-xs text-gray-500">
					<span>TormentNexus v{health?.version || "?"} · MIT License</span>
					<div className="flex items-center gap-4">
						<a
							href="https://github.com/MDMAtk/TormentNexus"
							className="hover:text-gray-300"
						>
							GitHub
						</a>
						<a
							href="https://hypernexus.site/docs"
							className="hover:text-gray-300"
						>
							Docs
						</a>
						<a
							href="https://www.npmjs.com/search?q=%40tormentnexus"
							className="hover:text-gray-300"
						>
							npm
						</a>
					</div>
				</div>
			</footer>
		</div>
	);
}
