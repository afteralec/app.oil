package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

func VerifyEmailPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
		}

		token := c.Query("t")
		key := email.VerificationKey(token)

		exists, err := i.Redis.Exists(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
		}
		if exists != 1 {
			c.Status(fiber.StatusNotFound)
			b := c.Locals(bind.Name).(fiber.Map)
			b["NotFoundMessage"] = "Sorry, it looks like this link has expired."
			b["NotFoundButtonLink"] = routes.Profile
			b["NotFoundButtonText"] = "Return to Profile"
			return c.Render("views/404", b, "views/layouts/standalone")
		}

		eid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			if err == redis.Nil {
				c.Status(fiber.StatusNotFound)
				return c.Render("views/404", c.Locals(bind.Name), "views/layouts/standalone")
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
		}
		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
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
				return c.Render("views/404", c.Locals(bind.Name), "views/layouts/standalone")
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
		}

		if e.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
		}
		if err == nil && ve.Verified {
			c.Status(fiber.StatusConflict)
			return c.Render("views/verify-email-409", c.Locals(bind.Name), "views/layouts/standalone")
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		un, err := username.Get(i, pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(bind.Name), "views/layouts/standalone")
		}

		b := c.Locals(bind.Name).(fiber.Map)
		b["VerifyToken"] = c.Query("t")
		b["Address"] = e.Address
		b["Username"] = un
		return c.Render("views/verify-email", b, "views/layouts/standalone")
	}
}

func VerifyEmail(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(shared.HeaderHXAcceptable, "true")
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

		_, err = i.Queries.MarkEmailVerified(context.Background(), id)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return c.Render("views/partials/verify/success", &fiber.Map{}, "")
	}
}
