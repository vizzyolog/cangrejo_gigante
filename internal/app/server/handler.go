package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

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

func (h *Handler) Handle(conn net.Conn) {
	defer conn.Close()

	h.log.Infof("New connection from %s", conn.RemoteAddr().String())

	nonce, err := h.powService.GenerateChallenge()
	if err != nil {
		h.log.Errorf("Failed to generate challenge: %v", err)
		if _, err := fmt.Fprintf(conn, "Server error\n"); err != nil {
			h.log.Errorf("Failed to send error to client: %v", err)
		}

		return
	}

	h.log.Infof("Sending challenge to client: %s:%d", nonce.Nonce, nonce.Difficulty)

	_, err = fmt.Fprintf(conn, "%s:%d\n", nonce.Nonce, nonce.Difficulty)
	if err != nil {
		h.log.Errorf("Failed to send challenge: %v", err)
		return
	}

	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			h.log.Warnf("Client closed connection prematurely.")
			return
		}
		h.log.Errorf("Failed to read solution: %v", err)
		return
	}

	data = strings.TrimSpace(data)
	h.log.Infof("Received solution from client: %s", data)

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		h.log.Warnf("Invalid solution format from client: %s", data)
		fmt.Fprintf(conn, "Invalid solution format\n")
		return
	}

	clientNonce, clientSolution := parts[0], parts[1]

	if clientNonce != nonce.Nonce {
		h.log.Warnf("Invalid nonce from client: expected '%s', got '%s'", nonce.Nonce, clientNonce)
		fmt.Fprintf(conn, "Invalid nonce\n")
		return
	}

	if !h.powService.VerifySolution(clientNonce, clientSolution) {
		h.log.Warnf("Invalid solution from client. Nonce='%s', Solution='%s'", clientNonce, clientSolution)
		fmt.Fprintf(conn, "Wrong PoW\n")
		return
	}

	quote := h.quoteService.GetRandomQuote()
	h.log.Infof("Sending quote to client: %s", quote)

	_, err = fmt.Fprintf(conn, "%s\n", quote)
	if err != nil {
		h.log.Errorf("Failed to send quote: %v", err)
		return
	}

	h.log.Info("Quote sent successfully, connection closing.")
}
