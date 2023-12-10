package handlers

import (
	"context"
	"database/sql"
	"net/mail"
	"slices"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

func RecoverPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover/password", c.Locals("bind"), "web/views/layouts/standalone")
	}
}

func RecoverPasswordSuccessPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover/password/success", c.Locals("bind"), "web/views/layouts/standalone")
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

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), r.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		emails, err := qtx.ListVerifiedEmails(context.Background(), p.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if len(emails) == 0 {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
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

		err = password.SetupRecovery(i, p.ID, r.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", routes.RecoverPasswordSuccess)
		return nil
	}
}
