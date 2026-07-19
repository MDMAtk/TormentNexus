#!/usr/bin/env python3
"""
Extract memories from Brain log, imported sessions, and Pi blocks.
Merge, dedupe, and append to L2 (agent_memory) and L3 (cold archive).
"""

import json
import sqlite3
import re
import os
from datetime import datetime

MEMORIES_PATH = r"C:/Users/hyper/.tormentnexus/agent_memory/memories.json"
L3_DB_PATH = r"C:/Users/hyper/.tormentnexus/l3_cold_archive.db"
BRAIN_LOG = r"C:/Users/hyper/workspace/tormentnexus/.memory/branches/main/log.md"
SESSION_DB = r"C:/Users/hyper/workspace/tormentnexus/go/tormentnexus.db"

print("=" * 60)
print("Memory Extraction & Merge Tool")
print("=" * 60)

# 1. Load existing L2 memories
with open(MEMORIES_PATH, encoding="utf-8") as f:
    existing_data = json.load(f)
existing = existing_data.get("memories", [])
existing_ids = set(m.get("id") for m in existing)
print(f"\nExisting L2: {len(existing)} memories")

# 2. Extract from Brain log
print("\n--- Extracting from Brain log ---")
with open(BRAIN_LOG, encoding="utf-8") as f:
    log = f.read()

turn_pattern = re.compile(r"(## Turn \d+ \| .+?)(?=\n## Turn |\Z)", re.DOTALL)
turns = turn_pattern.findall(log)
print(f"Found {len(turns)} turns")

new_memories = []
now = datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%SZ")

for i, turn in enumerate(turns):
    # Parse turn header
    header_match = re.match(r"## Turn (\d+) \| (.+?) \| (.+)", turn)
    if not header_match:
        continue
    turn_num = header_match.group(1)
    timestamp = header_match.group(2).strip()
    model = header_match.group(3).strip()

    # Extract Thinking content
    thinking = ""
    think_match = re.search(
        r"\*\*Thinking\*\*: (.+?)(?=\*\*Action\*\*|\*\*Observation\*\*|\n## Turn|\Z)",
        turn,
        re.DOTALL,
    )
    if think_match:
        thinking = think_match.group(1).strip()

    # Extract Observations (tool results)
    obs_match = re.search(
        r"\*\*Observation\*\*: (.+?)(?=\n## Turn|\Z)", turn, re.DOTALL
    )
    observation = ""
    if obs_match:
        observation = obs_match.group(1).strip()

    # Extract Actions
    action_match = re.search(
        r"\*\*Action\*\*: (.+?)(?=\*\*Observation\*\*|\n## Turn|\Z)", turn, re.DOTALL
    )
    action = ""
    if action_match:
        action = action_match.group(1).strip()

    # Determine memory type based on content
    combined = (thinking + " " + observation).lower()

    if "error" in combined or "fail" in combined:
        mem_type = "error"
    elif "success" in combined or "running" in combined or "healthy" in combined:
        mem_type = "system_state"
    elif "installed" in combined or "build" in combined or "deploy" in combined:
        mem_type = "operation"
    elif "memory" in combined or "remember" in combined or "learned" in combined:
        mem_type = "knowledge"
    else:
        mem_type = "long_term"

    # Skip if too short or empty
    if len(thinking) < 30:
        continue

    # Skip if very similar to existing (simple dedup by first 80 chars)
    content_preview = (thinking + " " + observation)[:200]
    is_dup = False
    for m in existing + new_memories:
        if m.get("content", "")[:200] == content_preview:
            is_dup = True
            break
    if is_dup:
        continue

    # Create memory entry
    mem_id = f"brain_{turn_num}_{timestamp.replace(':', '').replace('-', '')}"
    mem = {
        "id": mem_id,
        "type": mem_type,
        "namespace": "brain_log",
        "content": content_preview[:1000],
        "metadata": {
            "source": "brain_log",
            "model": model,
            "turn": turn_num,
            "timestamp": timestamp,
            "observation_preview": observation[:200] if observation else "",
        },
        "createdAt": timestamp,
        "accessedAt": now,
        "accessCount": 0,
    }
    new_memories.append(mem)

    if i > 0 and i % 500 == 0:
        print(
            f"  Processed {i}/{len(turns)} turns, extracted {len(new_memories)} new memories"
        )

print(f"Extracted {len(new_memories)} new memories from Brain log")

# 3. Extract from Gemini antigravity sessions
print("\n--- Gemini antigravity sessions ---")
antigravity_dir = r"C:/Users/hyper/.gemini/antigravity/brain"
if os.path.exists(antigravity_dir):
    sessions = [
        d
        for d in os.listdir(antigravity_dir)
        if os.path.isdir(os.path.join(antigravity_dir, d))
    ]
    print(f"Found {len(sessions)} sessions")
    for session_id in sessions:
        session_dir = os.path.join(antigravity_dir, session_id)
        for filename in ["walkthrough.md", "implementation_plan.md", "task.md"]:
            filepath = os.path.join(session_dir, filename)
            if os.path.exists(filepath):
                with open(filepath, encoding="utf-8", errors="replace") as f:
                    content = f.read()
                if len(content) > 50:
                    mem_id = (
                        f"antigravity_{session_id[:8]}_{filename.replace('.md', '')}"
                    )
                    mem = {
                        "id": mem_id,
                        "type": "knowledge",
                        "namespace": "antigravity",
                        "content": content[:1000],
                        "metadata": {
                            "source": "antigravity",
                            "session": session_id,
                            "file": filename,
                        },
                        "createdAt": now,
                        "accessedAt": now,
                        "accessCount": 0,
                    }
                    new_memories.append(mem)
    print(
        f"Extracted {sum(1 for m in new_memories if m['namespace'] == 'antigravity')} from antigravity"
    )

# 4. Extract from imported sessions
print("\n--- Imported session memories ---")
if os.path.exists(SESSION_DB):
    conn = sqlite3.connect(SESSION_DB)
    cur = conn.execute("SELECT * FROM imported_session_memories")
    cols = [d[0] for d in cur.description]
    rows = cur.fetchall()
    print(f"Found {len(rows)} session memories")
    for row in rows:
        row_dict = dict(zip(cols, row))
        mem_id = row_dict.get("uuid", f"session_{row_dict.get('memory_index', '')}")
        if mem_id not in existing_ids:
            content = row_dict.get("content", "")
            if len(content) > 20:
                new_memories.append(
                    {
                        "id": mem_id,
                        "type": row_dict.get("kind", "long_term"),
                        "namespace": "imported_session",
                        "content": content[:1000],
                        "metadata": json.loads(row_dict.get("metadata", "{}")),
                        "createdAt": row_dict.get("created_at", now),
                        "accessedAt": now,
                        "accessCount": 0,
                    }
                )
    conn.close()
    print(
        f"  Added {sum(1 for m in new_memories if m['namespace'] == 'imported_session')} from sessions"
    )

# 4. Merge into L2
all_memories = existing + new_memories
print("\n--- Merge ---")
print(f"Existing: {len(existing)}")
print(f"New: {len(new_memories)}")
print(f"Total: {len(all_memories)}")

# Deduplicate by ID
seen = set()
deduped = []
for m in all_memories:
    mid = m.get("id")
    if mid and mid not in seen:
        seen.add(mid)
        deduped.append(m)
print(f"After dedupe: {len(deduped)}")

# Write updated L2
output = {"version": 1, "savedAt": now, "memories": deduped}
with open(MEMORIES_PATH, "w", encoding="utf-8") as f:
    json.dump(output, f, indent=2, ensure_ascii=False)
print(f"Written to {MEMORIES_PATH}")

# 5. Append to L3 cold archive
print("\n--- Appending to L3 ---")
conn = sqlite3.connect(L3_DB_PATH)
cur = conn.execute("SELECT id FROM l3_cold_archive")
existing_l3 = set(r[0] for r in cur.fetchall())

batch = []
for m in new_memories:
    mid = m.get("id", "")
    if mid not in existing_l3:
        batch.append(
            (
                mid,
                "",
                m.get("type", "long_term"),
                m.get("namespace", "brain_log"),
                json.dumps([m.get("type", "memory")]),
                "",
                m.get("content", ""),
                0.5,
                0.0,
                now,
                m.get("createdAt", now),
                json.dumps(m.get("metadata", {})),
            )
        )

if batch:
    cur.executemany(
        """
        INSERT OR IGNORE INTO l3_cold_archive
        (id, session_id, kind, category, tags, source_url, content, importance, heat_score, archived_at, created_at, metadata)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    """,
        batch,
    )
    conn.commit()
    print(f"Added {len(batch)} new entries to L3")
else:
    print("No new entries for L3")

# Verify L3 FTS
conn.execute(
    "INSERT OR REPLACE INTO l3_cold_archive_fts(rowid, content) SELECT rowid, content FROM l3_cold_archive"
)
conn.commit()
print("FTS index rebuilt")

cur = conn.execute("SELECT COUNT(*) FROM l3_cold_archive")
print(f"L3 total: {cur.fetchone()[0]} entries")
conn.close()

print("\n" + "=" * 60)
print("Done!")
print(f"L2: {len(deduped)} memories")
print("L3: populated")
print("=" * 60)
