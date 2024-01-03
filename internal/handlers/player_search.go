package handlers

import (
	"context"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/views"
)

func SearchPlayer(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Search string `form:"search"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, c.Locals(constants.BindName), "views/layouts/standalone")
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		searchStr := fmt.Sprintf("%%%s%%", r.Search)
		players, err := i.Queries.SearchPlayersByUsername(context.Background(), searchStr)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		dest := c.Params("dest")
		if len(dest) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if dest == "player-permissions" {
			// TODO: Move this to a constant and inject it
			b := c.Locals(constants.BindName).(fiber.Map)
			b["Players"] = players

			c.Status(fiber.StatusOK)
			return c.Render("views/partials/player_permissions/search_results", b, "")
		}

		c.Status(fiber.StatusBadRequest)
		return nil
	}
}
