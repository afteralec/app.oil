package handlers

import (
	"context"
	"database/sql"
	"net/mail"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

const RecoverUsernameRoute = "/recover/username"

func RecoverUsernamePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover/username", c.Locals("bind"), "web/views/layouts/standalone")
	}
}

func RecoverUsername(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Email string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		e, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		ve, err := i.Queries.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return nil
		}

		err = username.Recover(i, ve)
		if err != nil {
			return nil
		}

		return nil
	}
}
