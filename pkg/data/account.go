package data

import "gorm.io/gorm"

type OtpType string

const (
	OtpTypeTOTP OtpType = "totp"
	OtpTypeHOTP OtpType = "hotp"
)

type Account struct {
	gorm.Model
	UserID    string
	Service   string
	Account   string
	Icon      string
	OtpType   OtpType
	Digits    uint
	Algorithm string
	Period    uint
	Counter   uint64

	Secret string
}
