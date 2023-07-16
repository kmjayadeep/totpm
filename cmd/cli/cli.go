package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/kmjayadeep/totpm/pkg/configfile"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/kmjayadeep/totpm/pkg/types"
	"github.com/pquerna/otp/totp"
)

var (
	app    = kingpin.New("totp", "Manage 2fa tokens")
	server = app.Flag("server", "Server address.").Default("http://localhost:3000").URL()

	login      = app.Command("login", "Login to the server")
	loginEmail = login.Flag("email", "Email address").String()
	loginPass  = login.Flag("password", "Password").String()

	otpList = app.Command("list", "List OTPs")
	otpCode = app.Command("code", "Show OTP token")
	otpSite = otpCode.Arg("name", "Name of the otp (optional)").String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// Register user
	case login.FullCommand():
		authLogin()

	case otpList.FullCommand():
		listOtps()

	case otpCode.FullCommand():
		getCode()
	}
}

func listOtps() {
	u := **server
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

func getCode() {
	u := **server
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
		if o.Name == *otpSite {
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

func authLogin() {
	u := **server
	u.Path = path.Join(u.Path, "api/auth/login")

	in := types.AuthInput{
		Email:    *loginEmail,
		Password: *loginPass,
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

	c := &configfile.Config{
		Token: token.AccessToken,
	}

	handleError(c.Write())
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
