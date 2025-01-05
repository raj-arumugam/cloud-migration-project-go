# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build
RUN go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]