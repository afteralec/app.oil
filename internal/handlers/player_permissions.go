package handlers

import (
	"context"
	"log"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/shared"
)

func PermissionsPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind))
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		ps, err := qtx.ListPlayerPermissions(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind))
		}

		perms, err := permission.MakePlayerPermissions(ps)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !perms.HasPermissionInSet(permission.ShowPermissionViewPermissions) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind))
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		return c.Render("views/player_permissions", b)
	}
}
