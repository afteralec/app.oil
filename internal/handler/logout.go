package handler

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/views"
)

func Logout(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		sess.Destroy()

		c.Append("HX-Redirect", route.Logout)
		return nil
	}
}

func LogoutPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.Logout, views.Bind(c), layout.Standalone)
	}
}
