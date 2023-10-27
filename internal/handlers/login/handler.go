package login

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
)

type LoginInput struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func New(s *session.Store, q *queries.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		i := new(LoginInput)

		if err := c.BodyParser(i); err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		ctx := context.Background()
		p, err := q.GetPlayerByUsername(ctx, i.Username)
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

		sess, err := s.Get(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		sess.Set("pid", p.ID)
		if err = sess.Save(); err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}
