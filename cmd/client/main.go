package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cangrejo_gigante/internal/app/client"
	"cangrejo_gigante/internal/config"
	"cangrejo_gigante/internal/domain/pow"
	"cangrejo_gigante/internal/infrastructure/network"
	"cangrejo_gigante/internal/logger"
)

func main() {
	log := logger.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("Failed to load config: %v", err)
		return
	}

	connManager := network.NewConnectionManager(cfg.Server.Address)
	powResolver := pow.NewPoWResolver(cfg.PoW.Difficulty)

	app := client.NewClient(connManager, powResolver, ctx, log)

	go func() {
		if err := app.Run(); err != nil {
			log.Errorf("Client failed: %v", err)
		}
		close(stopped)
	}()

	<-stopped

	log.Info("app stopped")
}
