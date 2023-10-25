package handlers

import (
	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/middleware"
)

func Logout(c *fiber.Ctx) error {
	sess, err := middleware.Sessions.Get(c)
	if err != nil {
		c.Status(500)
		return nil
	}

	sess.Destroy()

	c.Status(200)
	return nil
}
