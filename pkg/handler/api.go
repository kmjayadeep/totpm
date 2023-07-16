package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
	supa "github.com/nedpals/supabase-go"
)

func (h *Handler) GetSites(c *fiber.Ctx) error {
	user := c.Locals(LocalKeyUser).(*supa.User)

	sites := []data.Site{}
	if res := h.db.Where("user_id=?", user.ID).Find(&sites); res.Error != nil {
		return res.Error
	}

	return c.JSON(sites)
}

func (h *Handler) GetSite(c *fiber.Ctx) error {
	user := c.Locals(LocalKeyUser).(*supa.User)

	id := c.Params("id")

	site := data.Site{}
	if res := h.db.Where("user_id=? and id=?", user.ID, id).First(&site); res.Error != nil {
		return res.Error
	}

	return c.JSON(site)
}

func (h *Handler) AddSite(c *fiber.Ctx) error {
	user := c.Locals(LocalKeyUser).(*supa.User)

	site := data.Site{}
	if err := c.BodyParser(&site); err != nil {
		return err
	}
	site.UserID = user.ID

	if res := h.db.Create(&site); res.Error != nil {
		return res.Error
	}

	return c.JSON(fiber.Map{
		"ID": site.ID,
	})
}
