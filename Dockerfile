# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (required for go mod) and build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o url-shortener ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/url-shortener .

# Expose the port (matches docker-compose)
EXPOSE 8080

# Command to run the binary
CMD ["./url-shortener"]