package render

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
)

type NewAccountRequest struct {
	Service string `form:"service"`
	Account string `form:"account"`
	Secret  string `form:"secret"`
	OtpType string `form:"otpType"`
	Digits  uint   `form:"digits"`
}

func (h *Render) RenderNewAccount(c *fiber.Ctx) error {
	user, err := h.GetLoggedInUser(c)
	if err != nil {
		return err
	}

	if c.Method() == http.MethodGet {
		return c.Render("new", fiber.Map{})
	}

	a := NewAccountRequest{}

	if err := c.BodyParser(&a); err != nil {
		return err
	}

	acc := &data.Account{
		Service: a.Service,
		Account: a.Account,
		Secret:  a.Secret,
		OtpType: data.OtpType(a.OtpType),
		Digits:  a.Digits,
		UserID:  user.ID,
	}

	if res := h.db.Create(acc); res.Error != nil {
		return res.Error
	}

	return c.Redirect("/accounts")
}

func (h *Render) RenderAccounts(c *fiber.Ctx) error {
	user, err := h.GetLoggedInUser(c)
	if err != nil {
		return err
	}

	accs := []data.Account{}
	if res := h.db.Where("user_id=?", user.ID).Find(&accs); res.Error != nil {
		return res.Error
	}

	return c.Render("accounts", fiber.Map{
		"accounts": accs,
	})
}
