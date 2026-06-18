"""
Session Import Automation Script
==================================
Reads candidate session files discovered by the Go sidecar scanner
and imports them into the tormentnexus database via the tRPC bridge.

Usage:
    python scripts/import_sessions.py [--dry-run]
"""

import json
import os
import sys
import urllib.request
import urllib.error

GO_API = "http://localhost:4300"
TS_TRPC = "http://localhost:4100/trpc"


def get_candidates():
    """Fetch candidate list from Go sidecar."""
    url = f"{GO_API}/api/import/summary"
    try:
        with urllib.request.urlopen(url, timeout=10) as r:
            return json.loads(r.read()).get("data", {})
    except Exception as e:
        print(f"ERROR: Cannot reach Go sidecar at {GO_API}: {e}")
        sys.exit(1)


def import_session(data_json: str, dry_run: bool = False):
    """Import a session via TS control plane tRPC bridge."""
    payload = json.dumps(
        {
            "data": data_json,
            "merge": True,
            "dryRun": dry_run,
        }
    ).encode()

    req = urllib.request.Request(
        f"{GO_API}/api/session-export/import",
        data=payload,
        headers={"Content-Type": "application/json"},
    )
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            return json.loads(r.read())
    except urllib.error.HTTPError as e:
        return {"error": e.read().decode()[:200]}
    except Exception as e:
        return {"error": str(e)}


def main():
    dry_run = "--dry-run" in sys.argv

    print("=== Session Import Automation ===\n")
    print(f"Dry run: {dry_run}\n")

    data = get_candidates()
    count = data.get("count", 0)
    valid = data.get("validCount", 0)
    print(f"Candidates found: {count} total, {valid} valid\n")

    if count == 0:
        print("No candidates to import.")
        return

    print(f"By source tool: {data.get('bySourceTool', [])}")
    print(f"By format: {data.get('byFormat', [])}")
    print()

    # The Go sidecar doesn't expose file paths directly through the API.
    # Candidate files are in ~/.claude, ~/.aider, etc.
    # We scan known directories for session files matching the candidate count.

    import glob

    home = os.path.expanduser("~")
    workspace = os.getcwd()

    scan_roots = [
        os.path.join(home, ".claude"),
        os.path.join(home, ".aider"),
        os.path.join(workspace, ".claude"),
        os.path.join(workspace, ".aider"),
        os.path.join(home, "AppData", "Roaming", "Claude"),
    ]

    session_files = []
    for root in scan_roots:
        if os.path.isdir(root):
            for ext in ["*.jsonl", "*.json", "*.md"]:
                session_files.extend(
                    glob.glob(os.path.join(root, "**", ext), recursive=True)
                )

    # Filter to files that look like session exports
    candidates = []
    for f in session_files:
        name = os.path.basename(f).lower()
        if any(
            hint in name
            for hint in ["session", "chat", "conversation", "transcript", "history"]
        ):
            candidates.append(f)

    candidates = sorted(set(candidates))
    print(f"Found {len(candidates)} session files on disk\n")

    imported = 0
    errors = 0
    for fpath in candidates:
        fsize = os.path.getsize(fpath)
        print(
            f"  [{'+' if not dry_run else 'D'}] {os.path.relpath(fpath, home)} ({fsize:,} bytes)"
        )

        if dry_run:
            continue

        try:
            with open(fpath, "r", encoding="utf-8", errors="replace") as f:
                content = f.read()
        except Exception as e:
            print(f"    ERROR reading file: {e}")
            errors += 1
            continue

        if not content.strip():
            print("    SKIP: empty file")
            continue

        result = import_session(content, dry_run=False)
        if result.get("success"):
            imported += 1
            i = result.get("data", {}).get("imported", 0)
            m = result.get("data", {}).get("merged", 0)
            print(f"    OK: imported={i}, merged={m}")
        else:
            errors += 1
            err_msg = result.get("error", "unknown")
            print(f"    FAIL: {err_msg[:100]}")

    print("\n=== Summary ===")
    print(f"  Total candidates: {len(candidates)}")
    print(f"  Imported: {imported}")
    print(f"  Errors: {errors}")
    print(f"  Dry run: {dry_run}")


if __name__ == "__main__":
    main()
