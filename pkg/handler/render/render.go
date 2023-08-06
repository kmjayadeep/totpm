package render

import (
	"github.com/gofiber/fiber/v2/middleware/session"
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
