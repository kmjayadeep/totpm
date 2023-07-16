package configfile

import (
	"encoding/json"
	"os"
	"path"

	"github.com/kmjayadeep/totpm/pkg/types"
)

type Cache struct {
	Otps []types.Otp
}

func (c *Cache) Write() error {
	h, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	p := path.Join(h, "totp-cache")

	f, err := os.Create(p)
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(c)
}

func NewCache() *Cache {
	return &Cache{
		Otps: make([]types.Otp, 0),
	}
}

func ReadCache() (*Cache, error) {
	h, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	p := path.Join(h, "totp-cache")

	f, err := os.Open(p)

	if os.IsNotExist(err) {
		return &Cache{}, nil
	}

	if err != nil {
		return nil, err
	}

	c := &Cache{}

	if err := json.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}
