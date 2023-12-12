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
			return c.Render("web/views/partials/profile/email/err-internal", &fiber.Map{}, "")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/err-internal", &fiber.Map{}, "")
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		ec, err := qtx.CountEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/err-internal", &fiber.Map{}, "")
		}

		if ec >= shared.MaxEmailCount {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render("web/views/partials/profile/email/err-too-many-emails", &fiber.Map{}, "")
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("web/views/partials/profile/email/err-invalid-email", &fiber.Map{}, "")
		}

		e, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("web/views/partials/profile/email/err-invalid-email", &fiber.Map{}, "")
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/err-internal", &fiber.Map{}, "")
		}
		if err == nil && ve.Verified {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render("web/views/partials/profile/email/err-conflict", &fiber.Map{
				"Address": e.Address,
			}, "")
		}

		result, err := qtx.CreateEmail(
			context.Background(),
			queries.CreateEmailParams{Pid: pid.(int64), Address: e.Address},
		)
		if err != nil {
			if me, ok := err.(*mysql.MySQLError); ok {
				if me.Number == mysqlerr.ER_DUP_ENTRY {
					c.Append("HX-Retarget", "#add-email-error")
					c.Append("HX-Reswap", "innerHTML")
					c.Append(shared.HeaderHXAcceptable, "true")
					c.Status(fiber.StatusConflict)
					return c.Render("web/views/partials/profile/email/err-conflict", &fiber.Map{
						"Address": e.Address,
					}, "")
				}
			}
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/err-internal", &fiber.Map{}, "")
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/err-internal", &fiber.Map{}, "")
		}

		if err = tx.Commit(); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/err-internal", &fiber.Map{}, "")
		}

		if err = email.SendVerificationEmail(i, id, e.Address); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/err-internal", &fiber.Map{}, "")
		}

		c.Status(fiber.StatusCreated)
		return c.Render("web/views/partials/profile/email/unverified/new", &fiber.Map{
			"ID":      id,
			"Address": e.Address,
			"Created": true,
		}, "")
	}
}
