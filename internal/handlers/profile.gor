package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/shared"
)

// TODO: Add the Avatar section back into the profile with just Gravatar
// TODO: Add a section for changing your Username and Password
func ProfilePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals(shared.Bind), "web/views/layouts/standalone")
		}

		emails, err := i.Queries.ListEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals(shared.Bind), "web/views/layouts/standalone")
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Emails"] = emails
		b["VerifiedEmails"] = email.Verified(emails)
		b["GravatarEmail"] = "othertest@quack.ninja"
		b["GravatarHash"] = email.GravatarHash("after.alec@gmail.com")

		return c.Render("web/views/profile", b)
	}
}
