package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/middleware"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/username"
)

type Player struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func NewPlayer(c *fiber.Ctx) error {
	p := new(Player)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	u := username.Sanitize(p.Username)

	if !username.Validate(u) {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	if !password.Validate(p.Password) {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	pw_hash, err := password.Hash(p.Password)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	ctx := context.Background()
	result, err := queries.Q.CreatePlayer(ctx, queries.CreatePlayerParams{
		Username: u,
		PwHash:   pw_hash,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	pid, err := result.LastInsertId()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	sess, err := middleware.Sessions.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	sess.Set("pid", pid)
	if err = sess.Save(); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusCreated)
	return nil
}
