package sessiondata

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
)

const TwoHoursInSeconds = 120 * 60

func New(s *session.Store, q *queries.Queries, r *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := s.Get(c)
		if err != nil {
			log.Print(err)
			return c.Next()
		}

		pid := sess.Get("pid")
		if pid != nil {
			c.Locals("pid", pid)

			perms, err := permissions.List(q, r, pid.(int64))
			if err != nil {
				log.Print(err)
				return c.Next()
			}
			c.Locals("perms", perms)
		}

		return c.Next()
	}
}
