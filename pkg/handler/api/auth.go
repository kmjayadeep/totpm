package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/internal/authtoken"
	"github.com/kmjayadeep/totpm/internal/config"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/kmjayadeep/totpm/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthInput struct {
	Email    string
	Password string
}

func (h *Api) Login(c *fiber.Ctx) error {
	in := types.AuthInput{}
	if err := c.BodyParser(&in); err != nil {
		return err
	}

	user := &data.User{}

	if res := h.db.Where("email=?", in.Email).First(user); res.Error != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"Message": "Invalid credentails",
		})
	}

	token, exp, err := authtoken.GenerateJWTToken(user.ID, config.Get().AppKey)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"Token": token,
		"Exp":   exp,
	})
}
