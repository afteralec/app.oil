package session

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/service"
)

const TwoHoursInSeconds = 120 * 60

func New(i *service.Interfaces) fiber.Handler {
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
			c.Locals("theme", constant.ThemeDefault)
		} else {
			c.Locals("theme", theme)
		}

		return c.Next()
	}
}
