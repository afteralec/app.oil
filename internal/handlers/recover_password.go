package handlers

import (
	"context"
	"database/sql"
	"net/mail"
	"slices"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

const RecoverPasswordRoute = "/recover/password"

func RecoverPasswordPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover/password", c.Locals("bind"), "web/views/layouts/standalone")
	}
}

func RecoverPassword(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Username string `form:"username"`
		Email    string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		v := username.Validate(r.Username)
		if !v {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		_, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		// TODO: Transaction this up
		p, err := i.Queries.GetPlayerByUsername(context.Background(), r.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		emails, err := i.Queries.ListEmails(context.Background(), p.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if len(emails) == 0 {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		emailAddresses := []string{}
		for i := 0; i < len(emails); i++ {
			email := emails[i]
			emailAddresses = append(emailAddresses, email.Address)
		}

		if !slices.Contains(emailAddresses, r.Email) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		err = password.SetupRecovery(i.Redis, p.ID, r.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		return nil
	}
}
