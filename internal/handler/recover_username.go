package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/username"
	"petrichormud.com/app/internal/views"
)

func RecoverUsernamePage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.RecoverUsername, views.Bind(c), layout.Standalone)
	}
}

func RecoverUsernameSuccessPage(i *interfaces.Shared) fiber.Handler {
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

		b := views.Bind(c)
		b["EmailAddress"] = address
		return c.Render(views.RecoverUsernameSuccess, b, layout.Standalone)
	}
}

func RecoverUsername(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Email string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverUsernameErrInvalid, layout.None)
		}

		e, err := mail.ParseAddress(in.Email)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverUsernameErrInvalid, layout.None)
		}

		ve, err := i.Queries.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil {
			if err == sql.ErrNoRows {
				rusid, err := username.CacheRecoverySuccessEmail(i.Redis, e.Address)
				if err != nil {
					c.Status(fiber.StatusUnauthorized)
					c.Append(header.HXAcceptable, "true")
					return c.Render(partial.NoticeSectionError, partial.BindRecoverUsernameErrInternal, layout.None)
				}

				path := fmt.Sprintf("%s?t=%s", routes.RecoverUsernameSuccess, rusid)
				c.Append("HX-Redirect", path)
				return nil
			}
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverUsernameErrInternal, layout.None)
		}

		rusid, err := username.Recover(i, ve)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverUsernameErrInternal, layout.None)
		}

		path := fmt.Sprintf("%s?t=%s", routes.RecoverUsernameSuccess, rusid)
		c.Append("HX-Reswap", "none")
		c.Append("HX-Redirect", path)
		return nil
	}
}
