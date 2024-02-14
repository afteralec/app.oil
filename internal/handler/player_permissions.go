package handler

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/views"
)

func PlayerPermissionsPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layout.Standalone)
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		perms, ok := lperms.(permissions.PlayerGranted)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}

		if !perms.HasPermissionInSet(permissions.ShowPermissionViewPermissions) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		return c.Render(views.PlayerPermissions, views.Bind(c), layout.Main)
	}
}

func PlayerPermissionsDetailPage(i *interfaces.Shared) fiber.Handler {
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
			return c.Render(views.Login, views.Bind(c), layout.Standalone)
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		iperms, ok := lperms.(permissions.PlayerGranted)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}

		if !iperms.HasPermissionInSet(permissions.ShowPermissionViewPermissions) {
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

		perms := permissions.MakePlayerGranted(p.ID, pperms)
		allPerms := []fiber.Map{}
		for _, perm := range permissions.AllPlayer {
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

		b := views.Bind(c)
		b["Username"] = u
		b["Permissions"] = allPerms
		return c.Render(views.PlayerPermissionsDetail, b)
	}
}

func TogglePlayerPermission(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Grant bool `form:"issued"`
	}
	return func(c *fiber.Ctx) error {
		ipid := c.Locals("pid")
		if ipid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layout.Standalone)
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		iperms, ok := lperms.(permissions.PlayerGranted)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}
		if !iperms.Permissions[permissions.PlayerGrantAllPermissionsName] {
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
		_, ok = permissions.AllPlayerByTag[ptag]
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

		perms := permissions.MakePlayerGranted(pid, pperms)
		perm := permissions.AllPlayerByTag[ptag]
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
