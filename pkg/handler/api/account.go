package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/kmjayadeep/totpm/pkg/types"
)

func (h *Api) GetAccounts(c *fiber.Ctx) error {
	accounts := []data.Account{}
	if res := h.db.Find(&accounts); res.Error != nil {
		return res.Error
	}

	return c.JSON(accounts)
}

func (h *Api) AddAccount(c *fiber.Ctx) error {
	req := types.Account{}
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	acc := data.Account{
		Service:   req.Service,
		Account:   req.Account,
		Icon:      req.Icon,
		OtpType:   data.OtpType(req.OtpType),
		Digits:    req.Digits,
		Algorithm: req.Algorithm,
		Period:    req.Period,
		Counter:   req.Counter,
	}

	acc.SetSecret(req.Secret)

	if res := h.db.Create(&acc); res.Error != nil {
		return res.Error
	}

	return c.JSON(fiber.Map{
		"ID": acc.ID,
	})
}

func (h *Api) DeleteAccounts(c *fiber.Ctx) error {
	ids := strings.Split(c.Query("ids", ""), ",")

	if res := h.db.Where("id in (?)", ids).Delete(&data.Account{}); res.Error != nil {
		return res.Error
	}

	return c.SendStatus(fiber.StatusNoContent)
}
