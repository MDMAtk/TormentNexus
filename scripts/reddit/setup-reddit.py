#!/usr/bin/env python3
"""
Setup script for TormentNexus Reddit Agent
Creates credentials file and tests connection
"""

import os
import json
import sys


def setup():
    print("=" * 60)
    print("  TormentNexus Reddit Agent Setup")
    print("=" * 60)
    print()

    # Check for PRAW
    try:
        import praw

        print("  [+] PRAW installed")
    except ImportError:
        print("  [!] Installing PRAW...")
        os.system(f"{sys.executable} -m pip install praw requests --quiet")
        print("  [+] PRAW installed")

    # Check for MiMo API key
    mimo_key = os.environ.get("MIMO_API_KEY", "")
    if mimo_key:
        print("  [+] MIMO_API_KEY set")
    else:
        print("  [!] MIMO_API_KEY not set")
        print("      Set it with: export MIMO_API_KEY=your-key")

    print()
    print("  Reddit API Credentials")
    print("  " + "-" * 40)
    print()
    print("  To post to Reddit, you need API credentials:")
    print()
    print("  1. Go to: https://www.reddit.com/prefs/apps")
    print("  2. Click 'Create App' or 'Create Another App'")
    print("  3. Fill in:")
    print("     - Name: TormentNexusBot")
    print("     - Type: script")
    print("     - Redirect URI: http://localhost:8080")
    print("  4. Copy the client_id and client_secret")
    print()

    # Get credentials
    client_id = input("  Client ID: ").strip()
    client_secret = input("  Client Secret: ").strip()
    username = input("  Reddit Username: ").strip()
    password = input("  Reddit Password: ").strip()

    if not all([client_id, client_secret, username, password]):
        print("\n  [!] Missing credentials. Exiting.")
        return

    # Save credentials
    creds_dir = os.path.join(os.path.dirname(__file__), "..", "data")
    os.makedirs(creds_dir, exist_ok=True)
    creds_file = os.path.join(creds_dir, "reddit-creds.json")

    creds = {
        "client_id": client_id,
        "client_secret": client_secret,
        "username": username,
        "password": password,
    }

    with open(creds_file, "w") as f:
        json.dump(creds, f, indent=2)

    print(f"\n  [+] Credentials saved to {creds_file}")

    # Test connection
    print("\n  Testing connection...")
    try:
        import praw

        reddit = praw.Reddit(
            client_id=client_id,
            client_secret=client_secret,
            username=username,
            password=password,
            user_agent="TormentNexus/1.0 (AI research assistant)",
        )
        user = reddit.user.me()
        print(f"  [+] Logged in as: {user}")
        print("\n  Setup complete! Run the agent with:")
        print("  python3 scripts/reddit-agent-v2.py")
    except Exception as e:
        print(f"\n  [!] Connection failed: {e}")
        print("  Check your credentials and try again.")


if __name__ == "__main__":
    setup()
