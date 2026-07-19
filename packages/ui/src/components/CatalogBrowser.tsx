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
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
	Search,
	ExternalLink,
	Database,
	Wrench,
	Brain,
	Code,
	Users,
	FileText,
	RefreshCw,
} from "lucide-react";

interface CatalogEntry {
	title: string;
	url: string;
	description?: string;
	source: string;
	category: string;
}

interface CatalogStats {
	total: number;
	bySource: Record<string, number>;
	byCategory: Record<string, number>;
	skills: number;
	goHandlers: number;
}

interface CategoryInfo {
	label: string;
	description: string;
	count: number;
}

const categoryIcons: Record<string, any> = {
	mcp_server: Wrench,
	ai_dev_tool: Code,
	mcp_language: Code,
	mcp_category: Wrench,
	prompt: FileText,
	agent: Users,
	skill: Brain,
};

const categoryColors: Record<string, string> = {
	mcp_server: "purple",
	ai_dev_tool: "blue",
	mcp_language: "cyan",
	mcp_category: "amber",
	prompt: "orange",
	agent: "green",
	skill: "pink",
};

export function CatalogBrowser() {
	const [query, setQuery] = useState("");
	const [results, setResults] = useState<CatalogEntry[]>([]);
	const [stats, setStats] = useState<CatalogStats | null>(null);
	const [categories, setCategories] = useState<Record<string, CategoryInfo>>(
		{},
	);
	const [activeCategory, setActiveCategory] = useState("");
	const [loading, setLoading] = useState(false);
	const [total, setTotal] = useState(0);

	const fetchStats = useCallback(async () => {
		try {
			const res = await fetch("/api/backlog/stats");
			const data = await res.json();
			setStats(data);
		} catch {
			/* offline */
		}
	}, []);

	const fetchCategories = useCallback(async () => {
		try {
			const res = await fetch("/api/backlog/categories");
			const data = await res.json();
			setCategories(data.categories || {});
		} catch {
			/* offline */
		}
	}, []);

	const search = useCallback(async (q: string, cat: string) => {
		setLoading(true);
		try {
			const params = new URLSearchParams();
			if (q) params.set("q", q);
			if (cat) params.set("category", cat);
			params.set("limit", "50");

			const res = await fetch(`/api/backlog/search?${params}`);
			const data = await res.json();
			setResults(data.results || []);
			setTotal(data.total || 0);
		} catch {
			setResults([]);
		}
		setLoading(false);
	}, []);

	useEffect(() => {
		fetchStats();
		fetchCategories();
	}, [fetchStats, fetchCategories]);

	useEffect(() => {
		const timer = setTimeout(() => {
			search(query, activeCategory);
		}, 300);
		return () => clearTimeout(timer);
	}, [query, activeCategory, search]);

	return (
		<div className="space-y-6">
			{/* Header */}
			<div className="flex items-center justify-between">
				<div>
					<h2 className="text-2xl font-bold">The Index</h2>
					<p className="text-muted-foreground">
						{stats
							? `${stats.total.toLocaleString()} entries from ${Object.keys(stats.bySource).length} sources`
							: "Loading..."}
					</p>
				</div>
				<Button variant="outline" size="sm" onClick={fetchStats}>
					<RefreshCw className="h-4 w-4 mr-1" /> Refresh
				</Button>
			</div>

			{/* Stats */}
			{stats && (
				<div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-3">
					{Object.entries(categories).map(([key, cat]) => {
						const Icon = categoryIcons[key] || Database;
						const color = categoryColors[key] || "gray";
						return (
							<Card
								key={key}
								className={`cursor-pointer transition-all hover:border-${color}-500 ${activeCategory === key ? `border-${color}-500 bg-${color}-950/20` : ""}`}
								onClick={() =>
									setActiveCategory(activeCategory === key ? "" : key)
								}
							>
								<CardContent className="p-3 text-center">
									<Icon className={`h-5 w-5 mx-auto mb-1 text-${color}-400`} />
									<div className="text-lg font-bold">
										{cat.count.toLocaleString()}
									</div>
									<div className="text-xs text-muted-foreground">
										{cat.label}
									</div>
								</CardContent>
							</Card>
						);
					})}
				</div>
			)}

			{/* Search */}
			<div className="flex gap-2">
				<div className="relative flex-1">
					<Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
					<Input
						placeholder="Search MCP servers, skills, prompts..."
						value={query}
						onChange={(e) => setQuery(e.target.value)}
						className="pl-10"
					/>
				</div>
				{activeCategory && (
					<Button variant="outline" onClick={() => setActiveCategory("")}>
						Clear filter
					</Button>
				)}
			</div>

			{/* Results */}
			<div className="space-y-2">
				{loading ? (
					<div className="text-center py-8 text-muted-foreground">
						Searching...
					</div>
				) : results.length === 0 ? (
					<div className="text-center py-8 text-muted-foreground">
						{query
							? `No results for "${query}"`
							: "Enter a search term to explore the catalog"}
					</div>
				) : (
					<>
						<div className="text-sm text-muted-foreground mb-2">
							{total.toLocaleString()} results{" "}
							{activeCategory && `in ${categories[activeCategory]?.label}`}
						</div>
						<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
							{results.map((entry, i) => (
								<Card
									key={i}
									className="hover:border-purple-500/50 transition-colors"
								>
									<CardContent className="p-4">
										<div className="flex items-start justify-between gap-2">
											<div className="flex-1 min-w-0">
												<h3 className="font-medium text-sm truncate">
													{entry.title.replace(/\*\*/g, "")}
												</h3>
												{entry.description && (
													<p className="text-xs text-muted-foreground mt-1 line-clamp-2">
														{entry.description}
													</p>
												)}
												<div className="flex items-center gap-2 mt-2">
													<Badge variant="secondary" className="text-xs">
														{entry.source}
													</Badge>
													<Badge variant="outline" className="text-xs">
														{entry.category}
													</Badge>
												</div>
											</div>
											<a
												href={entry.url}
												target="_blank"
												rel="noopener noreferrer"
												className="text-muted-foreground hover:text-foreground"
											>
												<ExternalLink className="h-4 w-4" />
											</a>
										</div>
									</CardContent>
								</Card>
							))}
						</div>
					</>
				)}
			</div>
		</div>
	);
}
