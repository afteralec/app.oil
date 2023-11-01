package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
)

type Input struct {
	Username string `form:"username"`
}

func UsernameReserved(q *queries.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		i := new(Input)
		if err := c.BodyParser(i); err != nil {
			return err
		}

		ctx := context.Background()
		u, err := q.GetPlayerUsername(ctx, i.Username)
		// TODO: Figure out getting the error code from these
		if err != nil {
			return c.Render("web/views/htmx/player-free", fiber.Map{}, "")
		}

		if i.Username == u {
			return c.Render("web/views/htmx/player-reserved", fiber.Map{}, "")
		} else {
			return c.Render("web/views/htmx/player-free", fiber.Map{}, "")
		}
	}
}
