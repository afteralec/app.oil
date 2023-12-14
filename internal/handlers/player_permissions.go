package handlers

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/shared"
)

func PlayerPermissionsPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		perms, ok := lperms.(permission.PlayerPermissions)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		if !perms.HasPermissionInSet(permission.ShowPermissionViewPermissions) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		return c.Render("views/player_permissions", b)
	}
}

func PlayerPermissionsDetailPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		perms, ok := lperms.(permission.PlayerPermissions)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		if !perms.HasPermissionInSet(permission.ShowPermissionViewPermissions) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		u := c.Params("username")
		if len(u) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Username"] = u
		return c.Render("views/player_permissions_detail", b)
	}
}
