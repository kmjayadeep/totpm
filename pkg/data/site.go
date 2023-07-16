package data

import "gorm.io/gorm"

type Site struct {
	gorm.Model
	UserID  string
	Logo    string
	Website string
	Name    string
	Secret  string
}
