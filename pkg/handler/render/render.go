package render

import "gorm.io/gorm"

type Render struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Render {
	return &Render{
		db: db,
	}
}
