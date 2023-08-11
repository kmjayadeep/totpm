package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func Decrypt(key, v string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext, _ := hex.DecodeString(v)
	if (len(ciphertext) % aes.BlockSize) > 0 {
		return "", fmt.Errorf("input length is not a multiple of blocksize ")
	}

	iv := ciphertext[:aes.BlockSize]
	cbc := cipher.NewCBCDecrypter(c, iv)
	ciphertext = ciphertext[len(iv):]

	// decrypt in-place
	cbc.CryptBlocks(ciphertext, ciphertext)
	unpadded, err := pkcsUnpad(ciphertext, aes.BlockSize)
	if err != nil {
		return "", err
	}
	return string(unpadded), err
}

func Encrypt(key, v string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plaintext := []byte(v)

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	padded := pkcsPad(plaintext, aes.BlockSize)

	// put the IV at the beginning of the ciphertext
	encrypted := make([]byte, len(iv)+len(padded))
	copy(encrypted[:len(iv)], iv)
	cbc.CryptBlocks(encrypted[len(iv):], padded)

	return hex.EncodeToString(encrypted), nil
}

// very simple PKCS padding, as implemented in pgcrypto
func pkcsPad(input []byte, blockSize int) []byte {
	padLen := blockSize - (len(input) % blockSize)
	padded := make([]byte, len(input)+padLen)
	copy(padded, input)

	padding := padded[len(input):]
	for i := range padding {
		padding[i] = byte(padLen)
	}
	return padded
}

// .. and the reverse operation
func pkcsUnpad(input []byte, blockSize int) ([]byte, error) {
	if len(input)%blockSize != 0 {
		return nil, fmt.Errorf("input length %d not divisible by block size %d", len(input), blockSize)
	}
	if len(input) < blockSize {
		return nil, fmt.Errorf("input length %d is smaller than block size %d", len(input), blockSize)
	}
	padLen := int(input[len(input)-1])
	if padLen <= 0 || padLen > blockSize {
		return nil, fmt.Errorf("invalid padding length %d", padLen)
	}
	for pos, byte := range input[len(input)-padLen:] {
		if int(byte) != padLen {
			return nil, fmt.Errorf("padding byte %d at pos %d is not the same as padding length %d", byte, pos, padLen)
		}
	}
	return input[:len(input)-padLen], nil
}
