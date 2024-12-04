package network

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type ConnectionManager struct {
	Address string
	conn    net.Conn
}

func NewConnectionManager(address string) *ConnectionManager {
	return &ConnectionManager{
		Address: address,
		conn:    nil,
	}
}

func (cm *ConnectionManager) Connect() error {
	conn, err := net.Dial("tcp", cm.Address)
	if err != nil {
		return fmt.Errorf("failed to connect to server:%s, %w", cm.Address, err)
	}

	cm.conn = conn

	return nil
}

func (cm *ConnectionManager) Close() error {
	if err := cm.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

func (cm *ConnectionManager) Send(data string) error {
	_, err := fmt.Fprintf(cm.conn, "%s\n", data)

	return fmt.Errorf("failed to send data: %w", err)
}

func (cm *ConnectionManager) Receive() (string, error) {
	reader := bufio.NewReader(cm.conn)

	data, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read from connection: %w", err)
	}

	return strings.TrimSpace(data), nil
}
