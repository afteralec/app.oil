package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
	"petrichormud.com/app/internal/views"
)

func RecoverUsernamePage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("views/recover/username", c.Locals(constants.BindName), views.LayoutStandalone)
	}
}

func RecoverUsernameSuccessPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := c.Query("t")
		if len(t) == 0 {
			c.Redirect(routes.Home)
		}

		key := username.RecoverySuccessKey(t)
		address, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Redirect(routes.Home)
		}

		b := c.Locals(constants.BindName).(fiber.Map)
		b["EmailAddress"] = address

		return c.Render("views/recover/username/success", b, views.LayoutStandalone)
	}
}

func RecoverUsername(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Email string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render("views/partials/recover/username/err-invalid", c.Locals(constants.BindName), "")
		}

		e, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render("views/partials/recover/username/err-invalid", c.Locals(constants.BindName), "")
		}

		ve, err := i.Queries.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil {
			if err == sql.ErrNoRows {
				rusid, err := username.CacheRecoverySuccessEmail(i.Redis, e.Address)
				if err != nil {
					c.Status(fiber.StatusUnauthorized)
					c.Append(shared.HeaderHXAcceptable, "true")
					return c.Render("views/partials/recover/username/err-internal", c.Locals(constants.BindName), "")
				}

				path := fmt.Sprintf("%s?t=%s", routes.RecoverUsernameSuccess, rusid)
				c.Append("HX-Redirect", path)
				return nil
			}
			c.Status(fiber.StatusUnauthorized)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render("views/partials/recover/username/err", c.Locals(constants.BindName), "")
		}

		rusid, err := username.Recover(i, ve)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render("views/partials/recover/username/err-internal", c.Locals(constants.BindName), "")
		}

		path := fmt.Sprintf("%s?t=%s", routes.RecoverUsernameSuccess, rusid)
		c.Append("HX-Reswap", "none")
		c.Append("HX-Redirect", path)
		return nil
	}
}
