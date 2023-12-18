package bind

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/permissions"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		b := fiber.Map{
			"CSRF": c.Locals("csrf"),
			"PID":  c.Locals("pid"),
			// TODO: Starting in 2024, have this be 2023 - current year
			"CopyrightYear":   time.Now().Year(),
			"Title":           "Petrichor",
			"MetaContent":     "Petrichor MUD - a modern take on a classic MUD style of game.",
			"ShouldShowMenus": determineShouldShowMenus(c),
		}
		b = bind.CurrentView(b, c)
		c.Locals(bind.Name, b)

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
