package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/kmjayadeep/totpm/internal/config"
	"github.com/kmjayadeep/totpm/pkg/data"
	"github.com/kmjayadeep/totpm/pkg/handler"
	apihandler "github.com/kmjayadeep/totpm/pkg/handler/api"
	supa "github.com/nedpals/supabase-go"
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
	db.AutoMigrate(&data.Site{}, &data.Account{})
	supabase := supa.CreateClient(config.Get().SupabaseURL, config.Get().SupabaseKey)

	h := handler.NewHandler(db, supabase)
	api := apihandler.NewAPI(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Static("/assets", "./assets/dist")

	// Render pages
	app.Get("/home", h.RenderSites)
	app.Get("/home/:id", h.RenderSite)

	// APIs
	app.Post("/api/auth/signup", h.Signup)
	app.Post("/api/auth/login", h.Login)

	app.Get("/api/site", h.RequiresAuth, h.GetSites)
	app.Get("/api/site/:id", h.RequiresAuth, h.GetSite)
	app.Delete("/api/site/:id", h.RequiresAuth, h.DeleteSite)
	app.Post("/api/site", h.RequiresAuth, h.AddSite)

	// APIS compatible with 2fauth
	app.Get("/api/v1/twofaccounts", api.GetAccounts)
	app.Post("/api/v1/twofaccounts", api.AddAccount)
	app.Delete("/api/v1/twofaccounts", api.DeleteAccounts)

	app.Get("/api/v1/twofaccounts/:id/otp", api.GetAccountOTP)

	log.Fatal(app.Listen(":3000"))
}
