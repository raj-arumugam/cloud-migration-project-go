# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for downloading dependencies
RUN apk add --no-cache git

# Copy go mod files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cloud-migration

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/cloud-migration .

# Copy any additional configs if needed
COPY --from=builder /app/config.yaml .

RUN chmod +x /app/cloud-migration

EXPOSE 8080

# Add after CMD line
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Add signal handling for graceful shutdown
STOPSIGNAL SIGTERM

CMD ["/app/cloud-migration"]