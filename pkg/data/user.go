package data

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string
	PasswordHash string
}
