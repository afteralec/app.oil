package handlers

import (
	fiber "github.com/gofiber/fiber/v2"
)

const RecoverRoute = "/recover"

func RecoverPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover", c.Locals("b"), "web/views/layouts/standalone")
	}
}
