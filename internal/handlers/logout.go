package handlers

import (
	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/shared"
)

func Logout(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Status(500)
			return nil
		}

		sess.Destroy()

		c.Append("HX-Redirect", "/")
		return nil
	}
}
