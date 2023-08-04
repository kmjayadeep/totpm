package api

import (
	"gorm.io/gorm"
)

type Api struct {
	db *gorm.DB
}

func NewAPI(db *gorm.DB) *Api {
	return &Api{
		db: db,
	}
}
