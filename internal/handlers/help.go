package handlers

import (
	"context"
	"database/sql"
	"html/template"
	"slices"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/views"
)

func HelpPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		slugs, err := qtx.ListHelpSlugs(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		slug := c.Params("slug")
		if !slices.Contains(slugs, slug) {
			c.Status(fiber.StatusNotFound)
			return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
		}

		help, err := qtx.GetHelp(context.Background(), slug)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		html := template.HTML(help.Html)
		b := views.Bind(c)
		b["Content"] = html
		return c.Render(views.Help, b, layouts.Main)
	}
}
