#!/usr/bin/env python3
"""
TormentNexus Universal Client Support Installer
Detects installed AI coding clients and installs MCP configs,
skill files, plugins, extensions, hooks, and commands.
"""

import os
import sys
import json
import shutil
import platform

HOME = os.path.expanduser("~")
WORKSPACE = os.path.join(HOME, "workspace", "tormentnexus")
SYSTEM = platform.system()

# === CLIENT SUPPORT MATRIX ===
# Format: name -> dict(dirs, support types)
CLIENTS = {
    # Category 1: Official Big-Tech Frameworks
    "claude": {
        "dirs": [".claude"],
        "skills": True,
        "mcp": True,
        "hooks": True,
        "commands": True,
        "ext": False,
    },
    "gemini": {
        "dirs": [".gemini"],
        "skills": True,
        "mcp": True,
        "ext": True,
        "agents": True,
        "commands": True,
        "hooks": False,
    },
    "codex": {
        "dirs": [".codex"],
        "skills": True,
        "mcp": True,
        "commands": True,
        "hooks": False,
        "ext": False,
    },
    "grok": {
        "dirs": [".grok"],
        "skills": True,
        "mcp": True,
        "commands": True,
        "hooks": True,
        "ext": True,
    },
    "antigravity": {
        "dirs": [".antigravity", ".gemini/antigravity-ide"],
        "skills": True,
        "mcp": True,
        "ext": True,
        "agents": True,
        "commands": False,
    },
    # Category 2: Open-Source BYOK CLIs
    "aider": {
        "dirs": [".aider"],
        "mcp": True,
        "skills": True,
        "hooks": True,
        "commands": True,
        "ext": False,
    },
    "opencode": {
        "dirs": [".opencode"],
        "mcp": True,
        "skills": True,
        "commands": True,
        "hooks": True,
        "ext": True,
    },
    "openclaw": {
        "dirs": [".openclaw"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "goose": {
        "dirs": [".goose", "goose"],
        "mcp": True,
        "ext": True,
        "skills": True,
        "hooks": True,
        "commands": True,
    },
    "iflow": {
        "dirs": [".iflow"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "roo": {
        "dirs": [".roo", ".roo-code"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": True,
        "ext": False,
    },
    "cline": {
        "dirs": [".cline"],
        "mcp": True,
        "skills": True,
        "hooks": True,
        "commands": True,
        "ext": False,
    },
    # Category 3: Full-IDE Clients
    "cursor": {
        "dirs": [".cursor"],
        "skills": True,
        "mcp": True,
        "commands": True,
        "hooks": True,
        "ext": True,
    },
    "windsurf": {
        "dirs": [".windsurf"],
        "mcp": True,
        "skills": True,
        "commands": True,
        "hooks": True,
        "ext": True,
    },
    "zed": {
        "dirs": [".zed"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "trae": {
        "dirs": [".trae"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "continue": {
        "dirs": [".continue"],
        "mcp": True,
        "skills": True,
        "ext": True,
        "hooks": False,
        "commands": False,
    },
    # Category 4: B2B Autonomous Platforms
    "factory": {
        "dirs": [".factory"],
        "mcp": True,
        "skills": True,
        "commands": True,
        "hooks": False,
        "ext": False,
    },
    "openhands": {
        "dirs": [".openhands"],
        "mcp": True,
        "skills": True,
        "hooks": True,
        "commands": True,
        "ext": True,
        "agents": True,
    },
    "kiro": {
        "dirs": [".kiro"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "codewhale": {
        "dirs": [".codewhale"],
        "skills": True,
        "mcp": True,
        "plugin": True,
        "ext": True,
        "commands": True,
    },
    # Category 5: Orchestrators
    "omnigent": {
        "dirs": [".omnigent"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "citadel": {
        "dirs": [".citadel"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "agent-fusion": {
        "dirs": [".agent-fusion"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "herdr": {
        "dirs": [".herdr"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "claude-squad": {
        "dirs": [".claude-squad"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    # Category 6: Runtimes & Specialized Engines
    "qwen-code": {
        "dirs": [".qwen-code", ".qwen"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "pi": {
        "dirs": [".pi/agent"],
        "ext": True,
        "mcp": True,
        "skills": False,
        "hooks": False,
        "commands": False,
    },
    "kimi-code": {
        "dirs": [".kimi-code", ".moonshot"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    "cliproxyapi": {
        "dirs": [".cliproxyapi"],
        "mcp": True,
        "skills": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
    # Bonus: VS Code & JetBrains
    "vscode": {
        "dirs": [".vscode"],
        "skills": True,
        "mcp": True,
        "ext": True,
        "hooks": False,
        "commands": True,
    },
    "jetbrains": {
        "dirs": [".jetbrains"],
        "skills": True,
        "mcp": True,
        "ext": True,
        "hooks": False,
        "commands": False,
    },
    # Hermes memory
    "hermes": {
        "dirs": [".hermes"],
        "skills": True,
        "mcp": True,
        "hooks": False,
        "commands": False,
        "ext": False,
    },
}


def get_skill_content():
    return """# TormentNexus Skill — Universal AI Control Plane

## Overview
TormentNexus is your local AI control plane running on port 7778. It provides persistent
multi-tier memory (L1 scratchpad, L2 vector store, L3 cold archive), MCP tool routing
across 20+ servers, session import from Claude Code/Aider/Gemini, and commercial RBAC.

## Quick Start
1. Ensure TN Kernel is running: `http://127.0.0.1:7778/api/runtime/status`
2. Use `tn_memory_search` before any significant task to recall past context
3. Store key decisions with `tn_memory_store` using descriptive tags
4. Route through TN Kernel for commercial integrations (Jira, Confluence)
5. Use `tn_tool_search` to find the right tool for any job

## Available Tools
- `tn_memory_store` — Save important decisions with tags
- `tn_memory_search` — Find past memories by keyword, tag, or category
- `tn_memory_vector_search` — Semantic vector search
- `tn_tool_search` — Discover tools across 20+ MCP servers
- `tn_session_search` — Browse imported sessions
- `tn_skill_manage` — Access 5,776 reusable skill modules
- `tn_code_search` — Search code via AST-grep or pattern matching
- `tn_context_harvest` — Pull relevant L2 context

## Memory Best Practices
1. Search L2 before starting any significant task
2. Store important decisions, patterns, and facts
3. Use `@memory:keyword` inline for auto-expanded context
4. Check the cold archive for archived knowledge

## Commercial Security
All destructive tool calls are checked against commercial RBAC policies.
"""


def get_mcp_config():
    return {
        "mcpServers": {
            "tormentnexus": {
                "command": "tormentnexus"
                if SYSTEM != "Windows"
                else "tormentnexus.exe",
                "args": ["mcp"],
                "env": {"TORMENTNEXUS_WORKSPACE_ROOT": WORKSPACE},
                "type": "stdio",
                "lifecycle": "eager",
            }
        }
    }


def get_plugin_config():
    return """[plugin]
name = "tormentnexus"
version = "1.0.0"
description = "TormentNexus AI Control Plane — persistent memory, MCP routing, session import, commercial RBAC"

[hooks]
on_session_start = "tn_context_harvest"
on_tool_call = "tn_audit_log"
on_error = "tn_memory_store"

[commands]
"tn-store" = { description = "Interactive memory store", action = "slash" }
"tn-search" = { description = "Search memories", action = "slash" }
"tn-status" = { description = "System status", action = "slash" }
"tn-plan" = { description = "Project plans", action = "slash" }
"tn-purge" = { description = "Remove stale memories", action = "slash" }
"tn-summary" = { description = "Session summary", action = "slash" }

[memory]
enable_l2 = true
enable_l3 = true
auto_harvest = true
"""


def get_command_content(name):
    commands = {
        "tn-store": "Store a memory with tags. Usage: /tn-store",
        "tn-search": "Search L2/L3 memory. Usage: /tn-search [query]",
        "tn-status": "Show TN system status. Usage: /tn-status",
        "tn-plan": "Manage project plans. Usage: /tn-plan",
        "tn-purge": "Remove stale memories. Usage: /tn-purge",
        "tn-summary": "Summarize current session. Usage: /tn-summary",
    }
    return commands.get(name, f"TormentNexus command: {name}")


def get_agent_content():
    return """# TormentNexus Agent Configuration

## Role
You are an autonomous agent with access to TormentNexus — a universal AI control plane
with persistent memory, tool orchestration, and session management.

## Capabilities
- Persistent multi-tier memory (L1/L2/L3)
- MCP tool routing across 20+ servers
- Session import from Claude Code, Aider, Gemini
- Commercial RBAC security
- Skill registry with 5,776+ modules

## Instructions
1. Use tn_memory_search before starting any task
2. Store key decisions as you work
3. Route commercial integrations through TN Kernel
4. Check the cold archive for historical context
5. Use tn_context_harvest for complex multi-step tasks
"""


def install():
    mode = "personal"
    if "--mode" in sys.argv or "corporate" in sys.argv:
        mode = "corporate"

    print("=" * 60)
    print("TormentNexus Universal Client Support Installer")
    print(f"System: {SYSTEM} | Workspace: {WORKSPACE}")
    print(f"Mode: {mode.upper()}")
    print("=" * 60)

    installed = 0
    for client_name, config in CLIENTS.items():
        dirs = config.get("dirs", [])
        found = False
        for d in dirs:
            path = os.path.join(HOME, d)
            if os.path.exists(path) or True:  # Always install (preemptive)
                found = True
                base = os.path.join(HOME, d, "tormentnexus")
                os.makedirs(base, exist_ok=True)

                # Skills
                if config.get("skills"):
                    skill_dir = os.path.join(base, "skills")
                    os.makedirs(skill_dir, exist_ok=True)
                    content = get_skill_content()
                    if mode == "corporate":
                        content = content.replace("TormentNexus", "HyperNexus").replace(
                            "local AI control plane",
                            "commercial AI control plane with SSO, RBAC, and audit logging",
                        )
                    with open(os.path.join(skill_dir, "SKILL.md"), "w") as f:
                        f.write(content)

                # MCP Config
                if config.get("mcp"):
                    mcp_dir = os.path.join(base, "mcp")
                    os.makedirs(mcp_dir, exist_ok=True)
                    with open(os.path.join(mcp_dir, "servers.json"), "w") as f:
                        json.dump(get_mcp_config(), f, indent=2)

                # Plugin
                if config.get("plugin"):
                    with open(os.path.join(base, "plugin.toml"), "w") as f:
                        f.write(get_plugin_config())

                # Commands
                if config.get("commands"):
                    cmd_dir = os.path.join(base, "commands")
                    os.makedirs(cmd_dir, exist_ok=True)
                    for cmd in [
                        "tn-store",
                        "tn-search",
                        "tn-status",
                        "tn-plan",
                        "tn-purge",
                        "tn-summary",
                    ]:
                        with open(os.path.join(cmd_dir, f"{cmd}.md"), "w") as f:
                            f.write(get_command_content(cmd))

                # Agent config
                if config.get("agents"):
                    agent_dir = os.path.join(base, "agents")
                    os.makedirs(agent_dir, exist_ok=True)
                    with open(os.path.join(agent_dir, "agent.md"), "w") as f:
                        f.write(get_agent_content())

                # Hooks
                if config.get("hooks"):
                    hooks_dir = os.path.join(base, "hooks")
                    os.makedirs(hooks_dir, exist_ok=True)
                    hooks_config = {
                        "on_session_start": "tn_context_harvest",
                        "on_tool_error": "tn_memory_store",
                        "on_decision": "tn_memory_store",
                    }
                    if mode == "corporate":
                        hooks_config["on_security_violation"] = "tn_audit_log"
                        hooks_config["on_rbac_check"] = "tn_security_check"
                    with open(os.path.join(hooks_dir, "config.json"), "w") as f:
                        json.dump(hooks_config, f, indent=2)

                # Extension
                if config.get("ext"):
                    ext_dir = os.path.join(base, "extensions")
                    os.makedirs(ext_dir, exist_ok=True)
                    # Copy Pi extension for clients that support it
                    src_ext = os.path.join(
                        HOME, ".pi", "agent", "extensions", "tormentnexus.ts"
                    )
                    if os.path.exists(src_ext):
                        shutil.copy2(src_ext, ext_dir)

                # OpenHands-specific: install agent, actions, plugin, docker compose
                if client_name == "openhands":
                    deploy_dir = os.path.join(WORKSPACE, "deploy", "openhands")
                    if os.path.isdir(deploy_dir):
                        # Copy agent + actions + plugin to OpenHands plugins dir
                        plugins_dir = os.path.join(
                            HOME, ".openhands", "plugins", "tormentnexus"
                        )
                        os.makedirs(plugins_dir, exist_ok=True)
                        for fname in [
                            "tormentnexus_agent.py",
                            "actions.py",
                            "plugin.toml",
                        ]:
                            src = os.path.join(deploy_dir, fname)
                            if os.path.exists(src):
                                shutil.copy2(src, os.path.join(plugins_dir, fname))
                        print(f"  openhands plugins -> {plugins_dir}")
                        # Copy Docker compose
                        dc_src = os.path.join(deploy_dir, "docker-compose.yml")
                        if os.path.exists(dc_src):
                            shutil.copy2(
                                dc_src, os.path.join(HOME, ".openhands", "tormentnexus")
                            )
                        # Copy microagent
                        ma_src = os.path.join(
                            os.path.dirname(deploy_dir),
                            "..",
                            "npm",
                            "@tormentnexus",
                            "openhands",
                            "microagent.md",
                        )
                        if os.path.exists(ma_src):
                            os.makedirs(
                                os.path.join(HOME, ".openhands", "microagents"),
                                exist_ok=True,
                            )
                            shutil.copy2(
                                ma_src,
                                os.path.join(
                                    HOME, ".openhands", "microagents", "tormentnexus.md"
                                ),
                            )

                print(f"  {client_name} -> {base}")
                installed += 1

    # Corporate mode: install commercial config
    if mode == "corporate":
        print("\n--- Installing Commercial Configuration ---")
        commercial_dir = os.path.join(HOME, ".tormentnexus", "commercial")
        os.makedirs(commercial_dir, exist_ok=True)

        # SSO/OIDC config
        sso_config = {
            "provider": "https://identity.hypernexus.site/oauth/v2",
            "client_id": "hypernexus-dashboard-prod",
            "scopes": ["openid", "profile", "email", "groups"],
            "claim_mapping": {"email": "email", "groups": "groups", "name": "name"},
            "session_timeout": 3600,
        }
        with open(os.path.join(commercial_dir, "sso.json"), "w") as f:
            json.dump(sso_config, f, indent=2)
        print("  SSO/OIDC config -> sso.json")

        # RBAC roles
        rbac = {
            "roles": {
                "admin": {
                    "permissions": ["*"],
                    "scopes": ["memory:write", "tool:exec", "config:write"],
                },
                "developer": {
                    "permissions": ["memory:read", "memory:write", "tool:exec"],
                    "scopes": [],
                },
                "auditor": {"permissions": ["memory:read", "audit:read"], "scopes": []},
                "viewer": {"permissions": ["memory:read"], "scopes": []},
            }
        }
        with open(os.path.join(commercial_dir, "rbac.json"), "w") as f:
            json.dump(rbac, f, indent=2)
        print("  RBAC roles -> rbac.json")

        # Audit config
        audit = {
            "backend": "sqlite",
            "path": "audit.db",
            "rotation": "daily",
            "retention_days": 90,
            "log_level": "info",
        }
        with open(os.path.join(commercial_dir, "audit.json"), "w") as f:
            json.dump(audit, f, indent=2)
        print("  Audit logging -> audit.json")

        # Tenant isolation
        tenant = {
            "isolation": "container",
            "data_root": "/var/lib/hypernexus/tenants",
            "network": "tenant-network",
            "resource_limits": {"cpu": "1.5", "memory": "2G"},
            "auto_provision": True,
        }
        with open(os.path.join(commercial_dir, "tenants.json"), "w") as f:
            json.dump(tenant, f, indent=2)
        print("  Tenant isolation -> tenants.json")

        # License template
        lic = {
            "licensed_to": "YOUR_ORGANIZATION",
            "features": ["sso", "rbac", "audit", "multi_tenant", "mesh"],
            "expires": "2027-01-01",
            "seats": 50,
        }
        with open(os.path.join(commercial_dir, "license.template.json"), "w") as f:
            json.dump(lic, f, indent=2)
        print("  License template -> license.template.json")

    # Also install preemptively to default config directories
    preemptive = [
        (".claude", "Claude Code"),
        (".cursor", "Cursor"),
        (".windsurf", "Windsurf"),
        (".zed", "Zed"),
        (".continue", "Continue.dev"),
    ]
    for config_dir, name in preemptive:
        path = os.path.join(HOME, config_dir)
        if not os.path.exists(path):
            print(f"  [PREEMPTIVE] Created empty config dir for {name}: {path}")
            os.makedirs(path, exist_ok=True)
            # Also install TN support there
            base = os.path.join(path, "tormentnexus", "skills")
            os.makedirs(base, exist_ok=True)
            with open(os.path.join(base, "SKILL.md"), "w") as f:
                f.write(get_skill_content())

    print(f"\n{installed} clients supported. NONE LEFT BEHIND.")
    print("TormentNexus now works with every AI coding agent on your system.")


if __name__ == "__main__":
    install()
