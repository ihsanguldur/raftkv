package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Get("/kv/:key", h.Get)
	app.Put("/kv/:key", h.Put)
	app.Delete("/kv/:key", h.Delete)
}
