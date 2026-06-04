import os

ROOT_DIR = r"c:\Users\hyper\workspace"
IGNORED_DIRS = {
    "node_modules", ".venv", ".git", ".ruff_cache", ".pytest_cache", 
    "site-packages", "__pycache__", ".next", ".turbo"
}

def main():
    print("Searching for SKILL.md and skill.markdown files across workspace...")
    found_files = []
    for root, dirs, files in os.walk(ROOT_DIR):
        # Exclude ignored directories in-place
        dirs[:] = [d for d in dirs if d not in IGNORED_DIRS]
        for file in files:
            if file.lower() in ("skill.md", "skill.markdown"):
                full_path = os.path.join(root, file)
                found_files.append(full_path)
                
    print(f"Found {len(found_files)} files:")
    for f in found_files:
        print("-", f)

if __name__ == "__main__":
    main()
