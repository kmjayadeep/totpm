package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kmjayadeep/totpm/internal/totp"
	"github.com/kmjayadeep/totpm/pkg/configfile"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/kmjayadeep/totpm/pkg/types"
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

	if err := totp.ValidateSecretFormat(*secret); err != nil {
		c.handleErrorMsg(err, "Invalid totp secret")
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
	c.refreshCache()
	u := **c.server

	config, err := configfile.Read()
	c.handleError(err)
	cache, err := configfile.ReadCache()
	c.handleError(err)

	id := uint(0)

	idx, err := strconv.Atoi(*site)
	if err == nil {
		if idx > len(cache.Otps) {
			c.handleError(fmt.Errorf("Invalid otp id"))
		}
		id = cache.Otps[idx-1].ID
	} else {
		for _, s := range cache.Otps {
			if s.Name == *site {
				id = s.ID
				break
			}
		}
	}
	if id == 0 {
		c.handleError(fmt.Errorf("Invalid otp name"))
	}

	u.Path = path.Join(u.Path, "api/site/", strconv.Itoa(int(id)))
	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	c.handleError(err)

	req.Header.Set("x-access-token", config.Token)
	res, err := client.Do(req)
	c.handleError(err)

	if res.StatusCode == http.StatusUnauthorized {
		c.handleError(fmt.Errorf("Invalid token.. please login again"))
	}

	if res.StatusCode != http.StatusOK {
		c.handleErrorBody(res.Body, "Unable to get OTP")
	}

	otp := data.Site{}
	if err := json.NewDecoder(res.Body).Decode(&otp); err != nil {
		c.handleError(fmt.Errorf("Unable to get OTP"))
	}

	code, _, err := totp.GenerateCode(otp.Secret, time.Now())
	c.handleError(err)

	fmt.Println(code)
}

func (c *Cli) delete(site *string) {
	c.refreshCache()
	u := **c.server

	config, err := configfile.Read()
	c.handleError(err)
	cache, err := configfile.ReadCache()
	c.handleError(err)

	id := uint(0)

	idx, err := strconv.Atoi(*site)
	if err == nil {
		if idx > len(cache.Otps) {
			c.handleError(fmt.Errorf("Invalid otp id"))
		}
		id = cache.Otps[idx-1].ID
	} else {
		for _, s := range cache.Otps {
			if s.Name == *site {
				id = s.ID
				break
			}
		}
	}
	if id == 0 {
		c.handleError(fmt.Errorf("Invalid otp name"))
	}

	u.Path = path.Join(u.Path, "api/site/", strconv.Itoa(int(id)))
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	c.handleError(err)

	req.Header.Set("x-access-token", config.Token)
	res, err := client.Do(req)
	c.handleError(err)

	if res.StatusCode == http.StatusUnauthorized {
		c.handleError(fmt.Errorf("Invalid token.. please login again"))
	}

	if res.StatusCode != http.StatusOK {
		c.handleErrorBody(res.Body, "Unable to delete OTP")
	}

	c.PrintSuccess("OTP Deleted successfully")
}
