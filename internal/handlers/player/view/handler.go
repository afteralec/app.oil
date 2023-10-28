package viewplayer

import (
	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/queries"
)

func New(q *queries.Queries, r *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			return c.Redirect("/")
		}

		id := c.Params("id")
		b := c.Locals("bind").(fiber.Map)
		b["ID"] = id

		return c.Render("web/views/player", b)
	}
}
