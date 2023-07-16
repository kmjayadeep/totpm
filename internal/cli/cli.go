package cli

import (
	"net/url"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

type Cli struct {
	server **url.URL
}

func NewCli() *Cli {
	return &Cli{}
}

func (c *Cli) Run() {
	app := kingpin.New("totp", "Manage 2fa tokens")
	c.server = app.Flag("server", "Server address.").Default("http://localhost:3000").URL()

	login := app.Command("login", "Login to the server")
	loginEmail := login.Flag("email", "Email address").String()
	loginPass := login.Flag("password", "Password").String()

	otpList := app.Command("list", "List OTPs")
	otpCode := app.Command("code", "Show OTP token")
	otpName := otpCode.Arg("name", "Name of the otp (optional)").String()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// Register user
	case login.FullCommand():
		c.authLogin(loginEmail, loginPass)

	case otpList.FullCommand():
		c.listOtps()

	case otpCode.FullCommand():
		c.getCode(otpName)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
