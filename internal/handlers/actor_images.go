package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/permissions"
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
		b["ActorImages"] = pageActorImages
		b["PageHeader"] = fiber.Map{
			"Title":    "Actor Images",
			"SubTitle": "Actor images are where the primary properties for an actor are defined, like a template",
		}
		return c.Render(views.ActorImages, b)
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
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			return c.Render(partials.ActorImageFree, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layouts.CSRF)
		}
	}
}
