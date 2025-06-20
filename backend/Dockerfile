# backend/Dockerfile
# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies and leverage Docker cache.
# The go.sum\* wildcard handles the case where go.sum doesn't exist yet.
COPY go.mod go.sum* ./

# Copy the rest of the source code BEFORE running go mod tidy
# This ensures go mod tidy sees all import statements
COPY . .

# Tidy ensures go.sum is correct and all dependencies are downloaded.
# This is more robust than 'go mod download'.
RUN go mod tidy

# Build the application
# CGO_ENABLED=0 is important for creating a static binary for alpine
# -o /app/server creates the binary named 'server' in the /app directory
# This is the correct, idiomatic way for Dockerfiles
ENV CGO_ENABLED=0
RUN go build -o /app/server .

# Stage 2: Create a small production image
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/server .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"]