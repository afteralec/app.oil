package session

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/interfaces"
)

const TwoHoursInSeconds = 120 * 60

func New(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := i.Sessions.Get(c)
		if err != nil {
			return c.Next()
		}

		pid := sess.Get("pid")
		if pid != nil {
			c.Locals("pid", pid)
		}

		theme := sess.Get("theme")
		if theme == nil {
			c.Locals("theme", constants.ThemeDefault)
		} else {
			c.Locals("theme", theme)
		}

		return c.Next()
	}
}
