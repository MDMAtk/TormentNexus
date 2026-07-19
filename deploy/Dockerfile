# syntax=docker/dockerfile:1.7

# Stage 1: Builder
FROM node:20-alpine AS builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN apk add --no-cache python3 make g++
RUN corepack enable

WORKDIR /app

# Copy root config
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json turbo.json ./

# Copy package manifests for the minimal runtime workspace graph.
COPY packages/adk/package.json packages/adk/
COPY packages/agents/package.json packages/agents/
COPY packages/ai/package.json packages/ai/
COPY packages/enterprise/package.json packages/enterprise/
COPY packages/tormentnexus/package.json packages/tormentnexus/
COPY packages/mcp-registry/package.json packages/mcp-registry/
COPY packages/memory/package.json packages/memory/
COPY packages/search/package.json packages/search/
COPY packages/tools/package.json packages/tools/
COPY packages/tsconfig/package.json packages/tsconfig/
COPY packages/types/package.json packages/types/
COPY packages/ui/package.json packages/ui/
COPY apps/web/package.json apps/web/
COPY apps/tormentnexus-extension/package.json apps/tormentnexus-extension/
# Install dependencies (frozen lockfile)
RUN pnpm config set store-dir /pnpm/store \
 && pnpm config set fetch-retries 5 \
 && pnpm config set fetch-retry-mintimeout 10000 \
 && pnpm config set fetch-retry-maxtimeout 120000 \
 && pnpm config set network-concurrency 8

RUN --mount=type=cache,id=tormentnexus-pnpm-store-v2,target=/pnpm/store \
	pnpm install --frozen-lockfile --ignore-scripts

# Copy source code required to build the Core API and Web dashboard.
COPY packages/adk packages/adk
COPY packages/agents packages/agents
COPY packages/ai packages/ai
COPY packages/enterprise packages/enterprise
COPY packages/tormentnexus packages/tormentnexus
COPY packages/mcp-registry packages/mcp-registry
COPY packages/memory packages/memory
COPY packages/search packages/search
COPY packages/tools packages/tools
COPY packages/tsconfig packages/tsconfig
COPY packages/types packages/types
COPY packages/ui packages/ui
COPY apps/web apps/web
COPY apps/tormentnexus-extension apps/tormentnexus-extension
COPY scripts scripts

# Ignore any stale incremental TypeScript caches from the host checkout.
RUN find packages apps -name "*.tsbuildinfo" -delete

# Build workspace packages
RUN pnpm run build:workspace

# Stage 2: Core Runner
FROM node:20-alpine AS core
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app

# Copy built artifacts from builder
COPY --from=builder /app ./

# Expose Core port
EXPOSE 3000

# Start Core
CMD ["node", "packages/tormentnexus/dist/index.js"]

# Stage 3: Web Runner
FROM node:20-alpine AS web
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app

# Copy built artifacts from builder
COPY --from=builder /app ./

# Expose the internal Next.js listener port.
EXPOSE 3000

# Start Web
WORKDIR /app/apps/web
CMD ["pnpm", "start"]

# Stage 4: Go Sidecar Builder
FROM golang:1.25-alpine AS sidecar-builder
RUN apk add --no-cache git
WORKDIR /app
COPY go/go.mod go/go.sum ./go/
WORKDIR /app/go
RUN go mod download
WORKDIR /app
COPY go/ go/
WORKDIR /app/go
RUN go build -buildvcs=false -o /app/tormentnexus ./cmd/tormentnexus

# Stage 5: Go Sidecar Runner
FROM alpine:latest AS sidecar
RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /app
COPY --from=sidecar-builder /app/tormentnexus /app/tormentnexus
EXPOSE 7778
CMD ["/app/tormentnexus", "serve", "-port", "7778", "-host", "0.0.0.0"]
