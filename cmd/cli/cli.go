package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/pquerna/otp/totp"
)

var (
	app    = kingpin.New("totp", "Manage 2fa tokens")
	server = app.Flag("server", "Server address.").Default("http://localhost:3000").URL()

	login      = app.Command("login", "Login to the server")
	loginToken = login.Flag("token", "API token").String()

	otpList = app.Command("list", "List OTPs")
	otpCode = app.Command("code", "Show OTP token")
	otpSite = otpCode.Arg("name", "Name of the otp (optional)").String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// Register user
	case login.FullCommand():
		println(*loginToken)

	case otpList.FullCommand():
		listOtps()

	case otpCode.FullCommand():
		getCode()
	}
}

func listOtps() {
	u := **server
	u.Path = path.Join(u.Path, "api/site")
	res, err := http.Get(u.String())
	if err != nil {
		panic(err.Error() + "Unable to list OTPs")
	}

	if res.StatusCode == http.StatusUnauthorized {
		panic("Invalid token.. please login again")
	}

	if res.StatusCode != http.StatusOK {
		panic(err.Error() + "Unable to list OTPs")
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
		panic(err.Error() + "Unable to get OTP")
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
