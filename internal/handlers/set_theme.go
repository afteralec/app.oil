package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/util"
)

func SetTheme(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		theme := c.Params("theme")
		if len(theme) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if theme != constants.ThemeLight && theme != constants.ThemeDark {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		sess.Set("theme", theme)
		if err := sess.Save(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		toggleTheme := constants.ThemeDark
		if theme == constants.ThemeDark {
			toggleTheme = constants.ThemeLight
		} else {
			toggleTheme = constants.ThemeDark
		}
		themeText := "Light"
		if theme == constants.ThemeDark {
			themeText = "Dark"
		}

		toggleThemePath := routes.ThemePath(toggleTheme)

		b := fiber.Map{
			"Theme":           theme,
			"ThemeText":       themeText,
			"ToggleThemePath": toggleThemePath,
		}

		pid, err := util.GetPID(c)
		if err == nil {
			if err := i.Queries.UpdatePlayerSettingsTheme(context.Background(), queries.UpdatePlayerSettingsThemeParams{
				PID:   pid,
				Theme: theme,
			}); err != nil {
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusInternalServerError)
			}
		}

		return c.Render(partials.ThemeToggle, b, layouts.None)
	}
}
