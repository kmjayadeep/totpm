package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/kmjayadeep/totpm/internal/config"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/kmjayadeep/totpm/pkg/handler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	db, err := gorm.Open(postgres.Open(config.Get().DBConnectionString), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&data.Site{})

	h := handler.NewHandler(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Static("/assets", "./assets/dist")

	// Render pages
	app.Get("/home", h.RenderSites)
	app.Get("/home/:id", h.RenderSite)

	// APIs
	app.Get("/api/site", h.GetSites)
	app.Post("/api/site", h.AddSite)

	log.Fatal(app.Listen(":3000"))
}
