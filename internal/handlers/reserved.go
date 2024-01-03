package handlers

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/views"
)

func Reserved(i *shared.Interfaces) fiber.Handler {
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
				return c.Render("views/partials/register/player-free", fiber.Map{
					"CSRF": c.Locals("csrf"),
				}, views.LayoutCSRF)
			}
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/player-reserved-err", fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, views.LayoutCSRF)
		}

		if r.Username == p.Username {
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render("views/partials/register/player-reserved", fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, views.LayoutCSRF)
		} else {
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			return c.Render("views/partials/register/player-free", fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, views.LayoutCSRF)
		}
	}
}
