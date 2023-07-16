package handler

import (
	"github.com/gofiber/fiber/v2"
	supa "github.com/nedpals/supabase-go"
)

type AuthInput struct {
	Email    string
	Password string
}

func (h *Handler) Signup(c *fiber.Ctx) error {
	in := AuthInput{}
	if err := c.BodyParser(&in); err != nil {
		return err
	}

	_, err := h.supabase.Auth.SignUp(c.Context(), supa.UserCredentials{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	in := AuthInput{}
	if err := c.BodyParser(&in); err != nil {
		return err
	}

	u, err := h.supabase.Auth.SignIn(c.Context(), supa.UserCredentials{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		return err
	}

	return c.JSON(u)
}
