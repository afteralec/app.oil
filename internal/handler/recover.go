package handler

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/views"
)

func RecoverPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.Recover, views.Bind(c), layout.Standalone)
	}
}
