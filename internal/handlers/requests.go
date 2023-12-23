package handlers

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/shared"
)

func RequestFieldPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lpid := c.Locals("pid")
		if lpid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
		}
		pid, ok := lpid.(int64)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
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

		field := c.Params("field")
		if len(field) == 0 {
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

		if req.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permissions.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
			}
			if !iperms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}
		}

		_, ok = request.FieldsByType[req.Type]
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		comments, err := qtx.ListCommentsForRequestWithAuthor(context.Background(), rid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := fiber.Map{}
		if req.Type == request.TypeCharacterApplication {
			app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
			if err != nil {
				// TODO: This means that a Request was created without content - this is an error
				// We should instead insert a blank content row here, but deal with this later
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			b = request.BindCharacterApplicationFieldPage(b, &app, field)
		} else {
			// TODO: This means that there's a request in the database with an invalid type
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		viewsByField, ok := request.ViewsByFieldAndType[req.Type]
		if !ok {
			// TODO: Again, noteworthy because either a bad type or a missing register
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		view, ok := viewsByField[field]
		if !ok {
			// TODO: Noteworthy to handle and track
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b = request.BindRequestFieldPage(b, request.BindRequestFieldPageParams{
			PID:      pid,
			Field:    field,
			Request:  &req,
			Comments: comments,
		})

		return c.Render(view, b, "views/layouts/requests")
	}
}

func UpdateRequestField(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in := new(request.UpdateInput)
		if err := c.BodyParser(in); err != nil {
			log.Println(err)
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		lpid := c.Locals("pid")
		if lpid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
		}
		pid, ok := lpid.(int64)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
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

		field, err := in.GetField()
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

		// Handle a status update
		if field == request.FieldStatus {
			if !request.IsStatusValid(in.Status) {
				c.Status(fiber.StatusBadRequest)
				return nil
			}

			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			perms, ok := lperms.(permissions.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			ok = request.IsStatusUpdateOK(&req, perms, pid, in.Status)
			if !ok {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			if err = qtx.UpdateRequestStatus(context.Background(), queries.UpdateRequestStatusParams{
				ID:     rid,
				Status: in.Status,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			c.Append("HX-Refresh", "true")
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !request.IsFieldValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = in.UpdateField(qtx, &req, field); err != nil {
			if err == request.ErrInvalidInput {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return nil
	}
}
