package main

import (
	"backend/src/internal/config"
	"backend/src/internal/handler"
	"backend/src/internal/repository"
	"backend/src/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var repo repository.ILinksRepository

	switch cfg.Storage {
	case "postgres":
		pool, err := pgxpool.New(context.Background(), cfg.GetDBDSN())
		if err != nil {
			log.Fatal(err)
		}

		repo, err = repository.NewLinksRepository("postgres", pool)
		if err != nil {
			log.Fatal(err)
		}
	case "local":
		repo, err = repository.NewLinksRepository("local", nil)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown storage type: %w", cfg.Storage)
	}

	linkService := service.NewLinksService(repo)
	linkHandler := handler.NewLinkHandler(linkService)

	app := fiber.New()

	app.Post("/shorten", linkHandler.Shorten)
	app.Get("/:shortLink", linkHandler.Redirect)

	go func() {
		if err := app.Listen(":" + strconv.Itoa(cfg.ServerPort)); err != nil {
			log.Print(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	shutdowCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = app.ShutdownWithContext(shutdowCtx); err != nil {
		log.Print(err)
	}
}
