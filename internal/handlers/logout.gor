package handlers

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
)

func Logout(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		sess.Destroy()

		c.Append("HX-Redirect", "/logout")
		return nil
	}
}

func LogoutPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/logout", c.Locals("b"), "web/views/layouts/standalone")
	}
}
