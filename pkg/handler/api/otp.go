package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/internal/otp"
	"github.com/kmjayadeep/totpm/pkg/data"
)

func (h *Api) GetAccountOTP(c *fiber.Ctx) error {
	acc := data.Account{}
	id := c.Params("id")

	if res := h.db.First(&acc, id); res.Error != nil {
		return res.Error
	}

	code, err := otp.GenerateCode(acc)
	if err != nil {
		return err
	}

	return c.JSON(code)
}
