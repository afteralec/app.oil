package handlers

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/views"
)

func RecoverPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.Recover, c.Locals(constants.BindName), layouts.Standalone)
	}
}
