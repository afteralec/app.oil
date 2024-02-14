package handler

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"slices"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/view"
)

func HelpPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		header, err := qtx.ListHelpHeaders(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		help := []fiber.Map{}
		for _, header := range header {
			tags, err := qtx.GetTagsForHelpFile(context.Background(), header.Slug)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}

			help = append(help, fiber.Map{
				"Title":    header.Title,
				"Sub":      header.Sub,
				"Category": header.Category,
				"Tags":     tags,
				"Path":     route.HelpFilePath(header.Slug),
			})
		}

		b := view.Bind(c)
		b["Help"] = help
		return c.Render(view.Help, b, layout.Main)
	}
}

func HelpFilePage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		slugs, err := qtx.ListHelpSlugs(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		slug := c.Params("slug")
		if !slices.Contains(slugs, slug) {
			c.Status(fiber.StatusNotFound)
			return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
		}

		help, err := qtx.GetHelp(context.Background(), slug)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		relatedRecords, err := qtx.GetHelpRelated(context.Background(), slug)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		tags, err := qtx.GetTagsForHelpFile(context.Background(), slug)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		related := []fiber.Map{}
		for _, record := range relatedRecords {
			related = append(related, fiber.Map{
				"Title": record.RelatedTitle,
				"Sub":   record.RelatedSub,
				"Path":  route.HelpFilePath(record.RelatedSlug),
			})
		}

		html := template.HTML(help.HTML)
		b := view.Bind(c)
		b["Content"] = html
		b["Related"] = related
		// TODO: Once the help path can take a query string, save the last state of the session's help path
		b["HelpPath"] = route.Help
		b["HelpTitle"] = help.Title
		b["Sub"] = help.Sub
		b["Category"] = help.Category
		b["Tags"] = tags
		return c.Render(view.HelpFile, b, layout.Main)
	}
}

func SearchHelp(i *interfaces.Shared) fiber.Handler {
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
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusBadRequest)
			b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    "search-help-error",
				SectionClass: "py-4 px-6 w-[60%]",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			})
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}

		if len(in.Search) == 0 {
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusBadRequest)
			b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    "search-help-error",
				SectionClass: "py-4 px-6 w-[60%]",
				NoticeText: []string{
					"Please enter a search term to search.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			})
			return c.Render(partial.NoticeSectionWarn, b, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusInternalServerError)
			b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    "search-help-error",
				SectionClass: "py-4 px-6 w-[60%]",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			})
			return c.Render(partial.NoticeSectionError, b, layout.None)
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
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", "search-help-error")
				c.Status(fiber.StatusInternalServerError)
				b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    "search-help-error",
					SectionClass: "py-4 px-6 w-[60%]",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				})
				return c.Render(partial.NoticeSectionError, b, layout.None)
			}

			foundByTitle := []fiber.Map{}
			for _, hdr := range byTitle {
				tags, err := qtx.GetTagsForHelpFile(context.Background(), hdr.Slug)
				if err != nil {
					c.Append(header.HXAcceptable, "true")
					c.Append("HX-Retarget", "search-help-error")
					c.Status(fiber.StatusInternalServerError)
					b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
						SectionID:    "search-help-error",
						SectionClass: "py-4 px-6 w-[60%]",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					})
					return c.Render(partial.NoticeSectionError, b, layout.None)
				}

				foundByTitle = append(foundByTitle, fiber.Map{
					"Title":    hdr.Title,
					"Sub":      hdr.Sub,
					"Category": hdr.Category,
					"Tags":     tags,
					"Path":     route.HelpFilePath(hdr.Slug),
				})
			}

			var containsheaderB strings.Builder
			fmt.Fprintf(&containsheaderB, "Title Contains \"%s\"", in.Search)

			results = append(results, fiber.Map{
				"ResultSets": []fiber.Map{
					{
						"Header":  containsheaderB.String(),
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
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", "search-help-error")
				c.Status(fiber.StatusInternalServerError)
				b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    "search-help-error",
					SectionClass: "py-4 px-6 w-[60%]",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				})
				return c.Render(partial.NoticeSectionError, b, layout.None)
			}

			foundByContent := []fiber.Map{}
			for _, hdr := range byContent {
				tags, err := qtx.GetTagsForHelpFile(context.Background(), hdr.Slug)
				if err != nil {
					c.Append(header.HXAcceptable, "true")
					c.Append("HX-Retarget", "search-help-error")
					c.Status(fiber.StatusInternalServerError)
					b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
						SectionID:    "search-help-error",
						SectionClass: "py-4 px-6 w-[60%]",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					})
					return c.Render(partial.NoticeSectionError, b, layout.None)
				}

				foundByContent = append(foundByContent, fiber.Map{
					"Title":    hdr.Title,
					"Sub":      hdr.Sub,
					"Category": hdr.Category,
					"Tags":     tags,
					"Path":     route.HelpFilePath(hdr.Slug),
				})
			}

			var containsheaderB strings.Builder
			fmt.Fprintf(&containsheaderB, "Help Files Containing \"%s\"", in.Search)

			results = append(results, fiber.Map{
				"ResultSets": []fiber.Map{
					{
						"Header":  containsheaderB.String(),
						"Results": foundByContent,
					},
				},
			})
		}

		if in.Category {
			byCategory, err := qtx.SearchHelpByCategory(context.Background(), search)
			if err != nil {
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", "search-help-error")
				c.Status(fiber.StatusInternalServerError)
				b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    "search-help-error",
					SectionClass: "py-4 px-6 w-[60%]",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				})
				return c.Render(partial.NoticeSectionError, b, layout.None)
			}

			foundByCategory := []fiber.Map{}
			for _, hdr := range byCategory {
				tags, err := qtx.GetTagsForHelpFile(context.Background(), hdr.Slug)
				if err != nil {
					c.Append(header.HXAcceptable, "true")
					c.Append("HX-Retarget", "search-help-error")
					c.Status(fiber.StatusInternalServerError)
					b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
						SectionID:    "search-help-error",
						SectionClass: "py-4 px-6 w-[60%]",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					})
					return c.Render(partial.NoticeSectionError, b, layout.None)
				}

				foundByCategory = append(foundByCategory, fiber.Map{
					"Title":    hdr.Title,
					"Sub":      hdr.Sub,
					"Category": hdr.Category,
					"Tags":     tags,
					"Path":     route.HelpFilePath(hdr.Slug),
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
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", "search-help-error")
				c.Status(fiber.StatusInternalServerError)
				b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    "search-help-error",
					SectionClass: "py-4 px-6 w-[60%]",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				})
				return c.Render(partial.NoticeSectionError, b, layout.None)
			}

			foundByTags := []fiber.Map{}
			for _, hdr := range byTags {
				tags, err := qtx.GetTagsForHelpFile(context.Background(), hdr.Slug)
				if err != nil {
					c.Append(header.HXAcceptable, "true")
					c.Append("HX-Retarget", "search-help-error")
					c.Status(fiber.StatusInternalServerError)
					b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
						SectionID:    "search-help-error",
						SectionClass: "py-4 px-6 w-[60%]",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					})
					return c.Render(partial.NoticeSectionError, b, layout.None)
				}

				foundByTags = append(foundByTags, fiber.Map{
					"Title":    hdr.Title,
					"Sub":      hdr.Sub,
					"Category": hdr.Category,
					"Tags":     tags,
					"Path":     route.HelpFilePath(hdr.Slug),
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
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusInternalServerError)
			b := partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    "search-help-error",
				SectionClass: "py-4 px-6 w-[60%]",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			})
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}

		if len(results) == 0 {
			b := view.Bind(c)
			b["Search"] = in.Search
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", "search-help-error")
			c.Status(fiber.StatusNotFound)
			return c.Render(partial.HelpIndexSearchNoResults, b, layout.None)
		}

		b := view.Bind(c)
		b["Results"] = results
		return c.Render(partial.HelpIndexSearchResults, b, layout.None)
	}
}
