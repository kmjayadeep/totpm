package render

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
)

func (h *Render) RenderAccounts(c *fiber.Ctx) error {
	accs := []data.Account{}
	if res := h.db.Find(&accs); res.Error != nil {
		return res.Error
	}

	return c.Render("accounts", fiber.Map{
		"accounts": accs,
	})
}
