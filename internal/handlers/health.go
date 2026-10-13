package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/BYUI-Study-Space/study-space-api/internal/store"
)

type Handler struct {
	store *store.Store
	tz    *time.Location
}

// New builds the handlers; tz is the campus timezone used to reject past sessions.
func New(s *store.Store, tz *time.Location) *Handler {
	return &Handler{store: s, tz: tz}
}

func (h *Handler) Health(c fiber.Ctx) error {
	if err := h.store.Ping(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "database unavailable"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
