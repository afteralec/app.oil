package handlers

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/views"
)

func RecoverPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("views/recover", c.Locals(constants.BindName), views.LayoutStandalone)
	}
}
