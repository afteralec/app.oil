package handlers

import (
	"context"
	"database/sql"
	"net/mail"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

const (
	RecoverUsernameRoute        = "/recover/username"
	RecoverUsernameSuccessRoute = "/recover/username/success"
)

func RecoverUsernamePage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover/username", c.Locals("bind"), "web/views/layouts/standalone")
	}
}

func RecoverUsernameSuccessPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover/username/success", c.Locals("bind"), "web/views/layouts/standalone")
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
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = username.Recover(i, ve)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", "/recover/username/success")
		return nil
	}
}
