package views

import (
	"os"
	"strings"
	"time"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/util"
)

// TODO: Clean this up to split Menus and Footer stuff only when needed
func Bind(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"CSRF":          c.Locals("csrf"),
		"Theme":         c.Locals("theme"),
		"CopyrightYear": time.Now().Year(),
		"HeadTitle":     "Petrichor",
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

	perms, err := util.GetPermissions(c)
	if err != nil {
		return nav
	}
	if perms.HasPermission(permissions.PlayerReviewCharacterApplicationsName) {
		nav = append(nav, reviewMenu(c))
	}
	// TODO: Clean up this permissions check
	if perms.HasPermission(permissions.PlayerViewAllActorImagesName) {
		nav = append(nav, actorMenu(c))
	}
	if perms.HasPermission(permissions.PlayerViewAllRoomsName) {
		nav = append(nav, roomsMenu(c))
	}
	if perms.HasPermission(permissions.PlayerGrantAllPermissionsName) {
		nav = append(nav, permissionsMenu(c))
	}

	nav = append(nav, accountMenu(c))

	nav = append(nav, playButton())

	return nav
}

func playButton() fiber.Map {
	return fiber.Map{
		"Type": "Play",
		"Path": os.Getenv("PETRICHOR_PLAY_URL"),
	}
}

func themeButton(c *fiber.Ctx) fiber.Map {
	theme := c.Locals("theme")
	toggleTheme := constant.ThemeDark
	if theme == constant.ThemeDark {
		toggleTheme = constant.ThemeLight
	}
	themeText := "Light"
	if theme == constant.ThemeDark {
		themeText = "Dark"
	}

	return fiber.Map{
		"Type":            "Theme",
		"Theme":           theme,
		"ThemeText":       themeText,
		"ToggleThemePath": route.ThemePath(toggleTheme),
	}
}

func helpLink(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"Type":   "Link",
		"Path":   route.Help,
		"Text":   "Help",
		"Active": c.Path() == route.Help,
	}
}

func actorMenu(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"Type": "List",
		"Button": fiber.Map{
			"Label": "actor",
		},
		"Sections": []fiber.Map{
			{
				"Items": []fiber.Map{
					{
						"Label":  "Actor Images",
						"Path":   route.ActorImages,
						"Active": c.Path() == route.ActorImages,
					},
				},
			},
		},
		"Path":   route.Rooms,
		"Text":   "Rooms",
		"Active": strings.Contains(c.Path(), route.Rooms),
	}
}

func roomsMenu(c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"Type": "List",
		"Button": fiber.Map{
			"Label": "Rooms",
		},
		"Sections": []fiber.Map{
			{
				"Items": []fiber.Map{
					{
						"Label":  "Rooms",
						"Path":   route.Rooms,
						"Active": c.Path() == route.Rooms,
					},
				},
			},
		},
		"Path":   route.Rooms,
		"Text":   "Rooms",
		"Active": strings.Contains(c.Path(), route.Rooms),
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
						"Path":   route.Characters,
						"Active": c.Path() == route.Characters,
					},
					{
						"Label":  "Profile",
						"Path":   route.Profile,
						"Active": c.Path() == route.Profile,
					},
				},
			},
			{
				"Items": []fiber.Map{
					{
						"Label":  "Logout",
						"Path":   route.Logout,
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
						"Path":   route.CharacterApplications,
						"Active": c.Path() == route.CharacterApplications,
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
						"Path":   route.PlayerPermissions,
						"Active": c.Path() == route.PlayerPermissions,
					},
				},
			},
		},
	}
}
