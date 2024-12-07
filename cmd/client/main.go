package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"cangrejo_gigante/internal/app/client"
	"cangrejo_gigante/internal/config"
	"cangrejo_gigante/internal/domain/pow"
	"cangrejo_gigante/internal/infrastructure/network"
	"cangrejo_gigante/internal/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("Failed to load config: %v", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Client.Timeout)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Info("Received termination signal, shutting down...")
		close(stopped)
		cancel()
	}()

	connManager := network.NewConnectionManager(cfg.Server.Address)
	powResolver := pow.NewPoWResolver(cfg.PoW.Difficulty)

	app := client.NewClient(connManager, powResolver, log)

	go func() {
		if err := app.Run(ctx); err != nil {
			log.Errorf("Client failed: %v", err)
		}
	}()

	<-stopped

	log.Info("app stopped")
}
