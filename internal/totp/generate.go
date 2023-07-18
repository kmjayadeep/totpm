package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strings"
	"time"
)

// GenerateCode generates the TOTP based on the given secret and time
// Returns code, expire time and error
func GenerateCode(secret string, t time.Time) (string, time.Time, error) {

	// Convert the secret to bytes
	decodedKey, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	// Calculate the number of time steps
	timeStep := int64(30)
	counter := t.Unix() / timeStep
	nextCount := counter + 1 // for expire date

	// Convert the counter to bytes
	msg := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		msg[i] = byte(counter & 0xFF)
		counter = counter >> 8
	}

	// Generate the HMAC-SHA1
	hmacHash := hmac.New(sha1.New, decodedKey)
	_, err = hmacHash.Write(msg)
	if err != nil {
		return "", time.Time{}, err
	}
	hash := hmacHash.Sum(nil)

	// Truncate the hash to get the offset
	offset := hash[len(hash)-1] & 0x0F

	// Get the 4 bytes dynamic code
	dynamicCode := (int(hash[offset]&0x7F)<<24 |
		int(hash[offset+1]&0xFF)<<16 |
		int(hash[offset+2]&0xFF)<<8 |
		int(hash[offset+3]&0xFF))

	// Calculate the TOTP value
	digits := 6
	modulo := 1
	for i := 0; i < digits; i++ {
		modulo *= 10
	}
	otpValue := fmt.Sprintf("%0*d", digits, dynamicCode%modulo)

	// Calculate the remaining seconds until the TOTP code expires.
	expireTime := time.Unix(nextCount*timeStep, 0)

	return otpValue, expireTime, nil
}

func ValidateSecretFormat(s string) error {
	_, err := base32.StdEncoding.DecodeString(strings.ToUpper(s))
	return err
}
