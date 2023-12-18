package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/bind"
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
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
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

		row, err := qtx.GetCharacterApplication(context.Background(), rid)
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

		if row.Request.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if row.Request.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}
		}

		// TODO: After putting an application in review,
		// This should check if the request is in review
		// and then show the review version of this page.
		// The review version of this page essentially just allows
		// the reviewer to comment and finish the review.

		b := bind.RequestStatus(c.Locals(bind.Name).(fiber.Map), &row.Request)
		b = bind.RequestViewedBy(b, &row.Request, pid.(int64))
		b = bind.CharacterApplicationPaths(b, &row.Request)
		b = bind.CharacterApplicationContent(b, &row.CharacterApplicationContent)
		b["CharacterApplicationParts"] = character.MakeApplicationParts("name", &row.CharacterApplicationContent)
		b["NextLink"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10))

		if !request.IsEditable(&row.Request) {
			return c.Render("views/character/application/name/view", b)
		}

		if row.Request.PID != pid {
			return c.Render("views/character/application/name/view", b)
		}

		return c.Render("views/character/application/name/edit", b)
	}
}

func CharacterApplicationGenderPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
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

		row, err := qtx.GetCharacterApplication(context.Background(), rid)
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

		if row.Request.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if row.Request.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}
		}

		b := bind.RequestStatus(c.Locals(bind.Name).(fiber.Map), &row.Request)
		b = bind.RequestViewedBy(b, &row.Request, pid.(int64))
		b = bind.CharacterApplicationPaths(b, &row.Request)
		b = bind.CharacterApplicationContent(b, &row.CharacterApplicationContent)
		b = bind.CharacterApplicationGender(b, &row.CharacterApplicationContent)
		b["CharacterApplicationParts"] = character.MakeApplicationParts("gender", &row.CharacterApplicationContent)
		b["BackLink"] = routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
		b["NextLink"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))

		if !request.IsEditable(&row.Request) {
			return c.Render("views/character/application/gender/view", b)
		}

		if row.Request.PID != pid {
			return c.Render("views/character/application/gender/view", b)
		}

		return c.Render("views/character/application/gender/edit", b)
	}
}

func CharacterApplicationShortDescriptionPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
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

		row, err := qtx.GetCharacterApplication(context.Background(), rid)
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

		if row.Request.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if row.Request.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}
		}

		b := bind.RequestStatus(c.Locals(bind.Name).(fiber.Map), &row.Request)
		b = bind.RequestViewedBy(b, &row.Request, pid.(int64))
		b = bind.CharacterApplicationPaths(b, &row.Request)
		b = bind.CharacterApplicationContent(b, &row.CharacterApplicationContent)
		b["CharacterApplicationParts"] = character.MakeApplicationParts("sdesc", &row.CharacterApplicationContent)
		b["BackLink"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10))
		b["NextLink"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10))

		if !request.IsEditable(&row.Request) {
			return c.Render("views/character/application/sdesc/view", b)
		}

		if row.Request.PID != pid {
			return c.Render("views/character/application/sdesc/view", b)
		}

		return c.Render("views/character/application/sdesc/edit", b)
	}
}

func CharacterApplicationDescriptionPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
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

		row, err := qtx.GetCharacterApplication(context.Background(), rid)
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

		if row.Request.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if row.Request.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}
		}

		b := bind.RequestStatus(c.Locals(bind.Name).(fiber.Map), &row.Request)
		b = bind.RequestViewedBy(b, &row.Request, pid.(int64))
		b = bind.CharacterApplicationPaths(b, &row.Request)
		b = bind.CharacterApplicationContent(b, &row.CharacterApplicationContent)
		b["CharacterApplicationParts"] = character.MakeApplicationParts("description", &row.CharacterApplicationContent)
		b["BackLink"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
		b["NextLink"] = routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10))

		if !request.IsEditable(&row.Request) {
			return c.Render("views/character/application/description/view", b)
		}

		if row.Request.PID != pid {
			return c.Render("views/character/application/description/view", b)
		}

		return c.Render("views/character/application/description/edit", b)
	}
}

func CharacterApplicationBackstoryPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
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

		row, err := qtx.GetCharacterApplication(context.Background(), rid)
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

		if row.Request.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if row.Request.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}
		}

		b := bind.RequestStatus(c.Locals(bind.Name).(fiber.Map), &row.Request)
		b = bind.RequestViewedBy(b, &row.Request, pid.(int64))
		b = bind.CharacterApplicationPaths(b, &row.Request)
		b = bind.CharacterApplicationContent(b, &row.CharacterApplicationContent)
		b["CharacterApplicationParts"] = character.MakeApplicationParts("backstory", &row.CharacterApplicationContent)
		b["BackLink"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10))
		b["NextLink"] = routes.CharacterApplicationSummaryPath(strconv.FormatInt(rid, 10))

		if !request.IsEditable(&row.Request) {
			return c.Render("views/character/application/backstory/view", b)
		}

		if row.Request.PID != pid {
			return c.Render("views/character/application/backstory/view", b)
		}

		return c.Render("views/character/application/backstory/edit", b)
	}
}

func CharacterApplicationSummaryPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
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

		row, err := qtx.GetCharacterApplication(context.Background(), rid)
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

		if row.Request.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if row.Request.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permission.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
			}
			if !iperms.Permissions[permission.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}
		}

		b := bind.RequestStatus(c.Locals(bind.Name).(fiber.Map), &row.Request)
		b = bind.RequestViewedBy(b, &row.Request, pid.(int64))
		b = bind.CharacterApplicationPaths(b, &row.Request)
		b = bind.CharacterApplicationContent(b, &row.CharacterApplicationContent)
		b["CharacterApplicationParts"] = character.MakeApplicationParts("summary", &row.CharacterApplicationContent)
		b["BackLink"] = routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10))

		if !request.IsEditable(&row.Request) {
			return c.Render("views/character/application/backstory/view", b)
		}

		if row.Request.PID != pid {
			return c.Render("views/character/application/backstory/view", b)
		}

		return c.Render("views/character/application/backstory/edit", b)
	}
}
