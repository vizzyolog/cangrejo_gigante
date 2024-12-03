# Word of Wisdom TCP Client-Server

This project implements a client-server application over TCP, designed to solve a challenge-response "Proof of Work" (PoW) mechanism. Upon solving the challenge successfully, the server responds with a random quote from the "Word of Wisdom" collection.

---

## Features

1. **Proof of Work Protection:**
    - The server generates a PoW challenge (nonce and difficulty).
    - The client solves the PoW challenge by finding a solution that meets the required difficulty.
    - PoW ensures that computational work is required, protecting the server from DDoS attacks.

2. **Graceful Shutdown:**
    - Both the client and server handle termination signals (`SIGINT`, `SIGTERM`) and shut down gracefully.

3. **Timeouts and Context Management:**
    - The client uses `context.Context` for timeout control during the PoW challenge solution.

4. **Logger Integration:**
    - Custom logging is implemented for both client and server to track activity and debug issues.

5. **Extensible Design:**
    - Decoupled components for easy extension and testing (e.g., connection management, PoW, and quote generation).

6. **Docker Integration:**
    - Dockerfile is included for building and running both the client and server in containerized environments.

---

## How It Works

1. **Server Workflow:**
    - Listens for incoming TCP connections.
    - Sends a PoW challenge (nonce and difficulty) to the client.
    - Validates the client's solution against the PoW requirements.
    - Sends a random quote if the solution is valid, or an error message otherwise.

2. **Client Workflow:**
    - Connects to the server.
    - Receives the PoW challenge.
    - Computes a valid solution using brute-force.
    - Sends the solution back to the server.
    - Displays the server's response (a random quote or an error message).

---

## Installation

### Requirements

- Go 1.20+ installed on your system
- Docker (optional, for containerized execution)

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/vizzyolog/cangrejo_gigante.git
   cd cangrejo_gigante

2.	Install dependencies (if needed):
    
    ```bash
    go mod tidy
   
3.	Build the server and client:
    
    ```bash
   go build -o ./bin/server cmd/server/main.go
   go build -o ./bin/client cmd/client/main.go
   
4.	Run the server:
    
    ```bash
   ./server


5. Run the client:

    ```bash
    ./server

### Configuration
The application uses a YAML configuration file located at configs/config.yaml:
    
```yaml
    server:
      address: "localhost:8080"
      nonceTTL: 30s
    
    pow:
      difficulty: 20
    
    quotes:
      file_path: "configs/quotes.txt"