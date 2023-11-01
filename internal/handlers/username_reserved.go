package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/shared"
)

func UsernameReserved(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Username string `form:"username"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			return err
		}

		ctx := context.Background()
		u, err := i.Queries.GetPlayerUsername(ctx, r.Username)
		// TODO: Figure out getting the error code from these
		if err != nil {
			return c.Render("web/views/htmx/player-free", fiber.Map{}, "")
		}

		if r.Username == u {
			return c.Render("web/views/htmx/player-reserved", fiber.Map{}, "")
		} else {
			return c.Render("web/views/htmx/player-free", fiber.Map{}, "")
		}
	}
}
