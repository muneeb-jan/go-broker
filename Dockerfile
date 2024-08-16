# Stage 1: Build the Go binary
FROM golang:alpine3.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/go-broker/main.go

# Stage 2: Create a minimal container to run the Go binary
FROM alpine:latest

# Install ca-certificates for HTTPS support
RUN apk update && apk add --no-cache ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]

# Expose port (replace with your application's port if different)
EXPOSE 8080
