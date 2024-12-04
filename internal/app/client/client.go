package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cangrejo_gigante/internal/domain/pow"
	"cangrejo_gigante/internal/infrastructure/network"
	"cangrejo_gigante/internal/logger"
)

type Client struct {
	connManager *network.ConnectionManager
	powResolver *pow.Resolver
	log         logger.Logger
}

func NewClient(connManager *network.ConnectionManager, powResolver *pow.Resolver, log logger.Logger) *Client {
	return &Client{
		connManager: connManager,
		powResolver: powResolver,
		log:         log,
	}
}

var ErrInvalidChallengeFormat = errors.New("invalid challenge format")

func (c *Client) Run(ctx context.Context) error {
	c.log.Info("Starting client...")

	if err := c.checkContextCancellation(ctx); err != nil {
		return err
	}

	if err := c.connectToServer(); err != nil {
		return err
	}
	defer c.connManager.Close()

	data, err := c.receiveDataFromServer()
	if err != nil {
		return err
	}

	nonce, difficulty, err := c.parseChallengeData(data)
	if err != nil {
		return err
	}

	challenge := pow.Challenge{
		Nonce:      nonce,
		Difficulty: difficulty,
	}

	solution, err := c.solveChallenge(ctx, &challenge)
	if err != nil {
		return err
	}

	if err := c.sendSolution(solution); err != nil {
		return err
	}

	if err := c.receiveResponseFromServer(); err != nil {
		return err
	}

	return nil
}

func (c *Client) checkContextCancellation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		c.log.Warn("Context canceled before starting client")

		return fmt.Errorf("context canceled: %w", ctx.Err())
	default:
	}

	return nil
}

func (c *Client) receiveDataFromServer() (string, error) {
	data, err := c.connManager.Receive()
	if err != nil {
		c.log.Errorf("Failed to receive data: %v", err)

		return "", fmt.Errorf("failed to receive data: %w", err)
	}

	c.log.Infof("Raw data received from server: '%s'", data)

	return data, nil
}

func (c *Client) connectToServer() error {
	if err := c.connManager.Connect(); err != nil {
		c.log.Errorf("Failed to connect: %v", err)

		return fmt.Errorf("connection failed: %w", err)
	}

	return nil
}

func (c *Client) solveChallenge(ctx context.Context, challenge *pow.Challenge) (*pow.Solution, error) {
	solution, err := pow.SolveChallenge(ctx, challenge)
	if err != nil {
		return nil, fmt.Errorf("failed to solve PoW: %w", err)
	}

	c.log.Infof("Solved challenge: Nonce='%s', Solution='%s'", solution.Nonce, solution.Response)

	return solution, nil
}

func (c *Client) sendSolution(solution *pow.Solution) error {
	if err := c.connManager.Send(fmt.Sprintf("%s:%s", solution.Nonce, solution.Response)); err != nil {
		c.log.Errorf("Failed to send solution: %v", err)

		return fmt.Errorf("failed to send solution: %w", err)
	}

	return nil
}

func (c *Client) receiveResponseFromServer() error {
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

func (c *Client) parseChallengeData(data string) (string, int, error) {
	parts := strings.Split(data, ":")
	if len(parts) != pow.ExpectedDataPartsCount {
		return "", 0, fmt.Errorf("%w, %s", ErrInvalidChallengeFormat, data)
	}

	nonce := parts[0]

	difficulty, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid difficulty format: %w", err)
	}

	return nonce, difficulty, nil
}
