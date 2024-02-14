package handler

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

// TODO: Add the Avatar section back into the profile with just Gravatar
func ProfilePage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		emails, err := i.Queries.ListEmails(context.Background(), pid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		b := view.Bind(c)
		b["Emails"] = emails
		b["VerifiedEmails"] = email.Verified(emails)
		b["GravatarEmail"] = "othertest@quack.ninja"
		b["GravatarHash"] = email.GravatarHash("after.alec@gmail.com")
		b["ChangePasswordPath"] = route.PlayerPasswordPath(pid)
		return c.Render(view.Profile, b, layout.Main)
	}
}
