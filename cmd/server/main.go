package main

import (
	"context"

	"cangrejo_gigante/internal/app/server"
	"cangrejo_gigante/internal/config"
	"cangrejo_gigante/internal/domain/pow"
	"cangrejo_gigante/internal/domain/quote"
	"cangrejo_gigante/internal/infrastructure/network"
	"cangrejo_gigante/internal/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("Failed to Load config :%v", err)

		return
	}

	powService := pow.New(cfg.PoW.Difficulty, log)

	quoteService, err := quote.New(cfg.Quotes.FilePath)
	if err != nil {
		log.Errorf("Failed to create quoteService %v", err)

		return
	}

	ctx := context.Background()

	nonceStore := server.NewNonceStore(cfg.Server.NonceTTL)
	handler := server.NewHandler(powService, quoteService, nonceStore, cfg.Server.MaxDataSize, cfg.Server.MaxConn, log)
	tcpServer := network.NewTCPServer(cfg.Server.Address, handler.Handle)

	srv := server.New(powService, quoteService, tcpServer, log)
	if err := srv.Run(ctx); err != nil {
		log.Errorf("Failed to start server: %v", err)

		return
	}
}
