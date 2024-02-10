package handler

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/util"
)

func ResendEmailVerification(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := util.GetID(c)
		if err != nil {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrNoID,
				layouts.None,
			)
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrInternal(id),
				layouts.None,
			)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrInternal(id),
				layouts.None,
			)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(
					partial.NoticeSectionError,
					partial.BindProfileEmailResendVerificationErrInternal(id),
					layouts.None,
				)
			}
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrInternal(id),
				layouts.None,
			)
		}

		if e.Verified {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationInfoConflict(id),
				layouts.None,
			)
		}
		if e.PID != pid {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrInternal(id),
				layouts.None,
			)
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrInternal(id),
				layouts.None,
			)
		}
		if err == nil && ve.Verified {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrForbiddenAlreadyVerified(id),
				layouts.None,
			)
		}

		if err = tx.Commit(); err != nil {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrInternal(id),
				layouts.None,
			)
		}

		if err = email.SendVerificationEmail(i, id, e.Address); err != nil {
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(
				partial.NoticeSectionError,
				partial.BindProfileEmailResendVerificationErrInternal(id),
				layouts.None,
			)
		}

		return c.Render(
			partial.NoticeSectionSuccess,
			partial.BindProfileEmailResendVerificationSuccess(id),
			layouts.None,
		)
	}
}
