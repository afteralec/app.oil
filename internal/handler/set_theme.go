package handler

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/util"
)

func SetTheme(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		theme := c.Params("theme")
		if len(theme) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if theme != constant.ThemeLight && theme != constant.ThemeDark {
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

		toggleTheme := constant.ThemeDark
		if theme == constant.ThemeDark {
			toggleTheme = constant.ThemeLight
		} else {
			toggleTheme = constant.ThemeDark
		}
		themeText := "Light"
		if theme == constant.ThemeDark {
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
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusInternalServerError)
			}
		}

		return c.Render(partial.ThemeToggle, b, layout.None)
	}
}
