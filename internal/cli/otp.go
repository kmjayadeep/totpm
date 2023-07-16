package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/kmjayadeep/totpm/pkg/configfile"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/pquerna/otp/totp"
)

func (c *Cli) listOtps() {
	u := *c.server
	u.Path = path.Join(u.Path, "api/site")

	config, err := configfile.Read()
	handleError(err)

	client := &http.Client{}

	req, err := http.NewRequest("GET", u.String(), nil)
	handleError(err)

	req.Header.Set("x-access-token", config.Token)

	res, err := client.Do(req)
	handleError(err)

	if res.StatusCode == http.StatusUnauthorized {
		panic("Invalid token.. please login again")
	}

	if res.StatusCode != http.StatusOK {
		panic("Unable to list OTPs")
	}

	var sites []data.Site
	if err := json.NewDecoder(res.Body).Decode(&sites); err != nil {
		panic(err.Error() + "Unable to list OTPs")
	}

	for _, o := range sites {
		fmt.Printf("%d\t%s\n", o.ID, o.Name)
	}
}

func (c *Cli) getCode(site *string) {
	u := *c.server
	u.Path = path.Join(u.Path, "api/site")
	res, err := http.Get(u.String())
	if err != nil {
		panic(err.Error() + "Unable to get OTP")
	}

	if res.StatusCode == http.StatusUnauthorized {
		panic("Invalid token.. please login again")
	}

	if res.StatusCode != http.StatusOK {
		panic("Unable to get OTP")
	}

	var sites []data.Site
	if err := json.NewDecoder(res.Body).Decode(&sites); err != nil {
		panic(err.Error() + "Unable to get OTP")
	}

	for _, o := range sites {
		if o.Name == *site {
			code, err := totp.GenerateCode(o.Secret, time.Now())
			if err != nil {
				panic(err.Error() + "Unable to get OTP")
			}
			fmt.Println(code)
			return
		}
	}

	panic("Unable to find the given OTP")
}
