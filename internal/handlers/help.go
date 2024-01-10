package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
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

		headers, err := qtx.ListHelpHeaders(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		help := []fiber.Map{}
		for _, header := range headers {
			tags, err := qtx.GetTagsForHelpFile(context.Background(), header.Slug)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
			}

			help = append(help, fiber.Map{
				"Title":    header.Title,
				"Sub":      header.Sub,
				"Category": header.Category,
				"Tags":     tags,
				"Path":     routes.HelpFilePath(header.Slug),
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
				"Path":  routes.HelpFilePath(record.RelatedSlug),
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
		Search   string `form:"search"`
		Title    bool   `form:"title"`
		Content  bool   `form:"content"`
		Category bool   `form:"category"`
		Tags     bool   `form:"tags"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusBadRequest)
			b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "search-help-error",
				SectionClass: "py-4 px-6 w-[60%]",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			})
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}

		if len(in.Search) == 0 {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusBadRequest)
			b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "search-help-error",
				SectionClass: "py-4 px-6 w-[60%]",
				NoticeText: []string{
					"Please enter a search term to search.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			})
			return c.Render(partials.NoticeSectionWarn, b, layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusInternalServerError)
			b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "search-help-error",
				SectionClass: "py-4 px-6 w-[60%]",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			})
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		var sb strings.Builder
		fmt.Fprintf(&sb, "%%%s%%", in.Search)
		search := sb.String()

		results := []fiber.Map{}

		if in.Title {
			byTitle, err := qtx.SearchHelpByTitle(context.Background(), queries.SearchHelpByTitleParams{
				Slug:  search,
				Title: search,
			})
			if err != nil {
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Append("HX-Retarget", "search-help-error")
				c.Status(fiber.StatusInternalServerError)
				b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    "search-help-error",
					SectionClass: "py-4 px-6 w-[60%]",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				})
				return c.Render(partials.NoticeSectionError, b, layouts.None)
			}

			foundByTitle := []fiber.Map{}
			for _, header := range byTitle {
				tags, err := qtx.GetTagsForHelpFile(context.Background(), header.Slug)
				if err != nil {
					c.Append(shared.HeaderHXAcceptable, "true")
					c.Append("HX-Retarget", "search-help-error")
					c.Status(fiber.StatusInternalServerError)
					b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
						SectionID:    "search-help-error",
						SectionClass: "py-4 px-6 w-[60%]",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					})
					return c.Render(partials.NoticeSectionError, b, layouts.None)
				}

				foundByTitle = append(foundByTitle, fiber.Map{
					"Title":    header.Title,
					"Sub":      header.Sub,
					"Category": header.Category,
					"Tags":     tags,
					"Path":     routes.HelpFilePath(header.Slug),
				})
			}

			var containsHeaderSB strings.Builder
			fmt.Fprintf(&containsHeaderSB, "Title Contains \"%s\"", in.Search)

			results = append(results, fiber.Map{
				"ResultSets": []fiber.Map{
					{
						"Header":  containsHeaderSB.String(),
						"Results": foundByTitle,
					},
				},
			})
		}

		if in.Content {
			byContent, err := qtx.SearchHelpByContent(context.Background(), queries.SearchHelpByContentParams{
				Sub: search,
				Raw: search,
			})
			if err != nil {
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Append("HX-Retarget", "search-help-error")
				c.Status(fiber.StatusInternalServerError)
				b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    "search-help-error",
					SectionClass: "py-4 px-6 w-[60%]",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				})
				return c.Render(partials.NoticeSectionError, b, layouts.None)
			}

			foundByContent := []fiber.Map{}
			for _, header := range byContent {
				tags, err := qtx.GetTagsForHelpFile(context.Background(), header.Slug)
				if err != nil {
					c.Append(shared.HeaderHXAcceptable, "true")
					c.Append("HX-Retarget", "search-help-error")
					c.Status(fiber.StatusInternalServerError)
					b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
						SectionID:    "search-help-error",
						SectionClass: "py-4 px-6 w-[60%]",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					})
					return c.Render(partials.NoticeSectionError, b, layouts.None)
				}

				foundByContent = append(foundByContent, fiber.Map{
					"Title":    header.Title,
					"Sub":      header.Sub,
					"Category": header.Category,
					"Tags":     tags,
					"Path":     routes.HelpFilePath(header.Slug),
				})
			}

			var containsHeaderSB strings.Builder
			fmt.Fprintf(&containsHeaderSB, "Help Files Containing \"%s\"", in.Search)

			results = append(results, fiber.Map{
				"ResultSets": []fiber.Map{
					{
						"Header":  containsHeaderSB.String(),
						"Results": foundByContent,
					},
				},
			})
		}

		if in.Category {
			byCategory, err := qtx.SearchHelpByCategory(context.Background(), search)
			if err != nil {
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Append("HX-Retarget", "search-help-error")
				c.Status(fiber.StatusInternalServerError)
				b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    "search-help-error",
					SectionClass: "py-4 px-6 w-[60%]",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				})
				return c.Render(partials.NoticeSectionError, b, layouts.None)
			}

			foundByCategory := []fiber.Map{}
			for _, header := range byCategory {
				tags, err := qtx.GetTagsForHelpFile(context.Background(), header.Slug)
				if err != nil {
					c.Append(shared.HeaderHXAcceptable, "true")
					c.Append("HX-Retarget", "search-help-error")
					c.Status(fiber.StatusInternalServerError)
					b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
						SectionID:    "search-help-error",
						SectionClass: "py-4 px-6 w-[60%]",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					})
					return c.Render(partials.NoticeSectionError, b, layouts.None)
				}

				foundByCategory = append(foundByCategory, fiber.Map{
					"Title":    header.Title,
					"Sub":      header.Sub,
					"Category": header.Category,
					"Tags":     tags,
					"Path":     routes.HelpFilePath(header.Slug),
				})
			}

			var sb strings.Builder
			fmt.Fprintf(&sb, "Help Files With Category Containing \"%s\"", in.Search)

			results = append(results, fiber.Map{
				"ResultSets": []fiber.Map{
					{
						"Header":  sb.String(),
						"Results": foundByCategory,
					},
				},
			})
		}

		if in.Tags {
			byTags, err := qtx.SearchHelpByTags(context.Background(), search)
			if err != nil {
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Append("HX-Retarget", "search-help-error")
				c.Status(fiber.StatusInternalServerError)
				b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    "search-help-error",
					SectionClass: "py-4 px-6 w-[60%]",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				})
				return c.Render(partials.NoticeSectionError, b, layouts.None)
			}

			foundByTags := []fiber.Map{}
			for _, header := range byTags {
				tags, err := qtx.GetTagsForHelpFile(context.Background(), header.Slug)
				if err != nil {
					c.Append(shared.HeaderHXAcceptable, "true")
					c.Append("HX-Retarget", "search-help-error")
					c.Status(fiber.StatusInternalServerError)
					b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
						SectionID:    "search-help-error",
						SectionClass: "py-4 px-6 w-[60%]",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					})
					return c.Render(partials.NoticeSectionError, b, layouts.None)
				}

				foundByTags = append(foundByTags, fiber.Map{
					"Title":    header.Title,
					"Sub":      header.Sub,
					"Category": header.Category,
					"Tags":     tags,
					"Path":     routes.HelpFilePath(header.Slug),
				})
			}

			var sb strings.Builder
			fmt.Fprintf(&sb, "Help Files With Tags Containing \"%s\"", in.Search)

			results = append(results, fiber.Map{
				"ResultSets": []fiber.Map{
					{
						"Header":  sb.String(),
						"Results": foundByTags,
					},
				},
			})

		}

		if err := tx.Commit(); err != nil {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusInternalServerError)
			b := partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "search-help-error",
				SectionClass: "py-4 px-6 w-[60%]",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			})
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}

		if len(results) == 0 {
			b := views.Bind(c)
			b["Search"] = in.Search
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusNotFound)
			return c.Render(partials.HelpIndexSearchNoResults, b, layouts.None)
		}

		b := views.Bind(c)
		b["Results"] = results
		return c.Render(partials.HelpIndexSearchResults, b, layouts.None)
	}
}
