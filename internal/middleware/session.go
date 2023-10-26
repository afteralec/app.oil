package middleware

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Sessions *session.Store

func Session() fiber.Handler {
	// TODO: Update this config to be more secure. Will depend on environment.
	Sessions = session.New()
	return func(c *fiber.Ctx) error {
		sess, err := Sessions.Get(c)
		if err != nil {
			log.Print(err)
			return c.Next()
		}
		pid := sess.Get("pid")
		c.Locals("pid", pid)
		return c.Next()
	}
}
