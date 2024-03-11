package permissions

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/service"
)

func New(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			return c.Next()
		}

		// TODO: Cache this in Redis here and on login
		// TODO: This will require an API for updating the value in Redis on change
		ps, err := i.Queries.ListPlayerPermissions(context.Background(), pid.(int64))
		if err != nil {
			// TODO: Split up the bind variables into different middleware
			// That way the non-permissions required variables can be loaded here
			// And we can return a generic 500 here by returning early
			return c.Next()
		}
		if len(ps) == 0 {
			perms := player.NewPermissions(pid.(int64), []query.PlayerPermission{})
			c.Locals("perms", perms)
			return c.Next()
		}

		perms := player.NewPermissions(pid.(int64), ps)
		c.Locals("perms", perms)
		return c.Next()
	}
}
