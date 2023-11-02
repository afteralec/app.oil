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

const MaxEmailCount = 3

func AddEmail(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Email string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		ec, err := i.Queries.CountEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if ec >= MaxEmailCount {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			return c.Render("web/views/partials/profile/email/err-too-many-emails", &fiber.Map{}, "")
		}

		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		e, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			return c.Render("web/views/partials/profile/email/err-invalid-email", &fiber.Map{}, "")
		}

		_, err = i.Queries.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil {
			if err != sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}
		if err == nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "innerHTML")
			return c.Render("web/views/partials/profile/email/err-conflict", &fiber.Map{
				"Address": e.Address,
			}, "")
		}

		result, err := i.Queries.CreateEmail(
			context.Background(),
			queries.CreateEmailParams{Pid: pid.(int64), Address: e.Address},
		)
		if err != nil {
			if me, ok := err.(*mysql.MySQLError); ok {
				if me.Number == mysqlerr.ER_DUP_ENTRY {
					c.Append("HX-Retarget", "#add-email-error")
					c.Append("HX-Reswap", "innerHTML")
					return c.Render("web/views/partials/profile/email/err-conflict", &fiber.Map{
						"Address": e.Address,
					}, "")
				}
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = email.Verify(i.Redis, id, e.Address)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		return c.Render("web/views/partials/profile/email/new-email", &fiber.Map{
			"ID":      id,
			"Address": e.Address,
			"Created": true,
		}, "")
	}
}
