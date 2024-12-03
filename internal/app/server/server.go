package server

import (
	"context"
	"time"

	"cangrejo_gigante/internal/infrastructure/network"
	"cangrejo_gigante/internal/logger"
)

type Server struct {
	powService       PowService
	quoteService     QuoteService
	nonceStore       *nonceStore
	connectionServer network.ConnectionServer
	log              logger.Logger
}

func New(powService PowService, quoteService QuoteService, nonceTTL time.Duration, secret string, connSrv network.ConnectionServer, logger logger.Logger) *Server {
	return &Server{
		powService:       powService,
		quoteService:     quoteService,
		nonceStore:       newNonceStore(nonceTTL, []byte(secret)),
		connectionServer: connSrv,
		log:              logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.log.Info("Shutting down server...")
	}()

	s.log.Info("Server starting...")
	return s.connectionServer.ListenAndServe()
}
