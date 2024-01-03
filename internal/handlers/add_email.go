package handlers

import (
	"context"
	"database/sql"
	"net/mail"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/views"
)

func AddEmail(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Email string `form:"email"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.PartialProfileEmailErrInternal, &fiber.Map{}, "")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.PartialProfileEmailErrInternal, &fiber.Map{}, "")
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		ec, err := qtx.CountEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.PartialProfileEmailErrInternal, &fiber.Map{}, "")
		}

		if ec >= shared.MaxEmailCount {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(views.PartialProfileEmailErrTooMany, &fiber.Map{}, "")
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(views.PartialProfileEmailErrInvalid, &fiber.Map{}, "")
		}

		e, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(views.PartialProfileEmailErrInvalid, &fiber.Map{}, "")
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.PartialProfileEmailErrInternal, &fiber.Map{}, "")
		}
		if err == nil && ve.Verified {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(views.PartialProfileEmailErrConflict, &fiber.Map{
				"Address": e.Address,
			}, "")
		}

		result, err := qtx.CreateEmail(
			context.Background(),
			queries.CreateEmailParams{PID: pid.(int64), Address: e.Address},
		)
		if err != nil {
			if me, ok := err.(*mysql.MySQLError); ok {
				if me.Number == mysqlerr.ER_DUP_ENTRY {
					c.Append("HX-Retarget", "#add-email-error")
					c.Append("HX-Reswap", "innerHTML")
					c.Append(shared.HeaderHXAcceptable, "true")
					c.Status(fiber.StatusConflict)
					return c.Render(views.PartialProfileEmailErrConflict, &fiber.Map{
						"Address": e.Address,
					}, "")
				}
			}
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.PartialProfileEmailErrInternal, &fiber.Map{}, "")
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.PartialProfileEmailErrInternal, &fiber.Map{}, "")
		}

		if err = tx.Commit(); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.PartialProfileEmailErrInternal, &fiber.Map{}, "")
		}

		if err = email.SendVerificationEmail(i, id, e.Address); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.PartialProfileEmailErrInternal, &fiber.Map{}, "")
		}

		c.Status(fiber.StatusCreated)
		return c.Render(views.PartialProfileEmailNewUnverified, &fiber.Map{
			"ID":      id,
			"Address": e.Address,
			"Created": true,
		}, "")
	}
}
