package server

import (
	"context"
	"fmt"

	"cangrejo_gigante/internal/infrastructure/network"
	"cangrejo_gigante/internal/logger"
)

type Server struct {
	powService       PowService
	quoteService     QuoteService
	connectionServer network.ConnectionServer
	log              logger.Logger
}

func New(
	powService PowService,
	quoteService QuoteService,
	connSrv network.ConnectionServer,
	logger logger.Logger) *Server {
	return &Server{
		powService:       powService,
		quoteService:     quoteService,
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

	if err := s.connectionServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
