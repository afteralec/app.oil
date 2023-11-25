package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

const (
	RecoverUsernameRoute        = "/recover/username"
	RecoverUsernameSuccessRoute = "/recover/username/success"
)

func RecoverUsernamePage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/recover/username", c.Locals("bind"), "web/views/layouts/standalone")
	}
}

func RecoverUsernameSuccessPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := c.Query("t")
		if len(t) == 0 {
			c.Redirect(HomeRoute)
		}

		key := username.RecoverySuccessKey(t)
		ceid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Redirect(HomeRoute)
		}
		eid, err := strconv.ParseInt(ceid, 10, 64)
		if err != nil {
			c.Redirect(HomeRoute)
		}
		email, err := i.Queries.GetEmail(context.Background(), eid)
		if err != nil {
			c.Redirect(HomeRoute)
		}

		b := c.Locals("bind").(fiber.Map)
		b["EmailAddress"] = email.Address

		return c.Render("web/views/recover/username/success", b, "web/views/layouts/standalone")
	}
}

func RecoverUsername(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Email string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		e, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		ve, err := i.Queries.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rusid, err := username.Recover(i, ve)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		path := fmt.Sprintf("%s?t=%s", RecoverUsernameSuccessRoute, rusid)
		c.Append("HX-Redirect", path)
		return nil
	}
}
