package network

import (
	"fmt"
	"log"
	"net"
)

type TCPServer struct {
	address    string
	handleConn func(net.Conn)
}

func NewTCPServer(address string, handleConn func(net.Conn)) *TCPServer {
	return &TCPServer{
		address:    address,
		handleConn: handleConn,
	}
}

func (s *TCPServer) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start TCP server on %s: %w", s.address, err)
	}
	defer listener.Close()

	log.Printf("Listening on %s\n", s.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)

			continue
		}

		go s.handleConn(conn)
	}
}
