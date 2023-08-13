package render

import (
	"log"
	"net/http"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
)

type AccountRequest struct {
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

	a := AccountRequest{}

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

func (h *Render) RenderEditAccount(c *fiber.Ctx) error {
	user, err := h.GetLoggedInUser(c)
	if err != nil {
		return err
	}

	id := c.Params("id")
	acc := &data.Account{}
	if res := h.db.Where("id=? and user_id=?", id, user.ID).First(acc); res.Error != nil {
		return res.Error
	}

	secret, err := acc.GetSecret()
	if err != nil {
		return err
	}

	if c.Method() == http.MethodGet {
		return c.Render("edit", fiber.Map{
			"account": acc,
			"secret":  secret,
		})
	}

	a := AccountRequest{}

	if err := c.BodyParser(&a); err != nil {
		return err
	}

	acc.Service = a.Service
	acc.Account = a.Account
	acc.Digits = a.Digits
	acc.OtpType = data.OtpType(a.OtpType)
	acc.Icon = a.Icon

	err = acc.SetSecret(a.Secret)
	if err != nil {
		log.Println(err)
		return h.RenderError(c, "edit", http.StatusInternalServerError, "Unable to edit the account")
	}

	if res := h.db.UpdateColumns(acc); res.Error != nil {
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

	sort.Slice(accs, func(i, j int) bool {
		return accs[i].ID < accs[j].ID
	})

	return c.Render("accounts", fiber.Map{
		"accounts": accs,
	})
}

func (h *Render) RenderDeleteAccount(c *fiber.Ctx) error {
	user, err := h.GetLoggedInUser(c)
	if err != nil {
		return err
	}

	id := c.Params("id")
	acc := &data.Account{}
	if res := h.db.Where("id=? and user_id=?", id, user.ID).First(acc); res.Error != nil {
		return res.Error
	}

	if res := h.db.Delete(acc); res.Error != nil {
		return res.Error
	}

	c.Append("HX-Redirect", "/accounts")

	return c.SendStatus(http.StatusNoContent)
}
