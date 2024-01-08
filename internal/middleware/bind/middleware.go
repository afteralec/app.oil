package bind

import (
	"net/http"
	"time"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/util"
)

func themeMenu(c *fiber.Ctx) fiber.Map {
	theme := c.Locals("theme")
	toggleTheme := constants.ThemeDark
	if theme == constants.ThemeDark {
		toggleTheme = constants.ThemeLight
	}
	themeText := "Light"
	if theme == constants.ThemeDark {
		themeText = "Dark"
	}

	return fiber.Map{
		"Type":            "Theme",
		"Theme":           theme,
		"ThemeText":       themeText,
		"ToggleThemePath": routes.ThemePath(toggleTheme),
	}
}

func accountMenu(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"Type": "List",
		"Button": fiber.Map{
			"Label": "Account",
		},
		"Sections": []fiber.Map{
			{
				"Items": []fiber.Map{
					{
						"Label":  "Characters",
						"Path":   routes.Characters,
						"Active": c.Path() == routes.Characters,
					},
					{
						"Label":  "Profile",
						"Path":   routes.Profile,
						"Active": c.Path() == routes.Profile,
					},
				},
			},
			{
				"Items": []fiber.Map{
					{
						"Label":  "Logout",
						"Path":   routes.Logout,
						"Action": true,
					},
				},
			},
		},
	}
}

func reviewMenu(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"Type": "List",
		"Button": fiber.Map{
			"Label": "Review",
		},
		"Sections": []fiber.Map{
			{
				"Items": []fiber.Map{
					{
						"Label":  "Character Applications",
						"Path":   routes.CharacterApplications,
						"Active": c.Path() == routes.CharacterApplications,
					},
				},
			},
		},
	}
}

func permissionsMenu(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"Type": "List",
		"Button": fiber.Map{
			"Label": "Permissions",
		},
		"Sections": []fiber.Map{
			{
				"Items": []fiber.Map{
					{
						"Label":  "Player Permissions",
						"Path":   routes.PlayerPermissions,
						"Active": c.Path() == routes.PlayerPermissions,
					},
				},
			},
		},
	}
}

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		b := fiber.Map{
			// TODO: See if the CSRF token should be under a page-only bind
			"CSRF": c.Locals("csrf"),
			"PID":  c.Locals("pid"),
		}

		// TODO: Move this to a view bind
		if c.Method() == http.MethodGet {
			b["Theme"] = c.Locals("theme")
			b["CopyrightYear"] = time.Now().Year()
			b["Title"] = "Petrichor"
			b["MetaContent"] = "Petrichor MUD - a modern take on a classic"
			b["Path"] = c.Path()
			b["Menus"] = menus(c)
		}

		c.Locals(constants.BindName, b)
		return c.Next()
	}
}

func menus(c *fiber.Ctx) []fiber.Map {
	menus := []fiber.Map{
		themeMenu(c),
	}

	_, err := util.GetPID(c)
	if err != nil {
		return menus
	}

	menus = append(menus, accountMenu(c))

	perms, err := util.GetPermissions(c)
	if err != nil {
		return menus
	}

	if perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
		menus = append(menus, reviewMenu(c))
	}

	if perms.Permissions[permissions.PlayerGrantAllPermissionsName] {
		menus = append(menus, permissionsMenu(c))
	}

	return menus
}
