package handlers

import (
	"context"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
)

func EditEmail(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Email string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

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

		if !e.Verified {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if e.Pid != pid.(int64) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		_, err = qtx.DeleteEmail(context.Background(), id)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		result, err := qtx.CreatePlayerEmail(context.Background(), queries.CreatePlayerEmailParams{
			Email: r.Email,
			Pid:   pid.(int64),
		})
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		id, err = result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = tx.Commit()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return c.Render("web/views/partials/profile/email/unverified-email", &fiber.Map{
			"ID":       id,
			"Email":    r.Email,
			"Verified": false,
		}, "")
	}
}
