package handlers

import (
	fiber "github.com/gofiber/fiber/v2"
)

func HomePage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/index", c.Locals("bind"))
	}
}
