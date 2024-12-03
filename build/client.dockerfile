# Use the official Golang image for building
FROM golang:1.23-alpine as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the client binary
RUN go build -o client cmd/client/main.go

# Use a minimal image for running the client
FROM debian:bullseye-slim

# Copy the client binary
COPY --from=builder /app/client /client

# Expose the default port (if needed for debugging)
EXPOSE 8080

# Set the entrypoint to run the client
ENTRYPOINT ["/client"]