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

		perms, ok := lperms.(permission.PlayerIssued)
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

		iperms, ok := lperms.(permission.PlayerIssued)
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

		perms := permission.MakePlayerIssued(p.ID, pperms)
		allPerms := []fiber.Map{}
		for _, perm := range permission.AllPlayer {
			allPerms = append(allPerms, fiber.Map{
				"Name":   perm.Name,
				"Tag":    perm.Tag,
				"Title":  perm.Title,
				"About":  perm.About,
				"Link":   routes.PlayerPermissionsTogglePath(strconv.FormatInt(p.ID, 10), perm.Tag),
				"Issued": perms.Permissions[perm.Name],
			})
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Username"] = u
		b["Permissions"] = allPerms
		return c.Render("views/player_permissions_detail", b)
	}
}

func UpdatePlayerPermission(i *shared.Interfaces) fiber.Handler {
	type input struct {
		AssignAll                   string `form:"assign-all"`
		ReviewCharacterApplications string `form:"review-character-applications"`
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

		iperms, ok := lperms.(permission.PlayerIssued)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
		}

		if !iperms.Permissions[permission.PlayerAssignAllPermissionsName] {
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

		// TODO: What we're going to do here is build a diff of the incoming set of perms
		// and what the player currently has.
		// For every permission on the input that hasn't already been issued, we build
		// an issued history item and then issue the permission.
		// For every permission that's been issued that hasn't already been issued,
		// we build a revoked history item and revoke the permission.
		_ = permission.MakePlayerIssued(pid, pperms)

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return nil
	}
}

func TogglePlayerPermission(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Issued bool `form:"issued"`
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
		iperms, ok := lperms.(permission.PlayerIssued)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(shared.Bind), "views/layouts/standalone")
		}
		if !iperms.Permissions[permission.PlayerAssignAllPermissionsName] {
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

		log.Println(r)

		_ = permission.MakePlayerIssued(pid, pperms)

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return nil
	}
}
