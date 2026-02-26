# Multi-stage Docker build for mcp-server-planton.
#
# Local build (with private module access):
#   docker build --secret id=gh_pat,env=GH_PAT -t mcp-server-planton .
#
# CI builds pass the secret via docker/build-push-action's "secrets" input.

# ---- Build stage ----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum ./
RUN --mount=type=secret,id=gh_pat \
    git config --global url."https://$(cat /run/secrets/gh_pat)@github.com/".insteadOf "https://github.com/" && \
    GOPRIVATE=github.com/plantonhq/* go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /mcp-server-planton \
    ./cmd/mcp-server-planton

# ---- Runtime stage ----
FROM alpine:3.19

# CA certificates are needed for TLS connections to the Planton backend.
RUN apk --no-cache add ca-certificates

# Run as a non-root user.
RUN addgroup -g 1000 planton && \
    adduser -D -u 1000 -G planton planton

WORKDIR /app
COPY --from=builder /mcp-server-planton .
RUN chown -R planton:planton /app

USER planton

# Default HTTP port (overridable via PLANTON_MCP_HTTP_PORT).
EXPOSE 8080

# Health check for container orchestrators.
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./mcp-server-planton"]
