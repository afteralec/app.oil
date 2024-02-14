package handler

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
)

func DeleteEmail(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrUnauthorized, layout.None)
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
			}
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		if e.PID != pid.(int64) {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		if err := qtx.DeleteEmail(context.Background(), id); err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
			}
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		return c.Render(partial.ProfileEmailDeleteSuccess, &fiber.Map{
			"ID":      e.ID,
			"Address": e.Address,
		}, "")
	}
}
