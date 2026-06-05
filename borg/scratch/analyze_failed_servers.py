import sqlite3
import re
from collections import Counter

DB_PATH = r"c:\Users\hyper\workspace\borg\catalog.db"

def analyze_failures():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()
    
    # Query failed validation runs
    c.execute("""
        SELECT r.uuid, s.display_name, s.canonical_id, r.failure_class, r.findings_summary, rc.template
        FROM published_mcp_validation_runs r
        JOIN published_mcp_servers s ON r.server_uuid = s.uuid
        LEFT JOIN published_mcp_config_recipes rc ON s.uuid = rc.server_uuid
        WHERE r.outcome = 'failure'
        ORDER BY r.created_at DESC
        LIMIT 500
    """)
    
    runs = c.fetchall()
    conn.close()
    
    print(f"Analyzing last {len(runs)} failed validation runs...")
    
    patterns = [
        ("Missing Env Var", r"Missing required environment variable[s]?:\s*([A-Za-z0-9_,\s]+)"),
        ("Module Not Found (JS)", r"Cannot find module '([^']+)'"),
        ("File Not Found (Python/General)", r"can't open file '([^']+)': \[Errno 2\] No such file or directory"),
        ("File Not Found (FastMCP/Rust)", r"File not found:.*"),
        ("NPM Package 404", r"404 Not Found - GET .* - Not found"),
        ("NPM Package 404 Registry", r"The requested resource '([^']+)' could not be found"),
        ("Command Not Found", r"'([^']+)' is not recognized as an internal or external command"),
        ("Timeout", r"Connection timeout"),
        ("EADDRINUSE (Port Blocked)", r"address already in use"),
        ("Program Not Found", r"program not found"),
    ]
    
    categorized_failures = []
    uncategorized = []
    
    missing_env_vars = Counter()
    missing_modules = Counter()
    missing_files = Counter()
    missing_npm_packages = Counter()
    commands_missing = Counter()
    
    for r_uuid, name, cid, fail_class, findings, recipe_tpl in runs:
        findings_clean = findings.replace("\n", " ").strip()
        matched = False
        
        for pat_name, pat_regex in patterns:
            m = re.search(pat_regex, findings, re.IGNORECASE)
            if m:
                matched = True
                extracted = m.group(1) if m.groups() else ""
                categorized_failures.append((name, cid, pat_name, extracted, recipe_tpl, findings_clean))
                
                # Update specific counters
                if pat_name == "Missing Env Var":
                    # Split comma-separated vars
                    for v in re.split(r'[,\s]+', extracted):
                        v_clean = v.strip().replace(":", "")
                        if v_clean:
                            missing_env_vars[v_clean] += 1
                elif pat_name == "Module Not Found (JS)":
                    missing_modules[extracted] += 1
                elif pat_name == "File Not Found (Python/General)":
                    missing_files[extracted] += 1
                elif pat_name == "NPM Package 404 Registry":
                    missing_npm_packages[extracted] += 1
                elif pat_name == "Command Not Found":
                    commands_missing[extracted] += 1
                break
                
        if not matched:
            uncategorized.append((name, cid, findings_clean, recipe_tpl))
            
    print(f"\nCategorized Failures: {len(categorized_failures)}")
    print(f"Uncategorized Failures: {len(uncategorized)}\n")
    
    print("=== TOP MISSING ENV VARIABLES ===")
    for var, count in missing_env_vars.most_common(15):
        print(f"  {var}: {count} times")
        
    print("\n=== TOP MISSING JS MODULES ===")
    for mod, count in missing_modules.most_common(10):
        print(f"  {mod}: {count} times")
        
    print("\n=== TOP NPM 404 PACKAGES ===")
    for pkg, count in missing_npm_packages.most_common(10):
        print(f"  {pkg}: {count} times")
        
    print("\n=== TOP MISSING COMMANDS / EXECUTABLES ===")
    for cmd, count in commands_missing.most_common(10):
        print(f"  {cmd}: {count} times")

    print("\n=== TOP MISSING PATHS/FILES ===")
    for file, count in missing_files.most_common(10):
        print(f"  {file}: {count} times")

if __name__ == "__main__":
    analyze_failures()
