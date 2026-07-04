package app

import (
	"backend/src/internal/config"
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	switch cfg.Storage {
	case "postgres":
	case "local":
	}

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
