# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22.4-alpine AS builder

# Install certificates for HTTPS requests
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.version=${VERSION:-dev}" \
    -trimpath \
    -o server \
    cmd/main.go

# Runtime stage - scratch for minimal image
FROM scratch

# Copy timezone data and certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create non-root user
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy binary
COPY --from=builder /build/server /server

# Use non-root user (nobody:nobody)
USER 65534:65534

# Expose gRPC port
EXPOSE 50051

# Health check (note: scratch doesn't have shell, so this is informational)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/server", "-health"]

# Run the server
ENTRYPOINT ["/server"]
