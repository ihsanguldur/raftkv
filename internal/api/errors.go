package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ihsanguldur/raftkv/internal/service"
)

func writeError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrKeyNotFound):
		return c.Status(fiber.StatusNotFound).JSON(errorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrKeyRequired), errors.Is(err, service.ErrValueRequired):
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "internal server error"})
	}
}
