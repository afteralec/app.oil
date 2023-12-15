package permissions

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/shared"
)

func New(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			return c.Next()
		}

		// TODO: Cache this in Redis here and on login
		// TODO: This will require an API for updating the value in Redis on change
		ps, err := i.Queries.ListPlayerPermissions(context.Background(), pid.(int64))
		if err != nil {
			return c.Next()
		}
		if len(ps) == 0 {
			return c.Next()
		}

		perms := permission.MakePlayerIssued(pid.(int64), ps)

		c.Locals("perms", perms)
		return c.Next()
	}
}
