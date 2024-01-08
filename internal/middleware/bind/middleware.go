package bind

import (
	"html/template"
	"time"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/routes"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		theme := c.Locals("theme")
		toggleTheme := "dark"
		if theme == "dark" {
			toggleTheme = "light"
		}
		themeText := "Light"
		if theme == constants.ThemeDark {
			themeText = "Dark"
		}

		b := fiber.Map{
			"Theme":           theme,
			"ThemeText":       themeText,
			"ToggleTheme":     toggleTheme,
			"CSRF":            c.Locals("csrf"),
			"PID":             c.Locals("pid"),
			"CopyrightYear":   time.Now().Year(),
			"Title":           "Petrichor",
			"MetaContent":     "Petrichor MUD - a modern take on a classic MUD style of game.",
			"ShouldShowMenus": determineShouldShowMenus(c),
			"HomeView":        c.Path() == routes.Home,
			"ProfileView":     c.Path() == routes.Profile || c.Path() == routes.Me,
			"CharactersView":  c.Path() == routes.Characters,
			"PermissionsView": c.Path() == routes.PlayerPermissions,
			"IconWarn":        template.URL("ant-design:exclamation-circle-outlined"),
		}
		c.Locals(constants.BindName, b)

		return c.Next()
	}
}

type shouldShowMenus struct {
	Review                      bool
	ReviewCharacterApplications bool
	Permissions                 bool
}

func determineShouldShowMenus(c *fiber.Ctx) shouldShowMenus {
	lperms := c.Locals("perms")
	if lperms != nil {
		perms, ok := lperms.(permissions.PlayerGranted)
		if !ok {
			return shouldShowMenus{
				Review:                      false,
				ReviewCharacterApplications: false,
				Permissions:                 false,
			}
		}

		return shouldShowMenus{
			Review:                      perms.Permissions[permissions.PlayerReviewCharacterApplicationsName],
			ReviewCharacterApplications: perms.Permissions[permissions.PlayerReviewCharacterApplicationsName],
			Permissions:                 perms.Permissions[permissions.PlayerGrantAllPermissionsName],
		}
	}

	return shouldShowMenus{
		Review:                      false,
		ReviewCharacterApplications: false,
		Permissions:                 false,
	}
}
