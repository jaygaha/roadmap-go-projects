# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk update && \
    apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Create user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum .

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app ./cmd/api

# Final stage (reuse already fetched base to avoid extra pull)
FROM golang:1.24-alpine

# Import certs and zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary from builder
COPY --from=builder /build/app /app
RUN chmod +x /app

# Create and use unprivileged user
RUN adduser -D -g '' appuser
USER appuser

# Expose port
EXPOSE 8800

# Run the binary
CMD ["/app"]
