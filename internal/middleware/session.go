package middleware

import (
  "log"

  fiber "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func Session() fiber.Handler {
  store = session.New()
  return func(c *fiber.Ctx) error {
    sess, err := store.Get(c)
    if err != nil {
      log.Fatal(err)
    }
    pid := sess.Get("pid")
    c.Locals("pid", pid)
    return c.Next()
  }
}
