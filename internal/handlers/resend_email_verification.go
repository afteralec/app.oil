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
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			return c.Render("web/views/partials/profile/email/resend-404", &fiber.Map{
				"CSRF": c.Locals("csrf"),
				"ID":   id,
			}, "web/views/layouts/csrf")
		}
		if e.Verified {
			return c.Render("web/views/partials/profile/email/resend-conflict", &fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		}
		if e.Pid != pid.(int64) {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		}

		err = email.Verify(i.Redis, id, e.Email)
		if err != nil {
			return c.Render("web/views/partials/profile/email/resend-err", &fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, "web/views/layouts/csrf")
		}

		return c.Render("web/views/partials/profile/email/resend-success", &fiber.Map{
			"CSRF": c.Locals("csrf"),
		}, "web/views/layouts/csrf")
	}
}
