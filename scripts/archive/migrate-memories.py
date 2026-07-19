#!/usr/bin/env python3
"""
Migrate memories from workspace .tormentnexus to home-dir .tormentnexus.
Copies the 4K L2 memories to the correct path and populates L3 cold archive.
"""

import json
import sqlite3
import shutil
import os
from datetime import datetime

WORKSPACE_TN = r"C:\Users\hyper\workspace\tormentnexus\.tormentnexus"
HOME_TN = r"C:\Users\hyper\.tormentnexus"

print("=" * 60)
print("Memory Migration Tool")
print("=" * 60)

# 1. Copy agent_memory/memories.json to home dir
src = os.path.join(WORKSPACE_TN, "agent_memory", "memories.json")
dst = os.path.join(HOME_TN, "agent_memory", "memories.json")

if os.path.exists(src):
    os.makedirs(os.path.dirname(dst), exist_ok=True)
    shutil.copy2(src, dst)
    src_size = os.path.getsize(src)
    print(f"✅ Copied agent_memory: {src_size:,} bytes -> {dst}")
else:
    print(f"⚠️  Source not found: {src}")

# 2. Read memories
with open(src, "r", encoding="utf-8") as f:
    data = json.load(f)

memories = data.get("memories", [])
print(f"📦 Loaded {len(memories)} memories from agent_memory")

# 3. Populate L3 Cold Archive
l3_db = os.path.join(HOME_TN, "l3_cold_archive.db")
print(f"\n💾 Populating L3 Cold Archive: {l3_db}")

conn = sqlite3.connect(l3_db)
cur = conn.cursor()

# Check existing count
cur.execute("SELECT COUNT(*) FROM l3_cold_archive")
existing = cur.fetchone()[0]
print(f"   Existing L3 entries: {existing}")

if existing == 0:
    # Batch insert all memories
    batch = []
    now = datetime.utcnow().isoformat() + "Z"
    for m in memories:
        batch.append(
            (
                m.get("id", ""),
                m.get("session_id", ""),
                m.get("type", "memory"),
                m.get("type", "general"),
                json.dumps(m.get("tags", [])),
                "",
                m.get("content", ""),
                0.5,
                0.0,
                now,
                m.get("createdAt", now),
                json.dumps(m.get("metadata", {})),
            )
        )

    cur.executemany(
        """
        INSERT OR IGNORE INTO l3_cold_archive
        (id, session_id, kind, category, tags, source_url, content, importance, heat_score, archived_at, created_at, metadata)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    """,
        batch,
    )
    conn.commit()

    cur.execute("SELECT COUNT(*) FROM l3_cold_archive")
    total = cur.fetchone()[0]
    print(f"Inserted {total} memories into L3 cold archive")
else:
    print("   L3 already has data - skipping insert")

# 4. Also update sectioned_memory.json if empty
sectioned_src = os.path.join(WORKSPACE_TN, "sectioned_memory.json")
sectioned_dst = os.path.join(HOME_TN, "sectioned_memory.json")

with open(sectioned_src, "r", encoding="utf-8") as f:
    sectioned = json.load(f)

# Check if sections are empty
empty = all(len(s.get("entries", [])) == 0 for s in sectioned.get("sections", []))
if empty:
    print("\n📝 Sectioned memory is empty — populating from L2 memories")
    for section in sectioned.get("sections", []):
        section_name = section["section"]
        relevant = [
            m
            for m in memories
            if m.get("type") == section_name
            or (
                section_name == "general"
                and m.get("type", "") in ["", "general", "note"]
            )
        ]
        section["entries"] = [
            {
                "id": m.get("id"),
                "content": m.get("content", "")[:500],
                "metadata": m.get("metadata", {}),
                "createdAt": m.get("createdAt"),
            }
            for m in relevant[:50]
        ]  # Limit to 50 per section
        print(f"   {section_name}: {len(section['entries'])} entries")

    with open(sectioned_dst, "w", encoding="utf-8") as f:
        json.dump(sectioned, f, indent=2)
    print("✅ Updated sectioned_memory.json")

conn.close()
print("\n" + "=" * 60)
print("Migration complete!")
print(f"   L2 agent_memory: {len(memories)} memories")
print("   L3 cold archive: populated")
print("=" * 60)
