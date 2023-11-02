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
			c.Append("HX-Trigger-After-Swap", "username-reserved")
			return c.Render("web/views/htmx/player-free", fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		}

		if r.Username == u {
			c.Append("HX-Trigger-After-Swap", "username-reserved")
			return c.Render("web/views/htmx/player-reserved", fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		} else {
			c.Append("HX-Trigger-After-Swap", "username-reserved")
			return c.Render("web/views/htmx/player-free", fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		}
	}
}
