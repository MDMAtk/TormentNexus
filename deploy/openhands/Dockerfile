# OpenHands + TormentNexus Combined Docker Image
# Build: docker build -t tormentnexus-openhands:latest -f deploy/openhands/Dockerfile .
# Run:   docker compose -f deploy/openhands/docker-compose.yml up

FROM ghcr.io/all-hands-ai/openhands:latest AS openhands-base
FROM tormentnexus-core:latest AS tn-base

FROM openhands-base

# Copy TormentNexus binary and config from TN base
COPY --from=tn-base /app/tormentnexus /usr/local/bin/tormentnexus
COPY --from=tn-base /app/tormentnexus-core /usr/local/bin/tormentnexus-core

# Install TN OpenHands plugin
COPY deploy/openhands/plugins/ /home/openhands/.openhands/plugins/
COPY npm/@tormentnexus/openhands/ /opt/tormentnexus-openhands/

# Pre-configure TN MCP server
RUN mkdir -p /home/openhands/.openhands/microagents && \
    echo '{"mcpServers":{"tormentnexus":{"command":"tormentnexus","args":["mcp"],"env":{"TORMENTNEXUS_WORKSPACE_ROOT":"/opt/workspace_base"},"type":"stdio","lifecycle":"eager"}}}' \
    > /home/openhands/.openhands/mcp.json

# Copy microagent
COPY npm/@tormentnexus/openhands/microagent.md /home/openhands/.openhands/microagents/tormentnexus.md

# TN data volume
VOLUME ["/data/tormentnexus"]

EXPOSE 7778 3000
