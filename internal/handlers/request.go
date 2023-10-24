package handlers

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

func Request(c *fiber.Ctx) error {
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
}
