package authtoken

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWTToken(userId uint, key string) (string, int64, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
		Subject:   strconv.Itoa(int(userId)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, &claims)

	signed, err := token.SignedString([]byte(key))
	return signed, claims.ExpiresAt, err
}

func ParseJWTToken(tokenString, key string) (uint, error) {
	claims := &jwt.StandardClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return 0, nil
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, nil
	}

	return uint(userId), nil
}
