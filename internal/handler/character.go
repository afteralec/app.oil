package handler

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

func CharactersPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}
		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		// TODO: Make this a ListRequestsForPlayerByType query instead
		apps, err := qtx.ListCharacterApplicationsForPlayer(context.Background(), pid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c))
		}

		// TODO: Get this into a standard API on the request package
		summaries := []request.SummaryForQueue{}
		for _, app := range apps {
			content, err := request.Content(qtx, &app.Request)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summary, err := request.NewSummaryForQueue(request.SummaryForQueueParams{
				Request: &app.Request,
				Content: content,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summaries = append(summaries, summary)
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := view.Bind(c)
		b["RequestsPath"] = route.Requests
		b["CharacterApplicationSummaries"] = summaries
		b["HasCharacterApplications"] = len(apps) > 0
		return c.Render(view.Characters, b)
	}
}

func CharacterApplicationsQueuePage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		if !perms.HasPermission(player.PermissionReviewCharacterApplications.Name) {
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

		// TODO: Make this a "List Open Requests By Type" query
		apps, err := qtx.ListOpenCharacterApplications(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c))
		}

		summaries := []request.SummaryForQueue{}
		for _, app := range apps {
			content, err := request.Content(qtx, &app.Request)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summary, err := request.NewSummaryForQueue(request.SummaryForQueueParams{
				Request: &app.Request,
				Content: content,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summaries = append(summaries, summary)
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := view.Bind(c)
		// TODO: Move this length check down into the template
		b["ThereAreCharacterApplications"] = len(summaries) > 0
		b["CharacterApplicationSummaries"] = summaries
		return c.Render(view.CharacterApplicationQueue, b)
	}
}
