package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/shared"
)

func CharacterApplicationsQueuePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		perms, ok := lperms.(permissions.PlayerGranted)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
		}
		_, ok = perms.Permissions[permissions.PlayerReviewCharacterApplicationsName]
		if !ok {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		apps, err := qtx.ListOpenCharacterApplications(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(bind.Name))
		}

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
					return c.Render("views/500", c.Locals(bind.Name))
				}
				reviewer = p.Username
			}
			summaries = append(summaries, request.NewSummaryFromApplication(&app.Player, reviewer, &app.Request, &app.CharacterApplicationContent))
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := c.Locals(bind.Name).(fiber.Map)
		b["ThereAreCharacterApplications"] = len(summaries) > 0
		b["CharacterApplicationSummaries"] = summaries
		return c.Render("views/character/application/queue", b)
	}
}
