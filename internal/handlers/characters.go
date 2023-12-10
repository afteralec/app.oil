package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/shared"
)

func CharactersPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		reqs, err := i.Queries.ListCharacterApplicationsForPlayer(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"))
		}

		b := c.Locals("bind").(fiber.Map)
		b["CharacterApplications"] = reqs
		b["HasCharacterApplications"] = len(reqs) > 0
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
			Pid:  pid.(int64),
			Type: request.TypeCharacterApplication,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rid, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
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
			return nil
		}

		c.Status(fiber.StatusCreated)
		// TODO: Get this in a generator
		path := fmt.Sprintf("/characters/new/%d/name", rid)
		c.Append("HX-Redirect", path)
		return nil
	}
}

func UpdateCharacterApplication(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type request struct {
		Name        string `form:"name"`
		Gender      string `form:"gender"`
		Sdesc       string `form:"sdesc"`
		Description string `form:"description"`
		Backstory   string `form:"backstory"`
	}
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
				// TODO: Pull the request here and if the type isn't CharacterApplication, send back a 400
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:      app.Gender,
			Name:        app.Name,
			Sdesc:       app.Sdesc,
			Description: app.Description,
			Backstory:   app.Backstory,
			Vid:         app.Vid,
			Rid:         app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContent(context.Background(), queries.UpdateCharacterApplicationContentParams{
			// TODO: Get gender into a constant
			Gender:      r.Gender,
			Name:        r.Name,
			Sdesc:       r.Sdesc,
			Description: r.Description,
			Backstory:   r.Backstory,
			Vid:         app.Vid,
			Rid:         app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationName(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type request struct {
		Name string `form:"name"`
	}
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
				// TODO: Pull the request here and if the type isn't CharacterApplication, send back a 400
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:      app.Gender,
			Name:        app.Name,
			Sdesc:       app.Sdesc,
			Description: app.Description,
			Backstory:   app.Backstory,
			Vid:         app.Vid,
			Rid:         app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentName(context.Background(), queries.UpdateCharacterApplicationContentNameParams{
			Name: r.Name,
			Rid:  rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationGender(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type request struct {
		Gender string `form:"gender"`
	}
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
				// TODO: Pull the request here and if the type isn't CharacterApplication, send back a 400
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:      app.Gender,
			Name:        app.Name,
			Sdesc:       app.Sdesc,
			Description: app.Description,
			Backstory:   app.Backstory,
			Vid:         app.Vid,
			Rid:         app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentGender(context.Background(), queries.UpdateCharacterApplicationContentGenderParams{
			Gender: r.Gender,
			Rid:    rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationShortDescription(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type request struct {
		Sdesc string `form:"sdesc"`
	}
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
				// TODO: Pull the request here and if the type isn't CharacterApplication, send back a 400
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:      app.Gender,
			Name:        app.Name,
			Sdesc:       app.Sdesc,
			Description: app.Description,
			Backstory:   app.Backstory,
			Vid:         app.Vid,
			Rid:         app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentSdesc(context.Background(), queries.UpdateCharacterApplicationContentSdescParams{
			Sdesc: r.Sdesc,
			Rid:   rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationDescription(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type request struct {
		Description string `form:"description"`
	}
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
				// TODO: Pull the request here and if the type isn't CharacterApplication, send back a 400
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:      app.Gender,
			Name:        app.Name,
			Sdesc:       app.Sdesc,
			Description: app.Description,
			Backstory:   app.Backstory,
			Vid:         app.Vid,
			Rid:         app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentDescription(context.Background(), queries.UpdateCharacterApplicationContentDescriptionParams{
			Description: r.Description,
			Rid:         rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationBackstory(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type request struct {
		Backstory string `form:"backstory"`
	}
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
				// TODO: Pull the request here and if the type isn't CharacterApplication, send back a 400
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:      app.Gender,
			Name:        app.Name,
			Sdesc:       app.Sdesc,
			Description: app.Description,
			Backstory:   app.Backstory,
			Vid:         app.Vid,
			Rid:         app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentBackstory(context.Background(), queries.UpdateCharacterApplicationContentBackstoryParams{
			Backstory: r.Backstory,
			Rid:       rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationVersion(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input on the way in
	type request struct {
		Vid int64 `form:"vid"`
	}
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
				// TODO: Pull the request here and if the type isn't CharacterApplication, send back a 400
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:      app.Gender,
			Name:        app.Name,
			Sdesc:       app.Sdesc,
			Description: app.Description,
			Backstory:   app.Backstory,
			Vid:         app.Vid,
			Rid:         app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentVersion(context.Background(), queries.UpdateCharacterApplicationContentVersionParams{
			Vid: r.Vid,
			Rid: rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}
