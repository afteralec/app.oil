package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func NewCharacterApplication(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		count, err := qtx.CountOpenCharacterApplicationsForPlayer(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if count >= shared.MaxOpenCharacterApplications {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		result, err := qtx.CreateRequest(context.Background(), queries.CreateRequestParams{
			PID:  pid.(int64),
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

		if err = qtx.CreateCharacterApplicationContent(context.Background(), rid); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		path := routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
		c.Append("HX-Redirect", path)
		return nil
	}
}

func UpdateCharacterApplicationName(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Name string `form:"name"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
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

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		name := character.SanitizeName(r.Name)
		if !character.IsNameValid(name) {
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

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		err = qtx.CreateHistoryForCharacterApplication(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = qtx.UpdateCharacterApplicationContentName(context.Background(), queries.UpdateCharacterApplicationContentNameParams{
			Name: name,
			RID:  rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Status != request.StatusReady && character.IsApplicationReady(&app) {
			if err = qtx.MarkRequestReady(context.Background(), rid); err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					return nil
				}
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationGender(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Gender string `form:"gender"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
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

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		err = qtx.CreateHistoryForCharacterApplication(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		if !character.IsGenderValid(r.Gender) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = qtx.UpdateCharacterApplicationContentGender(context.Background(), queries.UpdateCharacterApplicationContentGenderParams{
			Gender: r.Gender,
			RID:    rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Status != request.StatusReady && character.IsApplicationReady(&app) {
			if err = qtx.MarkRequestReady(context.Background(), rid); err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					return nil
				}
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		c.Append("HX-Redirect", routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
		return nil
	}
}

func UpdateCharacterApplicationShortDescription(i *shared.Interfaces) fiber.Handler {
	type input struct {
		ShortDescription string `form:"sdesc"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
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

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		err = qtx.CreateHistoryForCharacterApplication(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		if !character.IsShortDescriptionValid(r.ShortDescription) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = qtx.UpdateCharacterApplicationContentShortDescription(context.Background(), queries.UpdateCharacterApplicationContentShortDescriptionParams{
			ShortDescription: r.ShortDescription,
			RID:              rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Status != request.StatusReady && character.IsApplicationReady(&app) {
			if err = qtx.MarkRequestReady(context.Background(), rid); err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					return nil
				}
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		c.Append("HX-Redirect", routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
		return nil
	}
}

func UpdateCharacterApplicationDescription(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Description string `form:"description"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
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

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		err = qtx.CreateHistoryForCharacterApplication(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		if !character.IsDescriptionValid(r.Description) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = qtx.UpdateCharacterApplicationContentDescription(context.Background(), queries.UpdateCharacterApplicationContentDescriptionParams{
			Description: r.Description,
			RID:         rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Status != request.StatusReady && character.IsApplicationReady(&app) {
			if err = qtx.MarkRequestReady(context.Background(), rid); err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					return nil
				}
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationBackstory(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Backstory string `form:"backstory"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
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

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		err = qtx.CreateHistoryForCharacterApplication(context.Background(), rid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !character.IsBackstoryValid(r.Backstory) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = qtx.UpdateCharacterApplicationContentBackstory(context.Background(), queries.UpdateCharacterApplicationContentBackstoryParams{
			Backstory: r.Backstory,
			RID:       rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Status != request.StatusReady && character.IsApplicationReady(&app) {
			if err = qtx.MarkRequestReady(context.Background(), rid); err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					return nil
				}
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", routes.CharacterApplicationPath(strconv.FormatInt(rid, 10)))
		c.Status(fiber.StatusOK)
		return nil
	}
}
