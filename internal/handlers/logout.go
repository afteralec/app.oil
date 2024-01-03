package handlers

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/views"
)

func Logout(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		sess.Destroy()

		c.Append("HX-Redirect", routes.Logout)
		return nil
	}
}

func LogoutPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.Logout, c.Locals("b"), layouts.Standalone)
	}
}
