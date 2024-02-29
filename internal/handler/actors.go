package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/actor"
	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

func ActorImagesPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		if !perms.HasPermission(player.PermissionViewAllActorImages.Name) {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		actorImages, err := i.Queries.ListActorImages(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		pageActorImages := []fiber.Map{}
		for _, actorImage := range actorImages {
			var sb strings.Builder
			fmt.Fprintf(&sb, "[%d] %s", actorImage.ID, actorImage.ShortDescription)
			pageActorImage := fiber.Map{
				"Title": sb.String(),
				"Name":  actorImage.Name,
				"Path":  route.ActorImagePath(actorImage.ID),
			}

			if perms.HasPermission(player.PermissionCreateActorImage.Name) {
				pageActorImage["EditPath"] = route.EditActorImagePath(actorImage.ID)
			}

			pageActorImages = append(pageActorImages, pageActorImage)
		}

		b := view.Bind(c)
		if perms.HasPermission(player.PermissionCreateActorImage.Name) {
			b["CreatePermission"] = true
		}
		b["ActorImages"] = pageActorImages
		b["PageHeader"] = fiber.Map{
			"Title":    "Actor Images",
			"SubTitle": "Actor images are where the primary properties for an actor are defined, like a template",
		}
		return c.Render(view.ActorImages, b)
	}
}

func EditActorImagePage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		aiid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		if !perms.HasPermission(player.PermissionCreateActorImage.Name) {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		actorImage, err := qtx.GetActorImage(context.Background(), aiid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		b := view.Bind(c)
		// TODO: Get a bind function for this
		b["NavBack"] = fiber.Map{
			"Path":  route.ActorImages,
			"Label": "Back to Actor Images",
		}
		// TODO: Get a bind function for this too
		b["PageHeader"] = fiber.Map{
			"Title":    actor.ImageTitleWithID(actorImage.Name, actorImage.ID),
			"SubTitle": "Update actor properties here",
		}
		// TODO: Write a bind function for this
		b["Name"] = actorImage.Name
		b["ShortDescription"] = actorImage.ShortDescription
		b["Description"] = actorImage.Description
		b["ShortDescriptionPath"] = route.ActorImageShortDescriptionPath(aiid)
		b["DescriptionPath"] = route.ActorImageDescriptionPath(aiid)
		return c.Render(view.EditActorImage, b)
	}
}

func ActorImagePage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		aiid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		if !perms.HasPermission(player.PermissionViewAllActorImages.Name) {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		actorImage, err := qtx.GetActorImage(context.Background(), aiid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		b := view.Bind(c)
		// TODO: Get a bind function for this
		b["NavBack"] = fiber.Map{
			"Path":  route.ActorImages,
			"Label": "Back to Actor Images",
		}
		// TODO: Get a bind function for this too
		b["PageHeader"] = fiber.Map{
			"Title": actor.ImageTitleWithID(actorImage.Name, actorImage.ID),
		}
		b["Name"] = actorImage.Name
		b["ShortDescription"] = actorImage.ShortDescription
		b["Description"] = actorImage.Description
		return c.Render(view.ActorImage, b)
	}
}

func NewActorImage(i *service.Interfaces) fiber.Handler {
	type input struct {
		Name string `form:"name"`
	}

	const sectionID string = "actor-image-create-error"

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		if !actor.IsImageNameValid(actor.SanitizeImageName(in.Name)) {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"The Image Name you entered isn't valid. Please try again.",
				},
				NoticeIcon: true,
			}), layout.None)
		}

		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"It looks like your session may have expired.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		if !perms.HasPermission(player.PermissionCreateActorImage.Name) {
			c.Status(fiber.StatusForbidden)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"You don't have the permission(s) necessary to create an Actor Image.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		result, err := qtx.CreateActorImage(context.Background(), query.CreateActorImageParams{
			Name:             in.Name,
			ShortDescription: actor.DefaultImageShortDescription,
			Description:      actor.DefaultImageDescription,
			Gender:           actor.DefaultImageGender,
		})
		if err != nil {
			if me, ok := err.(*mysql.MySQLError); ok {
				if me.Number == mysqlerr.ER_DUP_ENTRY {
					c.Status(fiber.StatusConflict)
					c.Append(header.HXAcceptable, "true")
					c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
					return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
						SectionID:    sectionID,
						SectionClass: "pt-2",
						NoticeText: []string{
							"That Actor Image name is already in use. Please choose another.",
						},
						NoticeIcon: true,
					}), layout.None)
				}
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		aiid, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		c.Status(fiber.StatusCreated)
		c.Append("HX-Redirect", route.EditActorImagePath(aiid))
		c.Append("HX-Reswap", "none")
		return nil
	}
}

func EditActorImageShortDescription(i *service.Interfaces) fiber.Handler {
	type input struct {
		ShortDescription string `form:"sdesc"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !actor.IsShortDescriptionValid(actor.SanitizeShortDescription(in.ShortDescription)) {
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
		if !perms.HasPermission(player.PermissionCreateActorImage.Name) {
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

		if err := qtx.UpdateActorImageShortDescription(context.Background(), query.UpdateActorImageShortDescriptionParams{
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
		b["ShortDescriptionPath"] = route.ActorImageShortDescriptionPath(actorImage.ID)
		b["NoticeSection"] = partial.BindNoticeSection(partial.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "actor-image-edit-short-description-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The short description has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partial.ActorImageEditShortDescription, b, layout.None)
	}
}

func EditActorImageDescription(i *service.Interfaces) fiber.Handler {
	type input struct {
		Description string `form:"desc"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !actor.IsDescriptionValid(actor.SanitizeDescription(in.Description)) {
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
		if !perms.HasPermission(player.PermissionCreateActorImage.Name) {
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

		if err := qtx.UpdateActorImageDescription(context.Background(), query.UpdateActorImageDescriptionParams{
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
		b["DescriptionPath"] = route.ActorImageDescriptionPath(actorImage.ID)
		b["NoticeSection"] = partial.BindNoticeSection(partial.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "actor-image-edit-description-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The description has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partial.ActorImageEditDescription, b, layout.None)
	}
}

func ActorImageNameReserved(i *service.Interfaces) fiber.Handler {
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
				return c.Render(partial.ActorImageFree, fiber.Map{
					"CSRF": c.Locals("csrf"),
				}, layout.CSRF)
			}
			c.Append("HX-Trigger-After-Swap", "ptrcr:actor-image-reserved")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.ActorImageReservedErr, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		}

		if in.Name == actorImage.Name {
			c.Append("HX-Trigger-After-Swap", "ptrcr:actor-image-reserved")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partial.ActorImageReserved, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		} else {
			c.Append("HX-Trigger-After-Swap", "ptrcr:actor-image-reserved")
			return c.Render(partial.ActorImageFree, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		}
	}
}
