package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

type Cli struct {
	server **url.URL
	debug  *bool
}

func NewCli() *Cli {
	return &Cli{}
}

func (c *Cli) Run() {
	app := kingpin.New("totp", "Manage 2fa tokens")
	c.server = app.Flag("server", "Server address.").Default("http://localhost:3000").URL()
	c.debug = app.Flag("debug", "Enable debug mode").Default("false").Bool()

	login := app.Command("login", "Login to the server")
	loginEmail := login.Flag("email", "Email address").String()
	loginPass := login.Flag("password", "Password").String()

	otpList := app.Command("list", "List OTPs")

	otpCode := app.Command("code", "Show OTP token")
	otpName := otpCode.Arg("name", "Name/ID of the otp").Required().String()

	otpDelete := app.Command("delete", "Delete OTP token")
	otpDeleteName := otpDelete.Arg("name", "Name/ID of the otp").Required().String()

	otpAdd := app.Command("add", "Add new totp")
	otpAddUri := otpAdd.Flag("uri", "otpauth:// Uri").String()
	otpAddName := otpAdd.Flag("name", "Identifier for the otp token").String()
	otpAddSecret := otpAdd.Flag("secret", "Secret value").String()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// Register user
	case login.FullCommand():
		c.authLogin(loginEmail, loginPass)

	case otpList.FullCommand():
		c.listOtps()

	case otpCode.FullCommand():
		c.getCode(otpName)

	case otpDelete.FullCommand():
		c.delete(otpDeleteName)

	case otpAdd.FullCommand():
		c.addOtp(otpAddUri, otpAddName, otpAddSecret)
	}
}

func (c *Cli) handleError(err error) {
	if err != nil {
		if *c.debug {
			fmt.Println("Error: " + err.Error())
		}
		fmt.Println("Unexpected error; Use --debug flag to view more details")
		os.Exit(1)
	}
}

func (c *Cli) handleErrorMsg(err error, msg string) {
	if err != nil {
		if *c.debug {
			fmt.Println("Error: " + err.Error())
		}
		fmt.Println("Error : " + msg + "; Use --debug flag to view more details")
		os.Exit(1)
	}
}

func (c *Cli) handleErrorBody(body io.ReadCloser, msg string) {
	if *c.debug {
		b, err := ioutil.ReadAll(body)
		c.handleError(err)
		fmt.Println("Body: " + string(b))
	}
	fmt.Println("Error : " + msg + "; Use --debug flag to view more details")
	os.Exit(1)
}
