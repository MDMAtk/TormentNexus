# TormentNexus 5-Minute Demo Script

## Setup (30s)

```bash
docker run -d -p 7778:7778 tormentnexus
curl http://127.0.0.1:7778/api/runtime/status | jq .data.version
```

## Memory (60s)

```bash
tn store "Deployment uses blue-green strategy on Hetzner" --tags deployment
tn search deployment
```

## Tools (60s)

```bash
tn tools
tn status
```

## AI Agent Integration (90s)

```bash
# Install for your favorite agent
python3 scripts/install-client-support.py
# Shows: "38 clients supported. NONE LEFT BEHIND."
```

## Cloud Dashboard (60s)

```
Open: https://hypernexus.site
Click: Sign Up → Stripe → provisioned
Login: https://cloud.hypernexus.site
```

## Summary

- One binary, one port, one API
- Works with every AI coding agent
- Persistent memory across sessions
- Self-host or cloud
