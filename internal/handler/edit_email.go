package handler

import (
	"context"
	"database/sql"
	"net/mail"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/headers"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/util"
)

func EditEmail(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Email string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInvalid, layouts.None)
		}

		ne, err := mail.ParseAddress(in.Email)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInvalid, layouts.None)
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrUnauthorized, layouts.None)
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInvalid, layouts.None)
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInvalid, layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "#profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(headers.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
			}
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
		}

		if e.PID != pid {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
		}

		if !e.Verified {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
		}

		if e.Address == ne.Address {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrConflictSame(ne.Address), layouts.None)
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), ne.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInternal, layouts.None)
		}
		if err == nil && ve.Verified {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrConflict(ne.Address), layouts.None)
		}

		if err := qtx.DeleteEmail(context.Background(), id); err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "#profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(headers.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
			}
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
		}

		result, err := qtx.CreateEmail(context.Background(), queries.CreateEmailParams{
			Address: ne.Address,
			PID:     pid,
		})
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
		}

		id, err = result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(headers.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileEditEmailErrInternal, layouts.None)
		}

		return c.Render(partials.ProfileEmailUnverified, &fiber.Map{
			"ID":       id,
			"Address":  ne.Address,
			"Verified": false,
		}, "")
	}
}
