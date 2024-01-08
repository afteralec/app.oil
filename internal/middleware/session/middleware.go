package session

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/shared"
)

const TwoHoursInSeconds = 120 * 60

func New(i *shared.Interfaces) fiber.Handler {
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
			c.Locals("theme", constants.ThemeLight)
		} else {
			c.Locals("theme", theme)
		}

		return c.Next()
	}
}
