package handlers

import (
	"context"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/shared"
)

// TODO: Do error partials for the unhappy paths here
func ResendEmailVerification(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return nil
		}
		if e.Verified {
			return c.Render("web/views/partials/profile/email/resend-conflict", &fiber.Map{}, "")
		}
		if e.Pid != pid.(int64) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		err = email.Verify(i.Redis, id, e.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return c.Render("web/views/partials/profile/email/resend-success", &fiber.Map{}, "")
	}
}
