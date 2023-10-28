package sessiondata

import (
	"context"
	"fmt"
	"log"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	redis "github.com/redis/go-redis/v9"

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

			permkey := fmt.Sprintf("perm:%v", pid)
			permKeyExists, err := r.Exists(context.Background(), permkey).Result()
			if err != nil {
				log.Print(err)
				return c.Next()
			}

			if permKeyExists == 1 {
				rperms, err := r.SMembers(context.Background(), permkey).Result()
				if err != nil {
					log.Print(err)
					return c.Next()
				}
				c.Locals("perms", rperms)
				return c.Next()
			} else {
				qpermrecords, err := q.ListPlayerPermissions(context.Background(), pid.(int64))
				if err != nil {
					log.Print(err)
					return c.Next()
				}

				var qperms []string
				for i := 0; i < len(qpermrecords); i++ {
					qpermrecord := qpermrecords[i]
					qperms = append(qperms, qpermrecord.Permission)
				}

				c.Locals("perms", qperms)
				r.SAdd(context.Background(), permkey, strings.Join(qperms, " "), TwoHoursInSeconds)
				return c.Next()
			}
		}

		return c.Next()
	}
}
