package handlers

import (
	"context"
	"slices"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
)

type LoginInput struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Login(s *session.Store, q *queries.Queries, r *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		i := new(LoginInput)
		if err := c.BodyParser(i); err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		p, err := q.GetPlayerByUsername(context.Background(), i.Username)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		v, err := password.Verify(i.Password, p.PwHash)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}
		if !v {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		pid := p.ID
		perms, err := permissions.List(q, r, pid)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}
		if !slices.Contains(perms, permissions.Login) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		sess, err := s.Get(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		sess.Set("pid", pid)
		if err = sess.Save(); err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}
