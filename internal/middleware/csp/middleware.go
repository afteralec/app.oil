package csp

import (
	"os"

	fiber "github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		csp := "default-src 'self'"
		if os.Getenv("PETRICHOR_APP_ENV") == "prod" {
			csp = "default-src https://petrichormud.com https://play.petrichormud.com"
		}

		c.Locals("csp", csp)
		c.Append("Content-Security-Policy", csp)

		return c.Next()
	}
}
