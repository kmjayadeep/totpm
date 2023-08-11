package security_test

import (
	"testing"

	"github.com/kmjayadeep/totpm/internal/security"
)

func TestEncryption(t *testing.T) {

	key := "IPIkJ1AlgmP7pDbksjuiAiQqG6HRlYnp"
	secret := "LQVKEEKDI3ELVSVDMVQGJNJFJRN4ZNCOB2JBSLJ55BH3OIHJZHGCWBJJZ5YERRUSUITS6NTDLGLZ436LQVFK5YN6QJB66UKF4BJWUFA"

	enc, err := security.Encrypt(key, secret)

	if err != nil {
		t.Fatalf("unable to encrypt")
	}

	secret2, err := security.Decrypt(key, enc)
	if err != nil {
		t.Fatalf("unable to Decrypt")
	}

	if secret != secret2 {
		t.Fatalf("got different value when decrypting, %s, %s, %s", secret, enc, secret2)
	}
}
