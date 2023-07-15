package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/xlzd/gotp"
)

func (h *Handler) RenderSite(c *fiber.Ctx) error {
	id := c.Params("id", "0")
	idx, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	sites := []data.Site{}
	if res := h.db.Find(&sites); res.Error != nil {
		return res.Error
	}

	if len(sites) == 0 {
		// 404
		return c.Render("home", fiber.Map{
			"sites": sites,
		})
	}

	current := sites[0]
	for _, s := range sites {
		if s.ID == uint(idx) {
			current = s
		}
	}

	totp := gotp.NewDefaultTOTP(current.Secret)
	code, exp := totp.NowWithExpiration()

	return c.Render("home", fiber.Map{
		"sites":   sites,
		"current": current,
		"code":    code,
		"exp":     exp,
	})
}

func (h *Handler) RenderSites(c *fiber.Ctx) error {
	sites := []data.Site{}
	if res := h.db.Find(&sites); res.Error != nil {
		return res.Error
	}

	return c.Render("home", fiber.Map{
		"sites": sites,
	})
}
