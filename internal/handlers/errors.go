// Package handlers contains the Fiber HTTP handlers.
package handlers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// ErrorHandler turns any error returned by a handler into a JSON response.
// Unexpected errors are logged and hidden behind a generic 500 message.
func ErrorHandler(c fiber.Ctx, err error) error {
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return c.Status(fe.Code).JSON(errorResponse{Error: fe.Message})
	}
	log.Printf("%s %s: %v", c.Method(), c.Path(), err)
	return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "internal server error"})
}

func validationError(c fiber.Ctx, fields map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "validation failed", Fields: fields})
}
