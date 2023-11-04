package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

const VerifyRoute = "/verify"

func VerifyPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("t")
		exists, err := i.Redis.Exists(context.Background(), token).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
		}
		if exists != 1 {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/401", c.Locals("bind"), "web/views/layouts/standalone")
		}

		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		eid, err := i.Redis.Get(context.Background(), token).Result()
		if err != nil {
			// TODO: Differentiate between "not found" error and connection error
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
		}
		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
		}
		e, err := i.Queries.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				// TODO: Make this a 404 page
				return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
		}

		_, err = i.Queries.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil {
			if err != sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
			}
		}
		if err == nil {
			c.Status(fiber.StatusConflict)
			// TODO: Build an "already verified" page for this
			return c.Render("web/views/401", c.Locals("bind"), "web/views/layouts/standalone")
		}

		un, err := username.Get(i.Redis, pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals("bind"), "web/views/layouts/standalone")
		}

		b := c.Locals("bind").(fiber.Map)
		b["VerifyToken"] = c.Query("t")
		b["Address"] = e.Address
		b["Username"] = un

		return c.Render("web/views/verify", b, "web/views/layouts/standalone")
	}
}

func Verify(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			// TODO: This should redirect them back to the login page for this token
			return nil
		}

		key := c.Query("t")
		if len(key) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		eid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: This same block of code is reused across a couple of email operations; peel it out?
		e, err := i.Queries.GetEmail(context.Background(), id)
		if err != nil {
			// TODO: Distinguish between "not found" and other errors
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil {
			if err != sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}
		if err == nil {
			c.Status(fiber.StatusConflict)
			return nil
		}

		_, err = i.Queries.MarkEmailVerified(context.Background(), id)
		if err != nil {
			// TODO: Distinguish between a "not found" error and a connection error
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return c.Render("web/views/partials/verify/success", &fiber.Map{}, "")
	}
}
