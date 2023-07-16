package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/types"
	supa "github.com/nedpals/supabase-go"
)

func (h *Handler) Signup(c *fiber.Ctx) error {
	in := types.AuthInput{}
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
	in := types.AuthInput{}
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
