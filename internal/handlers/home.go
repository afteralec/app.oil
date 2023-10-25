package handlers

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.Render("web/views/index", fiber.Map{
		"PID":           c.Locals("pid"),
		"CopyrightYear": time.Now().Year(),
		"Title":         "Petrichor",
		"MetaContent":   "Petrichor MUD - a modern take on a classic MUD style of game.",
	})
}
