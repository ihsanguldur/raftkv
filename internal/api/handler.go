package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ihsanguldur/raftkv/internal/service"
)

type Handler struct {
	svc *service.KVService
}

func NewHandler(svc *service.KVService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Get(c *fiber.Ctx) error {
	v, err := h.svc.Get(c.Params("key"))
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(valueResponse{Value: v})
}

func (h *Handler) Put(c *fiber.Ctx) error {
	var req putRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid body"})
	}
	key := strings.Clone(c.Params("key"))
	if err := h.svc.Put(key, req.Value); err != nil {
		return writeError(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	key := strings.Clone(c.Params("key"))
	if err := h.svc.Delete(key); err != nil {
		return writeError(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
