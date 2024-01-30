package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/actors"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/views"
)

func ActorImagesPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layouts.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		if !perms.HasPermission(permissions.PlayerViewAllActorImagesName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		actorImages, err := i.Queries.ListActorImages(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		pageActorImages := []fiber.Map{}
		for _, actorImage := range actorImages {
			var sb strings.Builder
			fmt.Fprintf(&sb, "[%d] %s", actorImage.ID, actorImage.ShortDescription)
			pageActorImage := fiber.Map{
				"Title": sb.String(),
				"Name":  actorImage.Name,
				"Path":  routes.ActorImagePath(actorImage.ID),
			}

			if perms.HasPermission(permissions.PlayerCreateActorImageName) {
				pageActorImage["EditPath"] = routes.EditActorImagePath(actorImage.ID)
			}

			pageActorImages = append(pageActorImages, pageActorImage)
		}

		b := views.Bind(c)
		if perms.HasPermission(permissions.PlayerCreateActorImageName) {
			b["CreatePermission"] = true
		}
		b["ActorImages"] = pageActorImages
		b["PageHeader"] = fiber.Map{
			"Title":    "Actor Images",
			"SubTitle": "Actor images are where the primary properties for an actor are defined, like a template",
		}
		return c.Render(views.ActorImages, b)
	}
}

func EditActorImagePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layouts.Standalone)
		}

		aiid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		if !perms.HasPermission(permissions.PlayerCreateActorImageName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		actorImage, err := qtx.GetActorImage(context.Background(), aiid)
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

		b := views.Bind(c)
		// TODO: Get a bind function for this
		b["NavBack"] = fiber.Map{
			"Path":  routes.ActorImages,
			"Label": "Back to Actor Images",
		}
		// TODO: Get a bind function for this too
		b["PageHeader"] = fiber.Map{
			"Title":    actors.ImageTitleWithID(actorImage.Name, actorImage.ID),
			"SubTitle": "Update actor properties here",
		}
		// TODO: Write a bind function for this
		b["Name"] = actorImage.Name
		b["ShortDescription"] = actorImage.ShortDescription
		b["Description"] = actorImage.Description
		b["ShortDescriptionPath"] = routes.ActorImageShortDescriptionPath(aiid)
		b["DescriptionPath"] = routes.ActorImageDescriptionPath(aiid)
		return c.Render(views.EditActorImage, b)
	}
}

func ActorImagePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layouts.Standalone)
		}

		aiid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		if !perms.HasPermission(permissions.PlayerViewAllActorImagesName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		actorImage, err := qtx.GetActorImage(context.Background(), aiid)
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

		b := views.Bind(c)
		// TODO: Get a bind function for this
		b["NavBack"] = fiber.Map{
			"Path":  routes.ActorImages,
			"Label": "Back to Actor Images",
		}
		// TODO: Get a bind function for this too
		b["PageHeader"] = fiber.Map{
			"Title": actors.ImageTitleWithID(actorImage.Name, actorImage.ID),
		}
		b["Name"] = actorImage.Name
		b["ShortDescription"] = actorImage.ShortDescription
		b["Description"] = actorImage.Description
		return c.Render(views.ActorImage, b)
	}
}

func NewActorImage(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Name string `form:"name"`
	}

	const sectionID string = "actor-image-create-error"

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if !actors.IsImageNameValid(in.Name) {
			c.Status(fiber.StatusBadRequest)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"The Image Name you entered isn't valid. Please try again.",
				},
				NoticeIcon: true,
			}), layouts.None)
		}

		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"It looks like your session may have expired.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if !perms.HasPermission(permissions.PlayerCreateActorImageName) {
			c.Status(fiber.StatusForbidden)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"You don't have the permission(s) necessary to create an Actor Image.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		result, err := qtx.CreateActorImage(context.Background(), queries.CreateActorImageParams{
			Name:             in.Name,
			ShortDescription: actors.DefaultImageShortDescription,
			Description:      actors.DefaultImageDescription,
			Gender:           actors.DefaultImageGender,
		})
		if err != nil {
			if me, ok := err.(*mysql.MySQLError); ok {
				if me.Number == mysqlerr.ER_DUP_ENTRY {
					c.Status(fiber.StatusConflict)
					c.Append(shared.HeaderHXAcceptable, "true")
					c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
					return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
						SectionID:    sectionID,
						SectionClass: "pt-2",
						NoticeText: []string{
							"That Actor Image name is already in use. Please choose another.",
						},
						NoticeIcon: true,
					}), layouts.None)
				}
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		aiid, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		c.Status(fiber.StatusCreated)
		c.Append("HX-Redirect", routes.EditActorImagePath(aiid))
		c.Append("HX-Reswap", "none")
		return nil
	}
}

func EditActorImageShortDescription(i *shared.Interfaces) fiber.Handler {
	type input struct {
		ShortDescription string `form:"sdesc"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !actors.IsShortDescriptionValid(in.ShortDescription) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !perms.HasPermission(permissions.PlayerCreateActorImageName) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		aiid, err := util.GetID(c)
		if err != nil {
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

		actorImage, err := qtx.GetActorImage(context.Background(), aiid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if actorImage.ShortDescription == in.ShortDescription {
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err := qtx.UpdateActorImageShortDescription(context.Background(), queries.UpdateActorImageShortDescriptionParams{
			ID:               actorImage.ID,
			ShortDescription: in.ShortDescription,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		actorImage, err = qtx.GetActorImage(context.Background(), aiid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := fiber.Map{}
		b["ShortDescription"] = actorImage.ShortDescription
		b["ShortDescriptionPath"] = routes.ActorImageShortDescriptionPath(actorImage.ID)
		b["NoticeSection"] = partials.BindNoticeSection(partials.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "actor-image-edit-short-description-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The short description has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partials.ActorImageEditShortDescription, b, layouts.None)
	}
}

func EditActorImageDescription(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Description string `form:"desc"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !actors.IsDescriptionValid(in.Description) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !perms.HasPermission(permissions.PlayerCreateActorImageName) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		aiid, err := util.GetID(c)
		if err != nil {
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

		actorImage, err := qtx.GetActorImage(context.Background(), aiid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if actorImage.Description == in.Description {
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err := qtx.UpdateActorImageDescription(context.Background(), queries.UpdateActorImageDescriptionParams{
			ID:          actorImage.ID,
			Description: in.Description,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		actorImage, err = qtx.GetActorImage(context.Background(), actorImage.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := fiber.Map{}
		b["Description"] = actorImage.Description
		b["DescriptionPath"] = routes.ActorImageDescriptionPath(actorImage.ID)
		b["NoticeSection"] = partials.BindNoticeSection(partials.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "actor-image-edit-description-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The description has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partials.ActorImageEditDescription, b, layouts.None)
	}
}

func ActorImageNameReserved(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Name string `form:"name"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		actorImage, err := i.Queries.GetActorImageByName(context.Background(), in.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Trigger-After-Swap", "ptrcr:actor-image-reserved")
				return c.Render(partials.ActorImageFree, fiber.Map{
					"CSRF": c.Locals("csrf"),
				}, layouts.CSRF)
			}
			c.Append("HX-Trigger-After-Swap", "ptrcr:actor-image-reserved")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.ActorImageReservedErr, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layouts.CSRF)
		}

		if in.Name == actorImage.Name {
			c.Append("HX-Trigger-After-Swap", "ptrcr:actor-image-reserved")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partials.ActorImageReserved, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layouts.CSRF)
		} else {
			c.Append("HX-Trigger-After-Swap", "ptrcr:actor-image-reserved")
			return c.Render(partials.ActorImageFree, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layouts.CSRF)
		}
	}
}
