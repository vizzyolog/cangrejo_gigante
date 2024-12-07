package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"cangrejo_gigante/internal/domain/pow"
	"cangrejo_gigante/internal/logger"
)

type Handler struct {
	powService   PowService
	quoteService QuoteService
	nonceStore   *NonceStore
	maxDataSize  int
	log          logger.Logger
	sem          chan struct{}
}

func NewHandler(
	powService PowService,
	quoteService QuoteService,
	nonceStore *NonceStore,
	maxDataSize int,
	maxConn int,
	log logger.Logger) *Handler {
	return &Handler{
		powService:   powService,
		quoteService: quoteService,
		nonceStore:   nonceStore,
		maxDataSize:  maxDataSize,
		log:          log,
		sem:          make(chan struct{}, maxConn),
	}
}

var ErrInvalidSolutionFormat = errors.New("invalid solution format")
var ErrInvalidSolution = errors.New("invalid or expired nonce")
var ErrWrongPow = errors.New("wrong PoW")
var ErrDataSizeExceeds = errors.New("data size exceeds limit")

func (h *Handler) Handle(conn net.Conn) {
	if !h.acquireSlot(conn) {
		return
	}
	defer conn.Close()

	challenge, err := h.generateAndSaveChallenge(conn)
	if err != nil {
		return
	}

	if err := h.sendChallenge(conn, challenge); err != nil {
		h.sendError(conn, "Failed to send challenge", err)

		return
	}

	if err := h.receiveAndVerifySolution(conn); err != nil {
		return
	}

	if err := h.sendQuoteToClient(conn); err != nil {
		return
	}

	h.log.Info("Quote sent successfully, connection closing.")
}

func (h *Handler) acquireSlot(conn net.Conn) bool {
	select {
	case h.sem <- struct{}{}:
		defer func() { <-h.sem }()

		return true
	default:
		h.sendError(conn, "Server is busy, try again later", nil)
		conn.Close()

		return false
	}
}

func (h *Handler) generateAndSaveChallenge(conn net.Conn) (*pow.Challenge, error) {
	h.log.Infof("New connection from %s", conn.RemoteAddr().String())

	challenge, err := h.powService.GenerateChallenge()
	if err != nil {
		h.sendError(conn, "Server error", err)

		return nil, fmt.Errorf("failed to generate challenge: %w", err)
	}

	if err := h.nonceStore.Save(challenge.Nonce); err != nil {
		h.sendError(conn, "Server error", err)

		return nil, fmt.Errorf("failed to save nonce: %w", err)
	}

	h.log.Infof("Sending challenge to client: %s:%d", challenge.Nonce, challenge.Difficulty)

	return challenge, nil
}

func (h *Handler) receiveAndVerifySolution(conn net.Conn) error {
	data, err := h.receiveDataFromClient(conn)
	if err != nil {
		return err
	}

	clientNonce, clientSolution, err := parseClientNonceAndSolutionFromData(data)
	if err != nil {
		h.sendError(conn, "Invalid solution format", nil)

		return err
	}

	if !h.nonceStore.IsValid(clientNonce) {
		h.sendError(conn, "Invalid or expired nonce", nil)

		return ErrInvalidSolution
	}

	if err := h.verifySolution(conn, clientNonce, clientSolution); err != nil {
		h.sendError(conn, "wrong solution", err)

		return err
	}

	h.nonceStore.MarkAsUsed(clientNonce)

	return nil
}

func (h *Handler) sendQuoteToClient(conn net.Conn) error {
	quote := h.quoteService.GetRandomQuote()

	h.log.Infof("Sending quote to client: %s", quote)

	if err := h.sendQuote(conn, quote); err != nil {
		h.sendError(conn, "Failed to send quote", err)

		return err
	}

	return nil
}

func (h *Handler) sendError(conn net.Conn, message string, err error) {
	h.log.Info("message: ", message)
	h.log.Info("error: ", err)

	if err != nil {
		h.log.Errorf("%s: %v", message, err)
	} else {
		h.log.Warn(message)
	}

	_, writeErr := fmt.Fprintf(conn, "%s\n", message)
	if writeErr != nil {
		h.log.Errorf("Failed to send error message to client: %v", writeErr)
	}
}

func (h *Handler) sendChallenge(conn net.Conn, challenge *pow.Challenge) error {
	_, err := fmt.Fprintf(conn, "%s:%d\n", challenge.Nonce, challenge.Difficulty)
	if err != nil {
		return fmt.Errorf("failed to send challenge: %w", err)
	}

	return nil
}

func (h *Handler) receiveDataFromClient(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	data, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			h.log.Warnf("Client closed connection prematurely.")

			return "", nil
		}

		h.log.Errorf("Failed to read solution: %v", err)

		return "", fmt.Errorf("failed to reciveDataFromClient: %w", err)
	}

	if len(data) > h.maxDataSize {
		return "", ErrDataSizeExceeds
	}

	return strings.TrimSpace(data), nil
}

func (h *Handler) verifySolution(
	conn net.Conn,
	clientNonce,
	clientSolution string) error {
	if !h.powService.VerifySolution(clientNonce, clientSolution) {
		h.sendError(conn, "Wrong PoW", nil)

		return ErrWrongPow
	}

	return nil
}

func (h *Handler) sendQuote(conn net.Conn, quote string) error {
	_, err := fmt.Fprintf(conn, "%s\n", quote)
	if err != nil {
		return fmt.Errorf("failed to send quote: %w", err)
	}

	return nil
}

func parseClientNonceAndSolutionFromData(data string) (string, string, error) {
	parts := strings.Split(data, ":")
	if len(parts) != pow.ExpectedDataPartsCount {
		return "", "", ErrInvalidSolutionFormat
	}

	clientNonce, clientSolution := parts[0], parts[1]

	return clientNonce, clientSolution, nil
}
