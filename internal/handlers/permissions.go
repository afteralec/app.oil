package handlers

import (
	"slices"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/shared"
)

const PermissionsRoute = "/permissions"

func PermissionsPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		perms, err := permissions.List(i, pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
		}
		if !slices.Contains(perms, permissions.ViewPermissions) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/401", c.Locals("bind"), "web/views/layouts/standalone")
		}

		return c.Render("web/views/permissions", c.Locals("bind"), "web/views/layouts/main")
	}
}
