package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/repository/local"
	"backend/internal/repository/postgres"
	"backend/internal/service"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo, cleanup, err := newRepository(cfg)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}
	defer cleanup()

	linksService := service.NewLinksService(repo)
	linksHandler := handler.NewLinksHandler(linksService, logger)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := repo.Ping(pingCtx); err != nil {
		log.Fatalf("storage unreachable: %v", err)
	}

	app := fiber.New()
	linksHandler.Register(app)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.Listen(fmt.Sprintf(":%d", cfg.ServerPort))
	}()

	select {
	case err := <-serverErr:
		log.Fatalf("server: %v", err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func newRepository(cfg *config.Config) (service.LinksRepository, func(), error) {
	switch cfg.Storage {
	case config.StoragePostgres:
		pool, err := pgxpool.New(context.Background(), cfg.GetDBDSN())
		if err != nil {
			return nil, nil, fmt.Errorf("pgx pool: %w", err)
		}
		return postgres.NewLinksRepository(pool), pool.Close, nil
	case config.StorageLocal:
		return local.NewLinksRepository(), func() {}, nil
	default:
		return nil, nil, fmt.Errorf("unknown storage type: %q", cfg.Storage)
	}
}
