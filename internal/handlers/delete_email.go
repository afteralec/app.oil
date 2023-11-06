package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
)

func DeleteEmail(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/partials/profile/email/delete/err-401", &fiber.Map{}, "")
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("web/views/partials/profile/email/delete/err-internal", &fiber.Map{}, "")
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("web/views/partials/profile/email/delete/err-internal", &fiber.Map{}, "")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/delete/err-internal", &fiber.Map{}, "")
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: Make this a 404 error?
				c.Append("HX-Retarget", "profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render("web/views/partials/profile/email/delete/err-internal", &fiber.Map{}, "")
			}
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/delete/err-internal", &fiber.Map{}, "")
		}

		if e.Pid != pid.(int64) {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render("web/views/partials/profile/email/delete/err-internal", &fiber.Map{}, "")
		}

		_, err = qtx.DeleteEmail(context.Background(), id)
		if err != nil {
			// TODO: Differentiate between not found and other errors
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/delete/err-internal", &fiber.Map{}, "")
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/delete/err-internal", &fiber.Map{}, "")
		}

		return c.Render("web/views/partials/profile/email/delete/success", &fiber.Map{
			"ID":      e.ID,
			"Address": e.Address,
		}, "")
	}
}
