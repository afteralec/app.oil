package home

import (
	fiber "github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/index", c.Locals("bind"))
	}
}
