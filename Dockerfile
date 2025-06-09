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

# Compile Go binary (with static build for Alpine)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# -----------------------------------------
# Runtime stage: Runs the compiled binary
# -----------------------------------------
FROM alpine:latest

# Set the working directory in the runtime container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server .

# Set environment variable (optional)
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./server"]
