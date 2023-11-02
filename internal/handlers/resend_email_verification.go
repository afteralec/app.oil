package handlers

import (
	"context"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/shared"
)

func ResendEmailVerification(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Append("HX-Refresh", "true")
			return nil
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{}, "")
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{}, "")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{}, "")
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			return c.Render("web/views/partials/profile/email/resend-404", &fiber.Map{
				"ID": id,
			}, "")
		}
		if e.Verified {
			return c.Render("web/views/partials/profile/email/resend-conflict", &fiber.Map{}, "")
		}
		if e.Pid != pid.(int64) {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{}, "")
		}

		err = email.Verify(i.Redis, id, e.Address)
		if err != nil {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{}, "")
		}

		return c.Render("web/views/partials/profile/email/resend-success", &fiber.Map{}, "")
	}
}
