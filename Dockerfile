FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go/go.mod go/go.sum ./
RUN go mod download

# Copy source code
COPY go/ .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -ldflags="-s -w" -o tormentnexus ./cmd/tormentnexus

# Runtime image
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/tormentnexus .

# Create data directory
RUN mkdir -p /root/.tormentnexus

# Expose port
EXPOSE 7778

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:7778/health || exit 1

# Run
ENTRYPOINT ["./tormentnexus"]
CMD ["serve"]
