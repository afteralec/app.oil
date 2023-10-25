package middleware

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func Session() fiber.Handler {
	Store = session.New()
	return func(c *fiber.Ctx) error {
		sess, err := Store.Get(c)
		if err != nil {
			log.Print(err)
			return c.Next()
		}
		pid := sess.Get("pid")
		c.Locals("pid", pid)
		return c.Next()
	}
}
