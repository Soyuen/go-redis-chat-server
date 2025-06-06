# -----------------------------------------
# Build stage: Compiles the Go application
# -----------------------------------------
FROM golang:1.21.13 AS builder

# Set the working directory inside the container
WORKDIR /app

# Install git (needed for fetching Go modules)
RUN apt-get update && apt-get install -y git

# Copy go.mod and go.sum to take advantage of caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application binary
RUN go build -o server ./cmd/main.go

# -----------------------------------------
# Runtime stage: Runs the compiled binary
# -----------------------------------------
FROM alpine:latest

# Set the working directory in the runtime container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server .

# Set environment variable (can also be set via docker-compose)
ENV PORT=8000

# Expose the application port
EXPOSE 8000

# Run the binary
CMD ["./server"]
