package handler

import (
	"github.com/gofiber/fiber/v2"
)

const LocalKeyUser = "user"

func (h *Handler) RequiresAuth(c *fiber.Ctx) error {
	token := c.Get("x-access-token")

	if token == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	u, err := h.supabase.Auth.User(c.Context(), token)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("user", u)
	return c.Next()
}
