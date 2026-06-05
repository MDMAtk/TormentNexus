import os
import re
import hashlib
import json
import yaml

ROOT_DIR = r"c:\Users\hyper\workspace"
TARGET_DIR = r"c:\Users\hyper\workspace\borg\.tormentnexus\skills"
MAIN_SKILL_MD = r"c:\Users\hyper\workspace\SKILL.md"

TARGET_DIRS = [
    ".agent",
    ".borg",
    "borg",
    "skillzhub",
    "mcp-superassistant",
    "onetool-mcp",
    "openclaw-dashboard",
    "opencode-autopilot",
    "superdawmcp",
    "tormentnexus",
]

IGNORED_DIRS = {
    "node_modules", ".venv", ".git", ".ruff_cache", ".pytest_cache", 
    "site-packages", "__pycache__", ".next", ".turbo"
}

def clean_content(text):
    """Normalize whitespace and lowercase to compare content core logic."""
    # Remove markdown titles and frontmatter
    text = re.sub(r"^---[\s\S]*?---", "", text)
    text = re.sub(r"^#\s+.*", "", text, flags=re.MULTILINE)
    # Remove spacing and normalize
    return "".join(text.split()).lower()

def extract_frontmatter(file_path):
    try:
        with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
            content = f.read()
            
        m = re.match(r"^---[\s\S]*?---", content)
        if m:
            fm_text = m.group(0).strip("-\n ")
            fm = yaml.safe_load(fm_text) or {}
            body = content[m.end():].strip()
            return fm, body, content
        else:
            # Try to parse title as name
            title_match = re.search(r"^#\s+(.*)", content, re.MULTILINE)
            name = title_match.group(1).strip() if title_match else os.path.basename(os.path.dirname(file_path))
            return {"name": name, "description": ""}, content.strip(), content
    except Exception as e:
        return {}, "", ""

def parse_main_skill_md():
    skills = []
    if not os.path.exists(MAIN_SKILL_MD):
        return skills
        
    try:
        with open(MAIN_SKILL_MD, "r", encoding="utf-8", errors="ignore") as f:
            content = f.read()
            
        pattern = r'(^##\s+(\d+)\.\s+([^\n]+)$)'
        matches = list(re.finditer(pattern, content, re.MULTILINE))
        
        for idx, match in enumerate(matches):
            name = match.group(3).strip()
            start_pos = match.end()
            end_pos = matches[idx+1].start() if idx + 1 < len(matches) else len(content)
            section_content = content[start_pos:end_pos].strip()
            section_content = re.sub(r'\n---+\n*$', '', section_content).strip()
            
            clean_name = name.lower().replace(" ", "_").replace("-", "_")
            skills.append({
                "name": name,
                "id": clean_name,
                "description": f"Main consolidated skill: {name}",
                "body": section_content,
                "raw": f"# {name}\n\n{section_content}",
                "source": MAIN_SKILL_MD
            })
    except Exception as e:
        print(f"Error parsing main SKILL.md: {e}")
    return skills

def main():
    print("=== DEDUPLICATING AND CATALOGING ALL SYSTEM-WIDE SKILLS ===")
    
    all_skills = []
    
    # 1. Parse main SKILL.md
    print("Parsing main consolidated SKILL.md...")
    main_skills = parse_main_skill_md()
    all_skills.extend(main_skills)
    print(f"Added {len(main_skills)} skills from main SKILL.md")
    
    # 2. Gather skill files dynamically
    found_files = []
    for file in os.listdir(ROOT_DIR):
        if file.lower() in ("skill.md", "skill.markdown"):
            if os.path.abspath(os.path.join(ROOT_DIR, file)) != os.path.abspath(MAIN_SKILL_MD):
                found_files.append(os.path.join(ROOT_DIR, file))
            
    for sub in TARGET_DIRS:
        sub_path = os.path.join(ROOT_DIR, sub)
        if not os.path.exists(sub_path):
            continue
        for root, dirs, files in os.walk(sub_path):
            dirs[:] = [d for d in dirs if d not in IGNORED_DIRS]
            for file in files:
                if file.lower() in ("skill.md", "skill.markdown"):
                    full_path = os.path.join(root, file)
                    if os.path.abspath(full_path) != os.path.abspath(MAIN_SKILL_MD):
                        found_files.append(full_path)
                        
    print(f"Discovered {len(found_files)} raw individual skill files to scrape.")
    
    # 3. Parse each file
    for idx, fpath in enumerate(found_files):
        if idx % 2000 == 0 and idx > 0:
            print(f"Processed {idx} files...")
            
        fm, body, raw = extract_frontmatter(fpath)
        if not fm:
            continue
            
        name = fm.get("name", "") or os.path.basename(os.path.dirname(fpath))
        desc = fm.get("description", "")
        
        clean_id = name.lower().replace(" ", "_").replace("-", "_")
        
        all_skills.append({
            "name": name,
            "id": clean_id,
            "description": desc,
            "body": body,
            "raw": raw,
            "source": fpath
        })
        
    print(f"Found total {len(all_skills)} raw skill definitions across all sources.")
    
    # 4. Deduplicate by content hash
    deduped = {}
    duplicates_log = []
    
    for s in all_skills:
        body_cleaned = clean_content(s["body"])
        if not body_cleaned:
            continue
            
        content_hash = hashlib.sha256(body_cleaned.encode('utf-8')).hexdigest()
        
        if content_hash in deduped:
            existing = deduped[content_hash]
            duplicates_log.append({
                "name": s["name"],
                "source": s["source"],
                "duplicates_with": existing["source"]
            })
            if s["source"] not in existing["sources"]:
                existing["sources"].append(s["source"])
        else:
            s["sources"] = [s["source"]]
            s["content_hash"] = content_hash
            deduped[content_hash] = s

    print(f"\nDeduplication complete:")
    print(f"  Unique Skills: {len(deduped)}")
    print(f"  Duplicate Skills Found: {len(duplicates_log)}")
    
    # 5. Write unique skills to target directory
    print(f"\nWriting unique skills to target directory: {TARGET_DIR}")
    os.makedirs(TARGET_DIR, exist_ok=True)
    
    catalog_index = []
    
    for s in deduped.values():
        skill_id = s["id"]
        safe_id = re.sub(r'[^a-zA-Z0-9_-]', '', skill_id).strip().lower()
        if not safe_id:
            safe_id = "unknown_skill_" + hashlib.md5(s["name"].encode()).hexdigest()[:6]
            
        skill_folder = os.path.join(TARGET_DIR, safe_id)
        os.makedirs(skill_folder, exist_ok=True)
        
        fm_data = {
            "name": s["name"],
            "description": s["description"] or f"Structured runbook for {s['name']}",
            "sources": s["sources"]
        }
        
        yaml_str = yaml.dump(fm_data, sort_keys=False, default_flow_style=False).strip()
        full_markdown = f"---\n{yaml_str}\n---\n\n{s['body']}"
        
        with open(os.path.join(skill_folder, "SKILL.md"), "w", encoding="utf-8") as f:
            f.write(full_markdown)
            
        catalog_index.append({
            "id": safe_id,
            "name": s["name"],
            "description": s["description"],
            "sources": s["sources"],
            "content_hash": s["content_hash"]
        })
        
    # Read existing index if there is one to merge external links
    existing_external_links = []
    report_path = os.path.join(TARGET_DIR, "catalog_index.json")
    if os.path.exists(report_path):
        try:
            with open(report_path, "r", encoding="utf-8") as f:
                old_data = json.load(f)
                existing_external_links = old_data.get("external_links", [])
        except Exception:
            pass
            
    # Write catalog report index including the external links
    with open(report_path, "w", encoding="utf-8") as f:
        json.dump({
            "total_raw_found": len(all_skills),
            "total_unique_local": len(deduped),
            "total_external_links": len(existing_external_links),
            "skills": catalog_index,
            "external_links": existing_external_links,
            "duplicates": duplicates_log
        }, f, indent=2)
        
    print(f"Catalog index written to: {report_path}")
    print("Catalog and deduplication complete!")

if __name__ == "__main__":
    main()
