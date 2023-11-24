package handlers

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

const ResetPasswordRoute = "/reset/password"

func ResetPasswordPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover/password", c.Locals("bind"), "web/views/layouts/standalone")
	}
}

func ResetPassword(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Username        string `form:"username"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirm"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		vu := username.Validate(r.Username)
		if !vu {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if r.Password != r.ConfirmPassword {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		vp := password.Validate(r.Password)
		if !vp {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		// TODO: Check the token in the URL here

		// TODO: Transaction this up
		_, err := i.Queries.GetPlayerByUsername(context.Background(), r.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Reset the password given the input here
		return nil
	}
}
