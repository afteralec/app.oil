package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/queries"
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

		perms, ok := lperms.(permission.PlayerGranted)
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
		Granted    bool
		Disabled   bool
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

		iperms, ok := lperms.(permission.PlayerGranted)
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

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		perms := permission.MakePlayerGranted(p.ID, pperms)
		allPerms := []fiber.Map{}
		for _, perm := range permission.AllPlayer {
			granted := perms.Permissions[perm.Name]
			disabled := true
			if granted {
				disabled = !iperms.CanRevokePermission(perm.Name)
			} else {
				disabled = !iperms.CanGrantPermission(perm.Name)
			}
			pm := fiber.Map{
				"Name":     perm.Name,
				"Tag":      perm.Tag,
				"Title":    perm.Title,
				"About":    perm.About,
				"Link":     routes.PlayerPermissionsTogglePath(strconv.FormatInt(p.ID, 10), perm.Tag),
				"Granted":  granted,
				"Disabled": disabled,
			}
			allPerms = append(allPerms, pm)
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Username"] = u
		b["Permissions"] = allPerms
		return c.Render("views/player_permissions_detail", b)
	}
}

func TogglePlayerPermission(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Grant bool `form:"issued"`
	}
	return func(c *fiber.Ctx) error {
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
		iperms, ok := lperms.(permission.PlayerGranted)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
		}
		if !iperms.Permissions[permission.PlayerGrantAllPermissionsName] {
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

		ptag := c.Params("tag")
		if len(ptag) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		_, ok = permission.AllPlayerByTag[ptag]
		if !ok {
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

		pperms, err := qtx.ListPlayerPermissions(context.Background(), pid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err = c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		perms := permission.MakePlayerGranted(pid, pperms)
		perm := permission.AllPlayerByTag[ptag]
		_, granted := perms.Permissions[perm.Name]

		if r.Grant && granted {
			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			c.Status(fiber.StatusConflict)
			return nil
		}

		if r.Grant && !granted {
			if !perms.CanGrantPermission(perm.Name) {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			params := queries.CreatePlayerPermissionIssuedChangeHistoryParams{
				IPID:       ipid.(int64),
				PID:        pid,
				Permission: perm.Name,
			}
			if err := qtx.CreatePlayerPermissionIssuedChangeHistory(context.Background(), params); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			_, err := qtx.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
				IPID:       ipid.(int64),
				PID:        pid,
				Permission: perm.Name,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			return nil
		}

		if !r.Grant && granted {
			if !perms.CanRevokePermission(perm.Name) {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			params := queries.CreatePlayerPermissionRevokedChangeHistoryParams{
				IPID:       ipid.(int64),
				PID:        pid,
				Permission: perm.Name,
			}
			if err := qtx.CreatePlayerPermissionRevokedChangeHistory(context.Background(), params); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if err := qtx.DeletePlayerPermission(context.Background(), queries.DeletePlayerPermissionParams{
				PID:        pid,
				Permission: perm.Name,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			return nil
		}

		if !r.Grant && !granted {
			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
	}
}
