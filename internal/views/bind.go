package views

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/util"
)

// TODO: Clean this up to split Menus and Footer stuff only when needed
func Bind(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"CSRF":          c.Locals("csrf"),
		"Theme":         c.Locals("theme"),
		"CopyrightYear": time.Now().Year(),
		"Title":         "Petrichor",
		"MetaContent":   "Petrichor MUD - a modern take on a classic",
		"Path":          c.Path(),
		"Nav":           nav(c),
	}
}

func nav(c *fiber.Ctx) []fiber.Map {
	nav := []fiber.Map{
		themeButton(c),
		helpLink(c),
	}

	_, err := util.GetPID(c)
	if err != nil {
		nav = append(nav, fiber.Map{"Type": "Login"})
		nav = append(nav, fiber.Map{"Type": "Register"})
		return nav
	}

	nav = append(nav, accountMenu(c))

	perms, err := util.GetPermissions(c)
	if err != nil {
		return nav
	}

	if perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
		nav = append(nav, reviewMenu(c))
	}

	if perms.Permissions[permissions.PlayerGrantAllPermissionsName] {
		nav = append(nav, permissionsMenu(c))
	}

	return nav
}

func themeButton(c *fiber.Ctx) fiber.Map {
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

// TODO: Re-do this as a dropdown menu with search and various links
func helpLink(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"Type":   "Link",
		"Path":   routes.Help,
		"Text":   "Help",
		"Active": c.Path() == routes.Help,
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
