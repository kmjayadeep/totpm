package data

import (
	"github.com/kmjayadeep/totpm/internal/config"
	"github.com/kmjayadeep/totpm/internal/security"
	"gorm.io/gorm"
)

type OtpType string

const (
	OtpTypeTOTP OtpType = "totp"
	OtpTypeHOTP OtpType = "hotp"
)

type Account struct {
	gorm.Model
	UserID    uint
	Service   string
	Account   string
	Icon      string
	OtpType   OtpType
	Digits    uint
	Algorithm string
	Period    uint
	Counter   uint64

	SecretEncrypted string
}

func (a *Account) GetSecret() (string, error) {
	return security.Decrypt(config.Get().AppKey, a.SecretEncrypted)
}

func (a *Account) SetSecret(secret string) error {
	s, err := security.Encrypt(config.Get().AppKey, secret)
	if err != nil {
		return err
	}
	a.SecretEncrypted = s
	return nil
}
