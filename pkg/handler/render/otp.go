package render

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/internal/otp"
	"github.com/kmjayadeep/totpm/pkg/data"
)

func (h *Render) RenderOtp(c *fiber.Ctx) error {
	user, err := h.GetLoggedInUser(c)
	if err != nil {
		return err
	}

	acc := data.Account{}
	id := c.Params("id")

	if res := h.db.Where("id=? and user_id=?", id, user.ID).Find(&acc); res.Error != nil {
		return res.Error
	}

	code, err := otp.GenerateCode(acc)
	if err != nil {
		return err
	}

	return c.Render("partials/otp", fiber.Map{
		"ID":  id,
		"Otp": code,
	}, "layouts/htmx")
}
