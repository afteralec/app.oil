package handler

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/views"
)

func CharactersPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layout.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		apps, err := qtx.ListCharacterApplicationsForPlayer(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c))
		}

		// TODO: Get this into a standard API on the request package
		summaries := []request.ApplicationSummary{}
		for _, app := range apps {
			reviewer := ""
			if app.Request.RPID > 0 {
				p, err := qtx.GetPlayer(context.Background(), app.Request.RPID)
				if err != nil {
					// TODO: Sort out this edge case
					// if err == sql.ErrNoRows {
					// TODO: Log this error here, this means we need to reset the reviewer and status on the request
					// }
					c.Status(fiber.StatusInternalServerError)
					return c.Render(views.InternalServerError, views.Bind(c))
				}
				reviewer = p.Username
			}
			summaries = append(summaries, request.NewSummaryFromApplication(&app.Player, reviewer, &app.Request, &app.CharacterApplicationContent))
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := views.Bind(c)
		b["RequestsPath"] = routes.Requests
		b["CharacterApplicationSummaries"] = summaries
		b["HasCharacterApplications"] = len(apps) > 0
		return c.Render(views.Characters, b)
	}
}
