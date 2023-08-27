package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"

	"github.com/AlecAivazis/survey/v2"

	"github.com/kmjayadeep/totpm/pkg/configfile"
	"github.com/kmjayadeep/totpm/pkg/types"
)

func emailValidator(val interface{}) error {
	if str, ok := val.(string); ok {
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(str) {
			return fmt.Errorf("Please enter a valid email address")
		}
	} else {
		// otherwise we cannot convert the value into a string and cannot enforce length
		return fmt.Errorf("Invalid value")
	}

	return nil
}

func (c *Cli) authLogin(email, pass *string) {
	u := *c.server
	u.Path = path.Join(u.Path, "api/auth/login")

	if *email == "" {
		prompt := &survey.Input{
			Message: "Email Address",
		}
		survey.AskOne(prompt, email, survey.WithValidator(emailValidator))
	}

	if *pass == "" {
		prompt := &survey.Password{
			Message: "Enter Password",
		}
		survey.AskOne(prompt, pass, survey.WithValidator(survey.Required))
	}

	in := types.AuthInput{
		Email:    *email,
		Password: *pass,
	}

	d, err := json.Marshal(in)
	c.handleError(err)

	res, err := http.Post(u.String(), "application/json", bytes.NewBuffer(d))
	c.handleError(err)

	if res.StatusCode == http.StatusUnauthorized {
		c.handleError(fmt.Errorf("Invalid creds.. please login again"))
	}

	if res.StatusCode != http.StatusOK {
		c.handleError(fmt.Errorf("Unable to Login"))
	}

	var token struct {
		Token string
		Exp   uint
	}
	err = json.NewDecoder(res.Body).Decode(&token)
	c.handleError(err)

	cfg := &configfile.Config{
		Token: token.Token,
	}

	c.handleError(cfg.Write())

	c.PrintSuccess("Successfully logged In")
}
