package main

import (
	"time"
  "embed"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/queries"
  "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	queries.Connect()
	queries.Build()
	config := configs.Fiber(viewsfs)
	app := fiber.New(config)

  app.Use(cors.New())
  app.Use(logger.New())

  app.Static("/", "./web/static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("web/views/index", fiber.Map{
			"CopyrightYear": time.Now().Year(),
			"Title":         "Hello, World!",
		}, "web/views/layouts/main")
	})

	app.Get("/request/:id", func(c *fiber.Ctx) error {
		return c.Render("web/views/request", fiber.Map{
			"CopyrightYear":    time.Now().Year(),
			"ID":               c.Params("id"),
			"Status":           "Ready",
			"Name":             "Test Character",
			"Backstory":        "This is a tragic backstory.\nWith a newline.",
			"ShortDescription": "test, testerly man",
			"Description":      "This is a test description.",
			"Class":            "Crafting",
			"Origin":           "LowQuarter",
			"Gender":           "Male",
		}, "web/views/layouts/main")
	})


	app.Listen(":8008")
}
