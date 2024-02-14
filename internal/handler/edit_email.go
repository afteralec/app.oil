package handler

import (
	"context"
	"database/sql"
	"net/mail"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/query"
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
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInvalid, layout.None)
		}

		ne, err := mail.ParseAddress(in.Email)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInvalid, layout.None)
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrUnauthorized, layout.None)
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInvalid, layout.None)
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInvalid, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "#profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
			}
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		if e.PID != pid {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		if !e.Verified {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		if e.Address == ne.Address {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrConflictSame(ne.Address), layout.None)
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), ne.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}
		if err == nil && ve.Verified {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrConflict(ne.Address), layout.None)
		}

		if err := qtx.DeleteEmail(context.Background(), id); err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "#profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
			}
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		result, err := qtx.CreateEmail(context.Background(), query.CreateEmailParams{
			Address: ne.Address,
			PID:     pid,
		})
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		id, err = result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		return c.Render(partial.ProfileEmailUnverified, &fiber.Map{
			"ID":       id,
			"Address":  ne.Address,
			"Verified": false,
		}, "")
	}
}
