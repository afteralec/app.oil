package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/middleware"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
)

type LoginInput struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Login(c *fiber.Ctx) error {
	i := new(LoginInput)

	if err := c.BodyParser(i); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	ctx := context.Background()
	p, err := queries.Q.GetPlayerByUsername(ctx, i.Username)
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

	sess, err := middleware.Sessions.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
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
