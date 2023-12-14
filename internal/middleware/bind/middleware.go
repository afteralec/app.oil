package bind

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(shared.Bind, fiber.Map{
			"CSRF":                      c.Locals("csrf"),
			"PID":                       c.Locals("pid"),
			"CopyrightYear":             time.Now().Year(),
			"Title":                     "Petrichor",
			"MetaContent":               "Petrichor MUD - a modern take on a classic MUD style of game.",
			"HomeView":                  c.Path() == routes.Home,
			"ProfileView":               c.Path() == routes.Profile || c.Path() == routes.Me,
			"CharactersView":            c.Path() == routes.Characters,
			"PermissionsView":           c.Path() == routes.PlayerPermissions,
			"ShouldShowPermissionsMenu": determineShouldShowPermissionsMenu(c),
		})

		return c.Next()
	}
}

func determineShouldShowPermissionsMenu(c *fiber.Ctx) bool {
	lperms := c.Locals("perms")
	if lperms != nil {
		perms, ok := lperms.(permission.PlayerPermissions)
		if !ok {
			return false
		}

		return perms.HasPermissionInSet(permission.ShowPermissionViewPermissions)
	}
	return false
}
