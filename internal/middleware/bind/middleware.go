package bind

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constant"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		b := fiber.Map{
			"PID": c.Locals("pid"),
		}
		c.Locals(constant.BindName, b)
		return c.Next()
	}
}
