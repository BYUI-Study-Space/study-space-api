// Package server wires the Fiber app: middleware and routes.
package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/BYUI-Study-Space/study-space-api/internal/handlers"
)

func New(allowedOrigins []string, h *handlers.Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "StudySpace API",
		ErrorHandler: handlers.ErrorHandler,
		BodyLimit:    16 * 1024,
	})

	app.Use(recoverer.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodOptions},
		AllowHeaders: []string{fiber.HeaderContentType},
	}))

	api := app.Group("/api")
	api.Get("/health", h.Health)
	api.Get("/locations", h.ListLocations)
	api.Get("/locations/:id", h.GetLocation)
	api.Post("/sessions", h.CreateSession)
	api.Get("/sessions/:id", h.GetSession)

	return app
}
