package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/BYUI-Study-Space/study-space-api/internal/models"
	"github.com/BYUI-Study-Space/study-space-api/internal/store"
)

// ListLocations handles GET /api/locations?wifi=true&group=true&noise=quiet.
func (h *Handler) ListLocations(c fiber.Ctx) error {
	var filter models.LocationFilter
	var err error
	if filter.Wifi, err = boolQuery(c, "wifi"); err != nil {
		return err
	}
	if filter.Group, err = boolQuery(c, "group"); err != nil {
		return err
	}
	if filter.Noise = c.Query("noise"); filter.Noise != "" && !models.NoiseLevels[filter.Noise] {
		return fiber.NewError(fiber.StatusBadRequest, `noise must be "quiet", "moderate" or "collaborative"`)
	}

	locations, err := h.store.ListLocations(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(locations)
}

// GetLocation handles GET /api/locations/:id.
func (h *Handler) GetLocation(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "location not found")
	}
	location, err := h.store.GetLocation(c.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "location not found")
	}
	if err != nil {
		return err
	}
	return c.JSON(location)
}

// boolQuery parses an optional boolean query parameter; it returns nil when the parameter is absent.
func boolQuery(c fiber.Ctx, key string) (*bool, error) {
	raw := c.Query(key)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, key+" must be true or false")
	}
	return &v, nil
}
