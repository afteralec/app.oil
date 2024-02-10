package handler

import (
	"context"
	"database/sql"
	"net/mail"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/queries"
)

func AddEmail(i *interfaces.Shared) fiber.Handler {
	type request struct {
		Email string `form:"email"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrUnauthorized, layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInternal, layouts.None)
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		ec, err := qtx.CountEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInternal, layouts.None)
		}

		if ec >= email.MaxCount {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrTooMany(), layouts.None)
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInvalid, layouts.None)
		}

		e, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInvalid, layouts.None)
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInternal, layouts.None)
		}
		if err == nil && ve.Verified {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrConflict(e.Address), layouts.None)
		}

		result, err := qtx.CreateEmail(
			context.Background(),
			queries.CreateEmailParams{PID: pid.(int64), Address: e.Address},
		)
		if err != nil {
			if me, ok := err.(*mysql.MySQLError); ok {
				if me.Number == mysqlerr.ER_DUP_ENTRY {
					c.Append("HX-Retarget", "#add-email-error")
					c.Append("HX-Reswap", "outerHTML")
					c.Append(header.HXAcceptable, "true")
					c.Status(fiber.StatusConflict)
					return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrConflict(e.Address), layouts.None)
				}
			}
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInternal, layouts.None)
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInternal, layouts.None)
		}

		if err = tx.Commit(); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInternal, layouts.None)
		}

		if err = email.SendVerificationEmail(i, id, e.Address); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindProfileAddEmailErrInternal, layouts.None)
		}

		c.Status(fiber.StatusCreated)
		return c.Render(partials.ProfileEmailNew, &fiber.Map{
			"ID":      id,
			"Address": e.Address,
			"Created": true,
		}, "")
	}
}
