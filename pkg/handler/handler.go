package handler

import (
	supa "github.com/nedpals/supabase-go"
	"gorm.io/gorm"
)

type Handler struct {
	db       *gorm.DB
	supabase *supa.Client
}

func NewHandler(db *gorm.DB, sc *supa.Client) *Handler {
	return &Handler{
		db:       db,
		supabase: sc,
	}
}
