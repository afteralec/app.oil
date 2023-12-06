package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
)

const (
	CharactersRoute              = "/characters"
	NewCharacterRoute            = "/characters/new"
	NewCharacterNameRoute        = "/characters/new/:id/name"
	NewCharacterGenderRoute      = "/characters/new/:id/gender"
	NewCharacterSdescRoute       = "/characters/new/:id/sdesc"
	NewCharacterDescriptionRoute = "/characters/new/:id/description"
	NewCharacterBackstoryRoute   = "/characters/new/:id/backstory"
)

func CharactersPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		apps, err := i.Queries.ListCharacterApplicationsForPlayer(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"))
		}

		b := c.Locals("bind").(fiber.Map)
		b["CharacterApplications"] = apps
		b["HasCharacterApplications"] = len(apps) > 0
		return c.Render("web/views/characters", b)
	}
}

func CharacterNamePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := c.Locals("bind").(fiber.Map)
		b["Name"] = app.Name
		return c.Render("web/views/characters/new/name", b, "web/views/layouts/standalone")
	}
}

func CharacterGenderPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := c.Locals("bind").(fiber.Map)
		b["Gender"] = app.Gender
		return c.Render("web/views/characters/new/gender", b, "web/views/layouts/standalone")
	}
}

func CharacterSdescPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := c.Locals("bind").(fiber.Map)
		b["Sdesc"] = app.Sdesc
		return c.Render("web/views/characters/new/gender", b, "web/views/layouts/standalone")
	}
}

func CharacterDescriptionPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := c.Locals("bind").(fiber.Map)
		b["Description"] = app.Description
		return c.Render("web/views/characters/new/gender", b, "web/views/layouts/standalone")
	}
}

func CharacterBackstoryPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := c.Locals("bind").(fiber.Map)
		b["Backstory"] = app.Backstory
		return c.Render("web/views/characters/new/gender", b, "web/views/layouts/standalone")
	}
}

func NewCharacterApplication(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		result, err := i.Queries.CreateRequest(context.Background(), queries.CreateRequestParams{
			Pid: pid.(int64),
			// TODO: Get this type into a constant
			Type: "CharacterApplication",
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Println("Error creating request")
			log.Println(err)
			return nil
		}

		rid, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Println("Error getting LastInsertId")
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContent(context.Background(), queries.CreateCharacterApplicationContentParams{
			// TODO: Get gender into a constant
			Gender:      "NonBinary",
			Name:        "",
			Sdesc:       "",
			Description: "",
			Backstory:   "",
			Rid:         rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Println("Error creating content")
			return nil
		}

		c.Status(fiber.StatusCreated)
		// TODO: Get this in a generator
		path := fmt.Sprintf("/character/new/%d/name", rid)
		c.Append("HX-Redirect", path)
		return nil
	}
}
