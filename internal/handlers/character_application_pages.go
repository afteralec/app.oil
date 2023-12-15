package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func CharacterApplicationNamePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
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

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			statuses := character.MakeApplicationPartStatuses("name", &app)
			b := c.Locals(shared.Bind).(fiber.Map)
			b["Name"] = app.Name
			b["CharacterApplicationNamePath"] = routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
			b["Statuses"] = statuses
			b["NextLink"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10))
			return c.Render("views/character/application/name/view", b, "views/layouts/standalone")
		}

		statuses := character.MakeApplicationPartStatuses("name", &app)
		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["CharacterApplicationNamePath"] = routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
		b["Statuses"] = statuses
		return c.Render("views/character/application/name/edit", b, "views/layouts/standalone")
	}
}

func CharacterApplicationGenderPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
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

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			statuses := character.MakeApplicationPartStatuses("gender", &app)
			b := c.Locals(shared.Bind).(fiber.Map)
			b["Name"] = app.Name
			b["Gender"] = app.Gender
			b["CharacterApplicationNamePath"] = routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
			b["Statuses"] = statuses
			b["BackLink"] = routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
			b["NextLink"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
			return c.Render("views/character/application/gender/view", b, "views/layouts/standalone")
		}

		statuses := character.MakeApplicationPartStatuses("gender", &app)
		gender := character.SanitizeGender(app.Gender)
		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["GenderNonBinary"] = character.GenderNonBinary
		b["GenderFemale"] = character.GenderFemale
		b["GenderMale"] = character.GenderMale
		b["Gender"] = gender
		b["GenderIsNonBinary"] = gender == character.GenderNonBinary
		b["GenderIsFemale"] = gender == character.GenderFemale
		b["GenderIsMale"] = gender == character.GenderMale
		b["CharacterApplicationGenderPath"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10))
		b["Statuses"] = statuses
		b["BackLink"] = routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
		return c.Render("views/character/application/gender/edit", b, "views/layouts/standalone")
	}
}

func CharacterApplicationShortDescriptionPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
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

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			statuses := character.MakeApplicationPartStatuses("sdesc", &app)
			b := c.Locals(shared.Bind).(fiber.Map)
			b["Name"] = app.Name
			b["ShortDescription"] = app.ShortDescription
			b["Statuses"] = statuses
			b["BackLink"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10))
			b["NextLink"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10))
			return c.Render("views/character/application/sdesc/view", b, "views/layouts/standalone")
		}

		statuses := character.MakeApplicationPartStatuses("sdesc", &app)
		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["ShortDescription"] = app.ShortDescription
		b["CharacterApplicationShortDescriptionPath"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
		b["Statuses"] = statuses
		b["BackLink"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10))
		return c.Render("views/character/application/sdesc/edit", b, "views/layouts/standalone")
	}
}

func CharacterApplicationDescriptionPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
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

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			statuses := character.MakeApplicationPartStatuses("description", &app)
			b := c.Locals(shared.Bind).(fiber.Map)
			b["Name"] = app.Name
			b["Description"] = app.Description
			b["Statuses"] = statuses
			b["BackLink"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
			b["NextLink"] = routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10))
			return c.Render("views/character/application/description/view", b, "views/layouts/standalone")
		}

		statuses := character.MakeApplicationPartStatuses("description", &app)
		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["Description"] = app.Description
		b["CharacterApplicationDescriptionPath"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10))
		b["Statuses"] = statuses
		b["BackLink"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
		return c.Render("views/character/application/description/edit", b, "views/layouts/standalone")
	}
}

func CharacterApplicationBackstoryPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
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

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			statuses := character.MakeApplicationPartStatuses("backstory", &app)
			b := c.Locals(shared.Bind).(fiber.Map)
			b["Name"] = app.Name
			b["Backstory"] = app.Backstory
			b["Statuses"] = statuses
			b["BackLink"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10))
			return c.Render("views/character/application/backstory/view", b, "views/layouts/standalone")
		}

		statuses := character.MakeApplicationPartStatuses("backstory", &app)
		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["Backstory"] = app.Backstory
		b["Statuses"] = statuses
		b["BackLink"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10))
		return c.Render("views/character/application/backstory/edit", b, "views/layouts/standalone")
	}
}

func CharacterApplicationReviewPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
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

		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
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

		statuses := character.MakeApplicationPartStatuses("review", &app)
		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["Statuses"] = statuses
		b["Ready"] = character.IsApplicationReady(&app)
		b["BackLink"] = routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10))
		b["SubmitCharacterApplicationPath"] = routes.SubmitCharacterApplicationPath(strconv.FormatInt(rid, 10))
		return c.Render("views/character/application/review", b, "views/layouts/standalone")
	}
}
