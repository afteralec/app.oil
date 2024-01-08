package handlers

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/shared"
)

func SetTheme(i *shared.Interfaces) fiber.Handler {
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

		b := fiber.Map{
			"Theme":       theme,
			"ThemeText":   themeText,
			"ToggleTheme": toggleTheme,
		}
		return c.Render(partials.ThemeToggle, b, layouts.None)
	}
}
