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

# Build the server binary
RUN go build -o server cmd/server/main.go

# Use a minimal image for running the server
FROM debian:bullseye-slim

# Copy the server binary
COPY --from=builder /app/server /server

# Create the configs directory and copy config files
RUN mkdir /configs
COPY configs/config.yaml /configs/config.yaml
COPY configs/quotes.txt /configs/quotes.txt

# Expose the server port
EXPOSE 8080

# Set the entrypoint to run the server
ENTRYPOINT ["/server"]