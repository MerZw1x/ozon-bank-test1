package main

import (
	"backend/src/internal/config"
	"backend/src/internal/handler"
	"backend/src/internal/repository/local"
	"backend/src/internal/repository/postgre"
	"backend/src/internal/service"
	"context"
	"fmt"
	"log"
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

	repo, cleanup, err := newRepository(cfg)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}
	defer cleanup()

	linksService := service.NewLinksService(repo)
	linksHandler := handler.NewLinksHandler(linksService)

	app := fiber.New()
	linksHandler.Register(app)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.ServerPort)
		if err := app.Listen(addr); err != nil {
			log.Printf("server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func newRepository(cfg *config.Config) (service.LinksRepository, func(), error) {
	switch cfg.Storage {
	case "postgres":
		pool, err := pgxpool.New(context.Background(), cfg.GetDBDSN())
		if err != nil {
			return nil, nil, fmt.Errorf("pgx pool: %w", err)
		}
		return postgre.NewLinksRepository(pool), pool.Close, nil
	case "local":
		return local.NewLinksRepository(), func() {}, nil
	default:
		return nil, nil, fmt.Errorf("unknown storage type: %q", cfg.Storage)
	}
}
