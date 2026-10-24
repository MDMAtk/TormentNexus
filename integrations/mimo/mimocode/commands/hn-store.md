# /hn-store — Store a memory fact

Store an important fact, decision, or finding into HyperNexus persistent memory.

Usage: `/hn-store <title> | <content> [| tag1,tag2]`

Examples:
- `/hn-store DB path | catalog.db was renamed to hypernexus.db in alpha.251 | bugfix,db`
- `/hn-store Deploy target | Production runs on Hetzner at 5.161.250.43 (ssh alias: hetzner) | deploy,infra`

After storing, confirm with the returned fact id if provided.
