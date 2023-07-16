package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"

	"github.com/kmjayadeep/totpm/pkg/configfile"
	"github.com/kmjayadeep/totpm/pkg/types"
)

func (c *Cli) authLogin(email, pass *string) {
	u := *c.server
	u.Path = path.Join(u.Path, "api/auth/login")

	in := types.AuthInput{
		Email:    *email,
		Password: *pass,
	}

	d, err := json.Marshal(in)
	handleError(err)

	res, err := http.Post(u.String(), "application/json", bytes.NewBuffer(d))
	handleError(err)

	if res.StatusCode == http.StatusUnauthorized {
		panic("Invalid creds.. please login again")
	}

	if res.StatusCode != http.StatusOK {
		panic("Unable to Login")
	}

	var token struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   uint   `json:"expires_in"`
	}
	err = json.NewDecoder(res.Body).Decode(&token)
	handleError(err)

	cfg := &configfile.Config{
		Token: token.AccessToken,
	}

	handleError(cfg.Write())
}
