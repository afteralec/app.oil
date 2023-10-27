package sessiondata

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func New(s *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := s.Get(c)
		if err != nil {
			log.Print(err)
			return c.Next()
		}

		pid := sess.Get("pid")
		c.Locals("pid", pid)
		return c.Next()
	}
}
