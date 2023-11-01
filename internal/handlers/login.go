package handlers

import (
	"context"
	"slices"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/shared"
)

func Login(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		p, err := i.Queries.GetPlayerByUsername(context.Background(), r.Username)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		v, err := password.Verify(r.Password, p.PwHash)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}
		if !v {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		pid := p.ID
		perms, err := permissions.List(i, pid)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}
		if !slices.Contains(perms, permissions.Login) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		sess, err := i.Sessions.Get(c)
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
