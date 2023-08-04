package data

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID    string
	Service   string
	Account   string
	Icon      string
	OtpType   string
	Digits    int
	Algorithm string
	Period    int
	Counter   int

	Secret string
}
