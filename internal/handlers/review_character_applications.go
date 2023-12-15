package handlers

import (
	"context"
	"log"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/shared"
)

func ReviewCharacterApplicationsPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		perms, ok := lperms.(permission.PlayerGranted)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		_, ok = perms.Permissions[permission.PlayerReviewCharacterApplicationsName]
		if !ok {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		apps, err := i.Queries.ListOpenCharacterApplications(context.Background())
		if err != nil {
			log.Println(err)
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
		b["ThereAreCharacterApplications"] = len(summaries) > 0
		b["CharacterApplicationSummaries"] = summaries
		return c.Render("views/character/applications/review", b)
	}
}
