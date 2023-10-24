package handlers

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.Render("web/views/index", fiber.Map{
		"CopyrightYear": time.Now().Year(),
		"Title":         "Hello, World!",
	}, "web/views/layouts/main")
}
