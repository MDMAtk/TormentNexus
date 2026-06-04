import os

ROOT_DIR = r"c:\Users\hyper\workspace"
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

def main():
    print("Starting fast targeted skill file search...")
    found_files = []
    
    # Also check the root directory itself (non-recursively)
    for file in os.listdir(ROOT_DIR):
        if file.lower() in ("skill.md", "skill.markdown"):
            found_files.append(os.path.join(ROOT_DIR, file))
            
    for sub in TARGET_DIRS:
        sub_path = os.path.join(ROOT_DIR, sub)
        if not os.path.exists(sub_path):
            continue
        print(f"Scanning: {sub_path}")
        for root, dirs, files in os.walk(sub_path):
            dirs[:] = [d for d in dirs if d not in IGNORED_DIRS]
            for file in files:
                if file.lower() in ("skill.md", "skill.markdown"):
                    full_path = os.path.join(root, file)
                    found_files.append(full_path)
                    
    print(f"\nFound {len(found_files)} skill files:")
    for f in found_files:
        print("-", f)

if __name__ == "__main__":
    main()
