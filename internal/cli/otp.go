package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kmjayadeep/totpm/pkg/configfile"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/kmjayadeep/totpm/pkg/types"
	"github.com/pquerna/otp/totp"
)

func (c *Cli) listOtps() {
	c.refreshCache()

	cache, err := configfile.ReadCache()
	c.handleError(err)

	fmt.Println("Id\tName\n============")

	for i, o := range cache.Otps {
		fmt.Printf("%d\t%s\n", i+1, o.Name)
	}
}

func (c *Cli) addOtp(uri, name, secret *string) {
	u := **c.server
	u.Path = path.Join(u.Path, "api/site")

	config, err := configfile.Read()
	c.handleError(err)

	if *name == "" {
		prompt := &survey.Input{
			Message: "Name",
		}
		survey.AskOne(prompt, name, survey.WithValidator(survey.Required))
	}

	if *secret == "" {
		prompt := &survey.Password{
			Message: "Secret",
		}
		survey.AskOne(prompt, secret, survey.WithValidator(survey.Required))
	}

	in := types.OtpInput{
		OtpAuthUri: *uri,
		Name:       *name,
		Secret:     *secret,
	}
	d, err := json.Marshal(&in)
	c.handleError(err)

	client := &http.Client{}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(d))
	c.handleError(err)

	req.Header.Set("x-access-token", config.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	c.handleError(err)

	if res.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		c.handleError(err)
		c.handleError(fmt.Errorf("Unable to add new totp; server response: " + string(b)))
	}

	c.refreshCache()

	c.PrintSuccess("Otp added successfully!")
}

func (c *Cli) getCode(site *string) {
	u := **c.server
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
