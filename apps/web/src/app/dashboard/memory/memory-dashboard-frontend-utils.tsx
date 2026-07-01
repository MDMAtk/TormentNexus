export interface MemoryRecord {
	id: string;
	session_id: string;
	memory_type: string;
	memory_kind: string;
	category: string;
	tags?: string;
	source_url?: string;
	content: string;
	importance: number;
	heat_score: number;
	last_accessed_at?: string;
	created_at?: string;
	metadata?: {
		source?: string;
		type?: string;
		memoryKind?: string;
		structuredObservation?: {
			type?: string;
			title?: string;
			subtitle?: string;
			narrative?: string;
			facts?: string[];
			concepts?: string[];
			filesRead?: string[];
			filesModified?: string[];
			toolName?: string;
		};
		structuredUserPrompt?: {
			role?: string;
			content?: string;
			promptNumber?: number;
			sessionId?: string;
			activeGoal?: string | null;
			lastObjective?: string | null;
		};
		structuredSessionSummary?: {
			name?: string;
			sessionId?: string;
			status?: string;
			cliType?: string;
			activeGoal?: string | null;
			lastObjective?: string | null;
			restartCount?: number;
		};
	};
}

export interface MemoryPivotAction {
	key?: string;
	label: string;
	handler: () => void;
	group?: string;
	description?: string;
	query?: string;
	mode?: MemorySearchMode;
	memoryId?: string;
	anchorTime?: string;
	windowSize?: number;
	limit?: number;
}

export interface RelatedMemoryRecord {
	id: string;
	content: string;
	relation_type: string;
	weight: number;
}

export type MemorySearchMode =
	| "all"
	| "fts"
	| "semantic"
	| "pivot"
	| "agent"
	| "facts"
	| "observations"
	| "prompts"
	| "session_summaries";

export const MEMORY_MODEL_PILLARS = [
	{ key: "importance", label: "Importance", color: "text-emerald-400" },
	{ key: "heat", label: "Heat Score", color: "text-amber-400" },
	{ key: "recency", label: "Recency", color: "text-blue-400" },
];

export const MEMORY_SEARCH_MODES: { key: MemorySearchMode; label: string }[] = [
	{ key: "fts", label: "Full-Text" },
	{ key: "semantic", label: "Semantic" },
	{ key: "pivot", label: "Pivot" },
	{ key: "agent", label: "Agent" },
];

export function getMemoryTitle(record: MemoryRecord): string {
	const content = record.content || "";
	const firstLine = content.split("\n")[0] || "";
	return firstLine.slice(0, 80) || record.id.slice(0, 30);
}

export function getMemoryBadgeLabel(record: MemoryRecord): string {
	return record.memory_kind || record.memory_type || "memory";
}

export interface MemoryDetailSection {
	title: string;
	body?: string;
	items?: string[];
}

export function getMemoryDetailSections(
	record: MemoryRecord,
): MemoryDetailSection[] {
	const sections: MemoryDetailSection[] = [];
	const observation = record.metadata?.structuredObservation;
	const prompt = record.metadata?.structuredUserPrompt;
	const summary = record.metadata?.structuredSessionSummary;

	if (observation) {
		if (observation.narrative) {
			sections.push({ title: "Narrative", body: observation.narrative });
		}
		if (observation.subtitle) {
			sections.push({ title: "Subtitle", body: observation.subtitle });
		}
		if (observation.facts && observation.facts.length > 0) {
			sections.push({ title: "Extracted facts", items: observation.facts });
		}
		if (observation.concepts && observation.concepts.length > 0) {
			sections.push({ title: "Concepts", items: observation.concepts });
		}
		if (observation.filesRead && observation.filesRead.length > 0) {
			sections.push({ title: "Files read", items: observation.filesRead });
		}
		if (observation.filesModified && observation.filesModified.length > 0) {
			sections.push({ title: "Files modified", items: observation.filesModified });
		}
	} else if (prompt) {
		sections.push({ title: "Prompt content", body: prompt.content || record.content });
		const anchors = [prompt.activeGoal, prompt.lastObjective].filter(Boolean) as string[];
		if (anchors.length > 0) {
			sections.push({ title: "Intent anchors", items: anchors });
		}
	} else if (summary) {
		if (summary.activeGoal) {
			sections.push({ title: "Active goal", body: summary.activeGoal });
		}
	} else {
		sections.push({
			title: "Metadata Details",
			items: [
				`ID: ${record.id || ""}`,
				`Session: ${record.session_id || ""}`,
				`Kind: ${record.memory_kind || ""}`,
				`Category: ${record.category || ""}`,
				`Heat Score: ${record.heat_score || 0}`,
				`Importance: ${record.importance || 0}`
			]
		});
	}
	return sections;
}

export function getMemoryModeHint(mode: MemorySearchMode, tier?: string): string {
	switch (mode) {
		case "fts":
			return `Full-text search across all ${tier || "L2"} memories using BM25`;
		case "semantic":
			return `Semantic vector similarity search across ${tier || "L2"} memories`;
		case "pivot":
			return "Pivot-based context retrieval";
		case "agent":
			return "Agent-specific memory search";
	}
	return `Search mode ${mode}`;
}

export interface MemoryPivotSection {
	title: string;
	actions: MemoryPivotAction[];
}

export function getMemoryPivotSections(
	record: MemoryRecord,
): MemoryPivotSection[] {
	return [];
}

export function getMemoryPreview(record: MemoryRecord, maxLen = 200): string {
	return (record.content || "").slice(0, maxLen);
}

export function getMemoryProvenance(record: MemoryRecord): string {
	return record.session_id || "unknown";
}

export function filterMemoryRecords(
	records: MemoryRecord[],
	query: string,
	kind?: string,
): MemoryRecord[] {
	if (!query && !kind) return records;
	return records.filter((r) => {
		if (kind && r.memory_kind !== kind) return false;
		if (query) {
			const q = query.toLowerCase();
			return (
				r.content.toLowerCase().includes(q) ||
				r.id.toLowerCase().includes(q) ||
				r.session_id.toLowerCase().includes(q)
			);
		}
		return true;
	});
}

export function groupMemoryWindowAroundAnchor(
	records: MemoryRecord[],
	anchorId: string,
	windowSize = 10,
): MemoryRecord[] {
	const idx = records.findIndex((r) => r.id === anchorId);
	if (idx < 0) return records.slice(0, windowSize);
	const start = Math.max(0, idx - Math.floor(windowSize / 2));
	return records.slice(start, start + windowSize);
}
export interface MemoryTimelineGroup {
	key: string;
	label: string;
	items: MemoryRecord[];
}

export function groupMemoryRecordsByDay(
	records: MemoryRecord[],
): MemoryTimelineGroup[] {
	const groups = new Map<string, MemoryTimelineGroup>();

	for (const memory of sortMemoryRecordsByTimestamp(records)) {
		const timestamp = getMemoryTimestamp(memory);
		if (!timestamp) continue;
		const date = new Date(timestamp);
		const key = `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`;
		const existing = groups.get(key);

		if (existing) {
			existing.items.push(memory);
			continue;
		}

		groups.set(key, {
			key,
			label: date.toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }),
			items: [memory],
		});
	}

	return Array.from(groups.values());
}

export function getMemoryRecordKey(record: MemoryRecord): string {
	return record.id;
}

export function getMemorySessionId(record: MemoryRecord): string {
	return record.session_id || "unknown";
}

export function getMemoryTimestamp(record: MemoryRecord): string {
	return record.created_at || record.last_accessed_at || "";
}

export function getRelatedMemoryRecords(
	record: MemoryRecord,
	_allRecords: MemoryRecord[],
): RelatedMemoryRecord[] {
	return [];
}

export function sortMemoryRecordsByTimestamp(
	records: MemoryRecord[],
	ascending = false,
): MemoryRecord[] {
	return [...records].sort((a, b) => {
		const tA = a.created_at ? new Date(a.created_at).getTime() : 0;
		const tB = b.created_at ? new Date(b.created_at).getTime() : 0;
		return ascending ? tA - tB : tB - tA;
	});
}
