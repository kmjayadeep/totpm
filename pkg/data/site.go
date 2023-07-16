package data

import "gorm.io/gorm"

type Site struct {
	gorm.Model
	UserID     string
	Name       string
	Secret     string
	OtpAuthUri string
	Logo       string
}
