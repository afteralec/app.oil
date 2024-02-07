package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/views"
)

// TODO: Add the Avatar section back into the profile with just Gravatar
// TODO: Add a section for changing your Username and Password
func ProfilePage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layouts.Standalone)
		}

		emails, err := i.Queries.ListEmails(context.Background(), pid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		b := views.Bind(c)
		b["Emails"] = emails
		b["VerifiedEmails"] = email.Verified(emails)
		b["GravatarEmail"] = "othertest@quack.ninja"
		b["GravatarHash"] = email.GravatarHash("after.alec@gmail.com")
		b["ChangePasswordPath"] = routes.PlayerPasswordPath(pid)
		return c.Render(views.Profile, b, layouts.Main)
	}
}
