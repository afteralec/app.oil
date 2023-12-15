package handlers

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/routes"
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

		perms, ok := lperms.(permission.PlayerIssuedPermissions)
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
	type playerPermissionDetail struct {
		Permission string
		About      string
		Link       string
		Issued     bool
	}
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

		iperms, ok := lperms.(permission.PlayerIssuedPermissions)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		if !iperms.HasPermissionInSet(permission.ShowPermissionViewPermissions) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		u := c.Params("username")
		if len(u) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), u)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("User not found")
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		pperms, err := qtx.ListPlayerPermissions(context.Background(), p.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		perms := permission.MakePlayerPermissions(p.ID, pperms)

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		allPermDetails := []playerPermissionDetail{}
		for _, pd := range permission.AllPlayerPermissionDetails {
			perm, about := pd[0], pd[1]
			allPermDetails = append(allPermDetails, playerPermissionDetail{
				Permission: perm,
				About:      about,
				Link:       routes.PlayerPermissionsPath(strconv.FormatInt(p.ID, 10)),
				Issued:     perms.Permissions[perm],
			})
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Username"] = u
		b["Permissions"] = allPermDetails
		b["PermissionsPath"] = routes.PlayerPermissionsPath(strconv.FormatInt(p.ID, 10))
		return c.Render("views/player_permissions_detail", b)
	}
}

func UpdatePlayerPermission(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println(string(c.Body()))

		ipid := c.Locals("pid")

		if ipid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		iperms, ok := lperms.(permission.PlayerIssuedPermissions)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		// TODO: This should be the permission to assign this permission
		if !iperms.HasPermissionInSet(permission.ShowPermissionViewPermissions) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		ppid := c.Params("id")
		if len(ppid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		pid, err := strconv.ParseInt(ppid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayer(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		pperms, err := qtx.ListPlayerPermissions(context.Background(), p.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_ = permission.MakePlayerPermissions(p.ID, pperms)

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		return c.Render("views/404", b)
	}
}
