package configfile

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	Token string
}

func (c *Config) Write() error {
	h, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	p := path.Join(h, ".totp")

	f, err := os.Create(p)
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(c)
}

func Read() (*Config, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	p := path.Join(h, ".totp")

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	c := &Config{}

	if err := json.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}
