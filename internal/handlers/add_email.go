package handlers

import (
	"context"
	"net/mail"
	"time"

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

		ec, err := i.Queries.CountPlayerEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if ec >= MaxEmailCount {
			c.Append("HX-Retarget", "#add-email-error")
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
			return c.Render("web/views/partials/profile/email/err-invalid-email", &fiber.Map{}, "")
		}

		result, err := i.Queries.CreatePlayerEmail(
			context.Background(),
			queries.CreatePlayerEmailParams{Pid: pid.(int64), Email: e.Address},
		)
		if err != nil {
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

		time.Sleep(3 * time.Second)

		c.Status(fiber.StatusCreated)
		return c.Render("web/views/partials/profile/email/new-email", &fiber.Map{
			"ID":      id,
			"Email":   e.Address,
			"Created": true,
		}, "")
	}
}
