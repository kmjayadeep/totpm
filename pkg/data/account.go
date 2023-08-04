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
	Period    uint
	Counter   uint64

	Secret string
}
