package render

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/kmjayadeep/totpm/pkg/data"
	"gorm.io/gorm"
)

type Render struct {
	db    *gorm.DB
	store *session.Store
}

func NewHandler(db *gorm.DB, s *session.Store) *Render {
	return &Render{
		db:    db,
		store: s,
	}
}

func (h *Render) GetLoggedInUser(c *fiber.Ctx) (*data.User, error) {
	if user, ok := c.Locals("user").(*data.User); ok {
		return user, nil
	}

	return nil, fmt.Errorf("Invalid data in the session")
}

func (h *Render) GetUserFromSession(c *fiber.Ctx) (*data.User, error) {
	sess, err := h.store.Get(c)
	if err != nil {
		return nil, err
	}
	u := sess.Get("user")
	if user, ok := u.(*data.User); ok {
		return user, nil
	}
	return nil, fmt.Errorf("Invalid data in the session")
}

func (h *Render) RequireAuth(c *fiber.Ctx) error {
	u, err := h.GetUserFromSession(c)
	if err != nil {
		return c.Redirect("/login")
	}

	user := &data.User{}

	if res := h.db.First(&user, u.ID); res.Error != nil {
		return c.Redirect("/login")
	}

	c.Locals("user", user)

	return c.Next()
}

func (h *Render) RenderError(c *fiber.Ctx, template string, status int, msg string) error {
	return c.Status(status).Render(template, fiber.Map{
		"error": msg,
	})
}
