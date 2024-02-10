package handler

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
)

func UsernameReserved(i *interfaces.Shared) fiber.Handler {
	type request struct {
		Username string `form:"username"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			return err
		}

		p, err := i.Queries.GetPlayerByUsername(context.Background(), r.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
				return c.Render(partials.PlayerFree, fiber.Map{
					"CSRF": c.Locals("csrf"),
				}, layouts.CSRF)
			}
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.PlayerReservedErr, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layouts.CSRF)
		}

		if r.Username == p.Username {
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partials.PlayerReserved, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layouts.CSRF)
		} else {
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			return c.Render(partials.PlayerFree, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layouts.CSRF)
		}
	}
}
