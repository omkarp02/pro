package test

import "github.com/gofiber/fiber/v2"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/test", h.test)
}

func (h *Handler) test(c *fiber.Ctx) error {
	return c.SendString("hello how are you doing")
}
