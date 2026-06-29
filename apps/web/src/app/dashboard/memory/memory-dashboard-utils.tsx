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
}

export interface MemoryPivotAction {
  label: string;
  handler: () => void;
  group?: string;
  description?: string;
  query?: string;
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

export type MemorySearchMode = "all" | "fts" | "semantic" | "pivot" | "agent" | "facts" | "observations" | "prompts" | "session_summaries";

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

export function getMemoryDetailSections(record: MemoryRecord): { label: string; value: string }[] {
  return [
    { label: "ID", value: record.id },
    { label: "Session", value: record.session_id },
    { label: "Kind", value: record.memory_kind },
    { label: "Category", value: record.category },
    { label: "Heat", value: String(record.heat_score?.toFixed(1) ?? "?") },
    { label: "Importance", value: String(record.importance?.toFixed(2) ?? "?") },
  ];
}

export function getMemoryModeHint(mode: MemorySearchMode): string {
  switch (mode) {
    case "fts": return "Full-text search across all L2 memories using BM25";
    case "semantic": return "Semantic vector similarity search";
    case "pivot": return "Pivot-based context retrieval";
    case "agent": return "Agent-specific memory search";
  }
}

export function getMemoryPivotSections(record: MemoryRecord): MemoryPivotAction[] {
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

export function groupMemoryRecordsByDay(records: MemoryRecord[]): Record<string, MemoryRecord[]> {
  const groups: Record<string, MemoryRecord[]> = {};
  for (const r of records) {
    const day = (r.created_at || "").slice(0, 10) || "unknown";
    if (!groups[day]) groups[day] = [];
    groups[day].push(r);
  }
  return groups;
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
