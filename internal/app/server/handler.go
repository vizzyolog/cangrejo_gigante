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
	log          logger.Logger
}

func NewHandler(powService PowService, quoteService QuoteService, log logger.Logger) *Handler {
	return &Handler{
		powService:   powService,
		quoteService: quoteService,
		log:          log,
	}
}

var ErrInvalidSolutionFormat = errors.New("invalid solution format")
var ErrWrongPow = errors.New("wrong PoW")

func (h *Handler) Handle(conn net.Conn) {
	defer conn.Close()

	h.log.Infof("New connection from %s", conn.RemoteAddr().String())

	nonce, err := h.powService.GenerateChallenge()
	if err != nil {
		h.sendError(conn, "Server error", err)

		return
	}

	h.log.Infof("Sending challenge to client: %s:%d", nonce.Nonce, nonce.Difficulty)

	if err := h.sendChallenge(conn, nonce); err != nil {
		h.sendError(conn, "Failed to send challenge", err)

		return
	}

	data, err := h.receiveDataFromClient(conn)
	if err != nil {
		return
	}

	clientNonce, clientSolution, err := parseClientNonceAndSolutionFromData(data)
	if err != nil {
		h.sendError(conn, "Invalid solution format", nil)

		return
	}

	if err := h.verifySolution(conn, nonce, clientNonce, clientSolution); err != nil {
		return
	}

	quote := h.quoteService.GetRandomQuote()

	h.log.Infof("Sending quote to client: %s", quote)

	if err := h.sendQuote(conn, quote); err != nil {
		h.sendError(conn, "Failed to send quote", err)

		return
	}

	h.log.Info("Quote sent successfully, connection closing.")
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

func (h *Handler) sendChallenge(conn net.Conn, nonce *pow.Challenge) error {
	_, err := fmt.Fprintf(conn, "%s:%d\n", nonce.Nonce, nonce.Difficulty)
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

	return strings.TrimSpace(data), nil
}

func (h *Handler) verifySolution(
	conn net.Conn,
	nonce *pow.Challenge,
	clientNonce,
	clientSolution string) error {
	if clientNonce != nonce.Nonce {
		h.sendError(conn, "Invalid nonce", nil)

		return ErrInvalidNonce
	}

	if !h.powService.VerifySolution(clientNonce, clientSolution) {
		h.sendError(conn, "Wrong PoW", nil)

		return ErrWrongPow
	}

	return nil
}

func (h *Handler) sendQuote(conn net.Conn, quote string) error {
	_, err := fmt.Fprintf(conn, "%s\n", quote)
	if err != nil {
		fmt.Errorf("failed to send quote: %w", err)
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
