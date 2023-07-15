package data

import "gorm.io/gorm"

type Site struct {
	gorm.Model
	Logo    string
	Website string
	Name    string
	Secret  string
}
