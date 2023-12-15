package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func CharactersPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		apps, err := i.Queries.ListCharacterApplicationsForPlayer(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind))
		}

		summaries := []character.ApplicationSummary{}
		for _, app := range apps {
			reviewer := ""
			if app.Request.RPID > 0 {
				p, err := i.Queries.GetPlayer(context.Background(), app.Request.RPID)
				if err != nil {
					// TODO: Sort out this edge case
					// if err == sql.ErrNoRows {
					// TODO: Log this error here, this means we need to reset the reviewer and status on the request
					// }
					c.Status(fiber.StatusInternalServerError)
					return c.Render("views/500", c.Locals(shared.Bind))
				}
				reviewer = p.Username
			}
			summaries = append(summaries, character.NewSummaryFromApplication(&app.Player, reviewer, &app.Request, &app.CharacterApplicationContent))
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		b["NewCharacterApplicationPath"] = routes.NewCharacterApplicationPath()
		b["CharacterApplicationSummaries"] = summaries
		b["HasCharacterApplications"] = len(apps) > 0
		return c.Render("views/characters", b)
	}
}
