FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go/go.mod go/go.sum ./
RUN go mod download

# Copy source code
COPY go/ .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -ldflags="-s -w" -o hypernexus-cloud ./cmd/cloud

# Runtime image
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata docker-cli

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/hypernexus-cloud .

# Create directories
RUN mkdir -p /var/lib/hypernexus /var/backups/hypernexus /data

# Expose ports
EXPOSE 7778 7779

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:7778/api/cloud/health || exit 1

# Run
ENTRYPOINT ["./hypernexus-cloud"]
CMD ["serve", "--cloud"]
