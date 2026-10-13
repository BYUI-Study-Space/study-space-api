// Command api runs the StudySpace HTTP API.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata" // the distroless image has no timezone database

	"github.com/gofiber/fiber/v3"

	"github.com/BYUI-Study-Space/study-space-api/internal/config"
	"github.com/BYUI-Study-Space/study-space-api/internal/database"
	"github.com/BYUI-Study-Space/study-space-api/internal/handlers"
	"github.com/BYUI-Study-Space/study-space-api/internal/server"
	"github.com/BYUI-Study-Space/study-space-api/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	tz, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		return fmt.Errorf("load timezone %q: %w", cfg.Timezone, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := database.Migrate(pool); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	app := server.New(cfg.AllowedOrigins, handlers.New(store.New(pool), tz))
	log.Printf("listening on :%s", cfg.Port)
	return app.Listen(":"+cfg.Port, fiber.ListenConfig{
		DisableStartupMessage: true,
		GracefulContext:       ctx,
	})
}
