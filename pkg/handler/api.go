package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
)

func (h *Handler) GetSites(c *fiber.Ctx) error {

	sites := []data.Site{}
	if res := h.db.Find(&sites); res.Error != nil {
		return res.Error
	}

	return c.JSON(sites)
}

func (h *Handler) AddSite(c *fiber.Ctx) error {

	site := data.Site{}

	if err := c.BodyParser(&site); err != nil {
		return err
	}

	if res := h.db.Create(&site); res.Error != nil {
		return res.Error
	}

	return c.JSON(fiber.Map{
		"ID": site.ID,
	})
}
