package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cangrejo_gigante/internal/domain/pow"
	"cangrejo_gigante/internal/infrastructure/network"
	"cangrejo_gigante/internal/logger"
)

type Client struct {
	connManager *network.ConnectionManager
	powResolver *pow.PoWResolver
	log         logger.Logger
	ctx         context.Context
}

func NewClient(connManager *network.ConnectionManager, powResolver *pow.PoWResolver, ctx context.Context, log logger.Logger) *Client {
	return &Client{
		connManager: connManager,
		powResolver: powResolver,
		log:         log,
		ctx:         ctx,
	}
}

func (c *Client) Run() error {
	c.log.Info("Starting client...")

	select {
	case <-c.ctx.Done():
		c.log.Warn("Context canceled before starting client")
		return c.ctx.Err()
	default:
	}

	if err := c.connManager.Connect(); err != nil {
		c.log.Errorf("Failed to connect: %v", err)
		return fmt.Errorf("connection failed: %w", err)
	}
	defer c.connManager.Close()

	data, err := c.connManager.Receive()
	if err != nil {
		c.log.Errorf("Failed to receive data: %v", err)
		return fmt.Errorf("failed to receive data: %w", err)
	}

	c.log.Infof("Raw data received from server: '%s'", data)

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid challenge format: %s", data)
	}

	challenge := pow.Challenge{
		Nonce: parts[0],
	}

	difficulty, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid difficulty format: %w", err)
	}
	challenge.Difficulty = difficulty

	c.log.Infof("Parsed challenge: Nonce='%s', Difficulty=%d", challenge.Nonce, challenge.Difficulty)

	solution, err := pow.SolveChallenge(c.ctx, &challenge)
	if err != nil {
		return fmt.Errorf("failed to solve PoW: %w", err)
	}

	c.log.Infof("Solved challenge: Nonce='%s', Solution='%s'", solution.Nonce, solution.Response)

	if err := c.connManager.Send(fmt.Sprintf("%s:%s", solution.Nonce, solution.Response)); err != nil {
		c.log.Errorf("Failed to send solution: %v", err)
		return fmt.Errorf("failed to send solution: %w", err)
	}

	response, err := c.connManager.Receive()
	if err != nil {
		if err.Error() == "EOF" {
			c.log.Info("Connection closed by server.")
			return nil
		}
		c.log.Errorf("Failed to receive response: %v", err)
		return fmt.Errorf("failed to receive response: %w", err)
	}

	c.log.Infof("Server response: %s", strings.TrimSpace(response))
	return nil
}
