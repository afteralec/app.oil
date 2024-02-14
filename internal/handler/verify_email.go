package handler

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/username"
	"petrichormud.com/app/internal/views"
)

func VerifyEmailPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layout.Standalone)
		}

		token := c.Query("t")
		key := email.VerificationKey(token)

		exists, err := i.Redis.Exists(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}
		if exists != 1 {
			c.Status(fiber.StatusNotFound)
			b := views.Bind(c)
			b["NotFoundMessage"] = "Sorry, it looks like this link has expired."
			b["NotFoundButtonLink"] = routes.Profile
			b["NotFoundButtonText"] = "Return to Profile"
			return c.Render(views.NotFound, b, layout.Standalone)
		}

		eid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			if err == redis.Nil {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}
		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}

		if e.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}
		if err == nil && ve.Verified {
			c.Status(fiber.StatusConflict)
			b := views.Bind(c)
			b["ErrMessageConflict"] = "That email has already been verified."
			return c.Render(views.Conflict, b, layout.Standalone)
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		un, err := username.Get(i, pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layout.Standalone)
		}

		b := views.Bind(c)
		b["VerifyToken"] = c.Query("t")
		b["Address"] = e.Address
		b["Username"] = un
		return c.Render(views.VerifyEmail, b, layout.Standalone)
	}
}

func VerifyEmail(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Refresh", "true")
			return nil
		}

		token := c.Query("t")
		if len(token) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		key := email.VerificationKey(token)
		eid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			if err == redis.Nil {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		e, err := i.Queries.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if e.Verified {
			c.Status(fiber.StatusConflict)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = i.Queries.MarkEmailVerified(context.Background(), id); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return c.Render(partial.VerifyEmailSuccess, &fiber.Map{}, "")
	}
}
