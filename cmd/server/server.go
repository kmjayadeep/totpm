package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/xlzd/gotp"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Static("/assets", "./assets/dist")

	app.Get("/home", func(c *fiber.Ctx) error {

		sites := []data.Site{{
			Name:   "amazon",
			Secret: "koo",
		}}

		return c.Render("home", fiber.Map{
			"sites": sites,
		})
	})

	app.Get("/home/:id", func(c *fiber.Ctx) error {

		_ = c.Params("id", "0")

		sites := []data.Site{{
			Name:   "amazon",
			Secret: "MFRGGCQ=",
		}}

		current := sites[0]

		totp := gotp.NewDefaultTOTP(current.Secret)
		code, exp := totp.NowWithExpiration()

		return c.Render("home", fiber.Map{
			"sites":   sites,
			"current": current,
			"code":    code,
			"exp":     exp,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
