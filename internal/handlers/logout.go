package handlers

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Logout(s *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := s.Get(c)
		if err != nil {
			c.Status(500)
			return nil
		}

		sess.Destroy()

		c.Status(200)
		return nil
	}
}
