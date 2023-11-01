package sessiondata

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/shared"
)

const TwoHoursInSeconds = 120 * 60

func New(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := i.Sessions.Get(c)
		if err != nil {
			log.Print(err)
			return c.Next()
		}

		pid := sess.Get("pid")
		if pid != nil {
			c.Locals("pid", pid)

			perms, err := permissions.List(i, pid.(int64))
			if err != nil {
				return c.Next()
			}
			c.Locals("perms", perms)
		}

		return c.Next()
	}
}
