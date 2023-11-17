package handlers

import (
	"net/mail"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
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

		_, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		return nil
	}
}
