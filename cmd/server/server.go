package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/kmjayadeep/totpm/pkg/data"
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

	app.Static("/assets", "./assets")

	app.Get("/home", func(c *fiber.Ctx) error {

		otps := []data.Totp{{
			Name:   "amazon",
			Secret: "koo",
		}}

		return c.Render("home", fiber.Map{
			"otps": otps,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
