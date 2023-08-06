package render

import (
	"encoding/gob"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func init() {
	gob.Register(&LoginRequest{})
}

func (h *Render) RenderRegister(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{})
}

func (h *Render) RenderLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func (h *Render) RenderLoginSubmit(c *fiber.Ctx) error {

	l := LoginRequest{}

	if err := c.BodyParser(&l); err != nil {
		return err
	}

	if l.Email != "" && l.Password != "" {
		sess, err := h.store.Get(c)
		if err != nil {
			return err
		}
		sess.Set("user", l)

		if err := sess.Save(); err != nil {
			return err
		}

		return c.Redirect("/accounts")
	}

	return c.Render("login", fiber.Map{})
}
