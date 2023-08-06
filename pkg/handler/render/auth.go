package render

import (
	"encoding/gob"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type RegisterRequest struct {
	Email           string `form:"email"`
	Password        string `form:"password"`
	PasswordConfirm string `form:"passwordConfirm"`
}

func init() {
	gob.Register(&data.User{})
}

func (h *Render) RenderRegister(c *fiber.Ctx) error {
	if c.Method() == http.MethodGet {
		return c.Render("register", fiber.Map{})
	}

	l := RegisterRequest{}

	if err := c.BodyParser(&l); err != nil {
		return err
	}

	if l.Password != l.PasswordConfirm {
		return c.Render("register", fiber.Map{
			"error": "Passwords doesnt match",
		})
	}
	user := &data.User{
		Email: l.Email,
	}

	count := int64(0)

	if res := h.db.Model(&user).Where("email=?", l.Email).Count(&count); res.Error != nil {
		return res.Error
	}

	if count > 0 {
		return c.Render("register", fiber.Map{
			"error": "Email already registered, please login",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(l.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)

	if res := h.db.Save(user); res.Error != nil {
		return res.Error
	}

	return c.Render("register", fiber.Map{
		"message": "Registration successful, please login",
	})
}

func (h *Render) RenderLogin(c *fiber.Ctx) error {
	if c.Method() == http.MethodGet {
		return c.Render("login", fiber.Map{})
	}

	l := LoginRequest{}

	if err := c.BodyParser(&l); err != nil {
		return err
	}

	user := &data.User{}

	if res := h.db.Where("email=?", l.Email).First(user); res.Error != nil {
		return res.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(l.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).Render("login", fiber.Map{
			"error": "Invalid login credentials",
		})
	}

	sess, err := h.store.Get(c)
	if err != nil {
		return err
	}
	sess.Set("user", user)

	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/accounts")
}
