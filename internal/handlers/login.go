package handlers

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
	"petrichormud.com/app/internal/views"
)

// TODO: When you log in or create an account, take the theme from the current session and set it
func Login(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), r.Username)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}

		v, err := password.Verify(r.Password, p.PwHash)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}
		if !v {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}

		pid := p.ID
		err = username.Cache(i.Redis, pid, p.Username)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}

		settings, err := qtx.GetPlayerSettings(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This means a player got created without settings
				c.Status(fiber.StatusInternalServerError)
				c.Append("HX-Retarget", "#login-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusUnauthorized)
				return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}

		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}

		sess.Set("pid", pid)
		theme := sess.Get("theme")
		if theme != nil {
			if err := qtx.UpdatePlayerSettingsTheme(context.Background(), queries.UpdatePlayerSettingsThemeParams{
				PID:   pid,
				Theme: theme.(string),
			}); err != nil {
				c.Append("HX-Retarget", "#login-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusUnauthorized)
				return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
			}
		} else {
			sess.Set("theme", settings.Theme)
		}
		if err = sess.Save(); err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partials.NoticeSectionError, partials.BindLoginErr, layouts.None)
		}

		c.Append("HX-Refresh", "true")
		return nil
	}
}

func LoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid != nil {
			return c.Redirect(routes.Home)
		}

		return c.Render(views.Login, views.Bind(c), layouts.Standalone)
	}
}
