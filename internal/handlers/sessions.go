package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/BYUI-Study-Space/study-space-api/internal/models"
	"github.com/BYUI-Study-Space/study-space-api/internal/store"
)

// CreateSession handles POST /api/sessions.
func (h *Handler) CreateSession(c fiber.Ctx) error {
	var in models.SessionInput
	if err := c.Bind().JSON(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "request body must be valid JSON")
	}
	in.Normalize()

	var loc *models.Location
	if in.LocationID != 0 {
		l, err := h.store.GetLocation(c.Context(), in.LocationID)
		switch {
		case errors.Is(err, store.ErrNotFound):
			// Validate reports the unknown location.
		case err != nil:
			return err
		default:
			loc = &l
		}
	}

	if fields := in.Validate(loc, time.Now().In(h.tz)); len(fields) > 0 {
		return validationError(c, fields)
	}

	session, err := h.store.CreateSession(c.Context(), in)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(session)
}

// GetSession handles GET /api/sessions/:id.
func (h *Handler) GetSession(c fiber.Ctx) error {
	session, err := h.store.GetSession(c.Context(), c.Params("id"))
	if errors.Is(err, store.ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "session not found")
	}
	if err != nil {
		return err
	}
	return c.JSON(session)
}
