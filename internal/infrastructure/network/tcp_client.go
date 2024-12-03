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
	return &ConnectionManager{Address: address}
}

func (cm *ConnectionManager) Connect() error {
	conn, err := net.Dial("tcp", cm.Address)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	cm.conn = conn
	return nil
}

func (cm *ConnectionManager) Close() error {
	if cm.conn != nil {
		return cm.conn.Close()
	}
	return nil
}

func (cm *ConnectionManager) Send(data string) error {
	_, err := fmt.Fprintf(cm.conn, "%s\n", data)
	return err
}

func (c *ConnectionManager) Receive() (string, error) {
	reader := bufio.NewReader(c.conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read from connection: %w", err)
	}
	return strings.TrimSpace(data), nil
}
