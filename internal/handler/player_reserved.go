package handler

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
)

func UsernameReserved(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Username string `form:"username"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			return err
		}

		p, err := i.Queries.GetPlayerByUsername(context.Background(), in.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
				return c.Render(partial.PlayerFree, fiber.Map{
					"CSRF": c.Locals("csrf"),
				}, layout.CSRF)
			}
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.PlayerReservedErr, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		}

		if in.Username == p.Username {
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partial.PlayerReserved, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		} else {
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			return c.Render(partial.PlayerFree, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		}
	}
}
