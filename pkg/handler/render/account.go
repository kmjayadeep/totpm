package render

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
)

type NewAccountRequest struct {
	Service string `form:"service"`
	Account string `form:"account"`
	Secret  string `form:"secret"`
	Icon    string `form:"icon"`
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
		OtpType: data.OtpType(a.OtpType),
		Digits:  a.Digits,
		UserID:  user.ID,
		Icon:    a.Icon,
	}

	err = acc.SetSecret(a.Secret)
	if err != nil {
		log.Println(err)
		return h.RenderError(c, "new", http.StatusInternalServerError, "Unable to add the account")
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
