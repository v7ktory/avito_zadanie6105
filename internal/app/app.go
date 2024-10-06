package app

import (
	"fmt"
	"log"
	sl "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.codenrock.com/tender/config"
	v1 "git.codenrock.com/tender/internal/controller/http/v1"
	"git.codenrock.com/tender/internal/repo"
	"git.codenrock.com/tender/internal/service"
	"git.codenrock.com/tender/pkg/postgres"
	"git.codenrock.com/tender/pkg/server"
)

func Run(configPath string) {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Set up logger
	setUpLogger(cfg.Log.Level)

	// Postgres connection
	sl.Info("Connecting to Postgres...")
	pg, err := postgres.New(cfg.Conn, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	// Repositories
	sl.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	sl.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos: repositories,
	}
	services := service.NewServices(deps)

	// Mux handler
	sl.Info("Initializing handlers and routes...")
	handler := http.NewServeMux()
	// setup handler validator as lib validator
	v1.NewRouter(handler, services)

	// HTTP server
	sl.Info("Starting http server...")
	sl.Debug("Server address", sl.Any("address", cfg.Adress))
	httpServer := server.New(handler, server.Address(cfg.Adress))

	// Waiting signal
	sl.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		sl.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		sl.Error("app - Run - httpServer.Notify: ", sl.Any("error", err.Error()))
	}

	// Graceful shutdown
	sl.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		sl.Error("app - Run - httpServer.Shutdown: ", sl.Any("error", err.Error()))
	}
}
