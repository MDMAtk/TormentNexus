# Discord Server Setup Guide

## Quick Setup (5 minutes)

### Step 1: Create Server

1. Open Discord
2. Click "+" to create a server
3. Choose "Create My Own"
4. Name it "TormentNexus"
5. Upload the logo (see `assets/logo.png`)

### Step 2: Create Channels

```
📋 INFORMATION
├── #announcements (read-only)
├── #rules (read-only)
├── #welcome
└── #faq (read-only)

💬 GENERAL
├── #general
├── #introductions
└── #off-topic

🛠️ DEVELOPMENT
├── #help
├── #bug-reports
├── #feature-requests
└── #pull-requests

🤖 MCP & TOOLS
├── #mcp-servers
├── #tool-showcase
└── #tool-requests

🧠 MEMORY & AI
├── #memory-system
├── #local-llms
└── #prompts

📣 UPDATES
├── #releases (webhook from GitHub)
└── #blog-posts (webhook from site)
```

### Step 3: Set Permissions

- @everyone: Read, Send Messages, Add Reactions
- @moderators: Manage Messages, Kick Members
- @admins: Administrator

### Step 4: Add Bots

- **GitHub Bot**: For release notifications
- **MEE6**: For moderation and leveling
- **Carl-bot**: For reaction roles

### Step 5: Set Up Webhooks

1. Go to Server Settings → Integrations
2. Create webhook for #releases:
   - URL: GitHub webhook
   - Events: Releases
3. Create webhook for #blog-posts:
   - URL: RSS feed webhook

### Step 6: Create Invite Link

1. Click "Invite People"
2. Set to "Never Expire"
3. Copy link: `https://discord.gg/tormentnexus`

### Step 7: Add to Landing Page

Update the landing page with the Discord link.

---

## Server Template

```json
{
  "name": "TormentNexus",
  "icon": "assets/logo.png",
  "channels": [
    { "name": "announcements", "type": "text", "read_only": true },
    { "name": "rules", "type": "text", "read_only": true },
    { "name": "welcome", "type": "text" },
    { "name": "faq", "type": "text", "read_only": true },
    { "name": "general", "type": "text" },
    { "name": "introductions", "type": "text" },
    { "name": "off-topic", "type": "text" },
    { "name": "help", "type": "text" },
    { "name": "bug-reports", "type": "text" },
    { "name": "feature-requests", "type": "text" },
    { "name": "pull-requests", "type": "text" },
    { "name": "mcp-servers", "type": "text" },
    { "name": "tool-showcase", "type": "text" },
    { "name": "tool-requests", "type": "text" },
    { "name": "memory-system", "type": "text" },
    { "name": "local-llms", "type": "text" },
    { "name": "prompts", "type": "text" },
    { "name": "releases", "type": "text" },
    { "name": "blog-posts", "type": "text" }
  ]
}
```
