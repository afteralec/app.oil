package handler

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/views"
)

// TODO: Add a main notification section to the main layout so we can notify the player
// TODO: i.e., if they have no email addresses set
func HomePage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.Home, views.Bind(c))
	}
}
