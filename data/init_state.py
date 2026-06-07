import sqlite3, os
os.makedirs("data", exist_ok=True)
conn = sqlite3.connect("data/assimilation_state.db")
conn.executescript("""
CREATE TABLE IF NOT EXISTS mcp_servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    github_url TEXT DEFAULT '',
    score INTEGER DEFAULT 0,
    status TEXT DEFAULT 'pending',
    go_file TEXT DEFAULT '',
    tools_exposed TEXT DEFAULT '[]',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    category TEXT DEFAULT 'general',
    frontmatter TEXT DEFAULT '',
    content TEXT DEFAULT '',
    version INTEGER DEFAULT 1,
    similarity_score INTEGER DEFAULT 100,
    canonical_id INTEGER REFERENCES skills(id),
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS hermes_addons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    category TEXT DEFAULT 'community',
    description TEXT DEFAULT '',
    go_file TEXT DEFAULT '',
    skill_name TEXT DEFAULT '',
    status TEXT DEFAULT 'pending',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS prompt_library (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    category TEXT DEFAULT 'general',
    content TEXT NOT NULL,
    tags TEXT DEFAULT '[]',
    usage_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_mcp_status   ON mcp_servers(status);
CREATE INDEX IF NOT EXISTS idx_skills_status ON skills(status);
CREATE INDEX IF NOT EXISTS idx_hermes_status ON hermes_addons(status);
""")
conn.commit()
print("State DB ready:", os.path.abspath("data/assimilation_state.db"))
conn.close()
