package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"slices"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
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

		headers, err := qtx.ListHelpTitleAndSub(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		help := []fiber.Map{}
		for _, header := range headers {
			help = append(help, fiber.Map{
				"Title": header.Title,
				"Sub":   header.Sub,
				"Path":  routes.HelpFilePath(header.Slug),
			})
		}

		b := views.Bind(c)
		b["Help"] = help
		return c.Render(views.Help, b, layouts.Main)
	}
}

func HelpFilePage(i *shared.Interfaces) fiber.Handler {
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

		relatedRecords, err := qtx.GetHelpRelated(context.Background(), slug)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		related := []fiber.Map{}
		for _, record := range relatedRecords {
			related = append(related, fiber.Map{
				"Title": record.RelatedTitle,
				"Sub":   record.RelatedSub,
				"Path":  routes.HelpFilePath(record.Related),
			})
		}

		html := template.HTML(help.HTML)
		b := views.Bind(c)
		b["Content"] = html
		b["Related"] = related
		// TODO: Once the help path can take a query string, save the last state of the session's help path
		b["HelpPath"] = routes.Help
		return c.Render(views.HelpFile, b, layouts.Main)
	}
}

func SearchHelp(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Search string `form:"search"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		var sb strings.Builder
		fmt.Fprintf(&sb, "%%%s%%", in.Search)
		search := sb.String()

		byTitle, err := qtx.SearchHelpByTitle(context.Background(), queries.SearchHelpByTitleParams{
			Slug:  search,
			Title: search,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		log.Println(len(byTitle))

		byContent, err := qtx.SearchHelpByContent(context.Background(), queries.SearchHelpByContentParams{
			Sub: search,
			Raw: search,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if len(byTitle) == 0 && len(byContent) == 0 {
			c.Status(fiber.StatusNotFound)
			return nil
		}

		foundByTitle := []fiber.Map{}
		for _, header := range byTitle {
			foundByTitle = append(foundByTitle, fiber.Map{
				"Title": header.Title,
				"Sub":   header.Sub,
				"Path":  routes.HelpFilePath(header.Slug),
			})
		}

		foundByContent := []fiber.Map{}
		for _, header := range byContent {
			foundByContent = append(foundByContent, fiber.Map{
				"Title": header.Title,
				"Sub":   header.Sub,
				"Path":  routes.HelpFilePath(header.Slug),
			})
		}

		var contentHeaderSB strings.Builder
		fmt.Fprintf(&contentHeaderSB, "Help Files Containing \"%s\"", in.Search)
		var titleHeaderSB strings.Builder
		fmt.Fprintf(&titleHeaderSB, "Best Matches for \"%s\"", in.Search)

		b := views.Bind(c)
		b["ByTitle"] = foundByTitle
		b["ByContent"] = foundByContent
		b["TitleHeader"] = titleHeaderSB.String()
		b["ContentHeader"] = contentHeaderSB.String()
		return c.Render(partials.HelpIndexSearchResults, b, layouts.None)
	}
}
