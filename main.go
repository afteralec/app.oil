package main

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/queries"
)

func main() {
	queries.Connect()
	queries.Build()
	config := configs.Fiber()
	app := fiber.New(config)

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

	app.Static("/", "./web/static")

	app.Listen(":8008")
}
