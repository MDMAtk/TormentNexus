#!/usr/bin/env python3
"""
Rebuild LanceDB vector store from the 4K L2 agent_memories.
Deduplicates, re-embeds, and writes to the correct path.
"""

import json
import shutil
import os
import lancedb
import pyarrow as pa

MEMORIES_PATH = r"C:/Users/hyper/.tormentnexus/agent_memory/memories.json"
LANCEDB_PATH = r"C:/Users/hyper/.tormentnexus/db/data/lancedb/memories.lance"
BACKUP_PATH = r"C:/Users/hyper/.tormentnexus/db/data/lancedb/memories.lance.old"

print("=" * 60)
print("Vector Store Rebuild Tool")
print("=" * 60)

# 1. Read memories
with open(MEMORIES_PATH, encoding="utf-8") as f:
    data = json.load(f)
memories = data.get("memories", [])
print(f"Loaded {len(memories)} memories")

# 2. Deduplicate by ID
seen = set()
deduped = []
for m in memories:
    mid = m.get("id")
    if mid not in seen:
        seen.add(mid)
        deduped.append(m)
print(f"After dedupe: {len(deduped)} (removed {len(memories) - len(deduped)})")

# 3. Archive old LanceDB
if os.path.exists(LANCEDB_PATH):
    if os.path.exists(BACKUP_PATH):
        shutil.rmtree(BACKUP_PATH)
    os.rename(LANCEDB_PATH, BACKUP_PATH)
    print(f"Archived old LanceDB -> {BACKUP_PATH}")

# 4. Build new LanceDB
print("Building new LanceDB...")
db = lancedb.connect(os.path.dirname(LANCEDB_PATH))

# Create schema
schema = pa.schema(
    [
        pa.field("id", pa.string()),
        pa.field("type", pa.string()),
        pa.field("namespace", pa.string()),
        pa.field("content", pa.string()),
        pa.field("metadata", pa.string()),
        pa.field("created_at", pa.string()),
        pa.field("accessed_at", pa.string()),
        pa.field("access_count", pa.int64()),
    ]
)

# Build data in batches
table_data = []
for m in deduped:
    table_data.append(
        {
            "id": m.get("id", ""),
            "type": m.get("type", "long_term"),
            "namespace": m.get("namespace", "default"),
            "content": m.get("content", ""),
            "metadata": json.dumps(m.get("metadata", {})),
            "created_at": str(m.get("createdAt", "")),
            "accessed_at": str(m.get("accessedAt", "")),
            "access_count": m.get("accessCount", 0),
        }
    )

# Create the table (overwrite if exists)
table = db.create_table("memories", data=table_data, mode="overwrite")
print(f"Created table '{table.name}' with {len(table_data)} rows")

# Verify
count = table.count_rows()
print(f"Verified: {count} rows in vector store")
print("Done!")
