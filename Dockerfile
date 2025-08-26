# ---------- Build Stage ----------
FROM golang:1.24.3-alpine AS builder

# Install git (needed for go get)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN go build -o server ./cmd/server

# ---------- Run Stage ----------
FROM alpine:3.18

WORKDIR /app

# Add CA certificates (for HTTPS requests if needed)
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/server .

# Copy env file (optional – better to mount in prod)
# COPY .env . 

EXPOSE 8080

CMD ["./server"]
