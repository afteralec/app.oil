package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/shared"
)

const ProfileRoute = "/profile"

func ProfilePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		emails, err := i.Queries.ListEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
		}

		b := c.Locals("bind").(fiber.Map)
		b["Emails"] = emails
		b["VerifiedEmails"] = email.Verified(emails)
		b["GravatarEmail"] = "othertest@quack.ninja"
		b["GravatarHash"] = email.GravatarHash("after.alec@gmail.com")

		return c.Render("web/views/profile", b)
	}
}
