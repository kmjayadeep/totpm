package otp

import (
	"fmt"
	"strings"
	"time"

	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

func algorithm(algo string) otp.Algorithm {
	a := strings.ToLower(algo)
	switch a {
	case "md5":
		return otp.AlgorithmMD5
	case "sha256":
		return otp.AlgorithmSHA256
	case "sha512":
		return otp.AlgorithmSHA512
	default:
		return otp.AlgorithmSHA1
	}
}

// TODO return expire date
func GenerateCode(a data.Account) (string, error) {
	sec, err := a.GetSecret()
	if err != nil {
		return "", err
	}

	if a.OtpType == data.OtpTypeTOTP {
		return totp.GenerateCodeCustom(sec, time.Now(), totp.ValidateOpts{
			Period:    a.Period,
			Digits:    otp.Digits(a.Digits),
			Algorithm: algorithm(a.Algorithm),
		})
	}

	if a.OtpType == data.OtpTypeHOTP {
		return hotp.GenerateCodeCustom(sec, a.Counter, hotp.ValidateOpts{
			Digits:    otp.Digits(a.Digits),
			Algorithm: algorithm(a.Algorithm),
		})
	}

	return "", fmt.Errorf("Invalid otp type")
}
