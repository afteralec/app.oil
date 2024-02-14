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
	"petrichormud.com/app/internal/player/username"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/view"
)

func UsernameReserved(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Username string `form:"username"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			return err
		}

		p, err := i.Queries.GetPlayerByUsername(context.Background(), in.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
				return c.Render(partial.PlayerFree, fiber.Map{
					"CSRF": c.Locals("csrf"),
				}, layout.CSRF)
			}
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.PlayerReservedErr, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		}

		if in.Username == p.Username {
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partial.PlayerReserved, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		} else {
			c.Append("HX-Trigger-After-Swap", "ptrcr:username-reserved")
			return c.Render(partial.PlayerFree, fiber.Map{
				"CSRF": c.Locals("csrf"),
			}, layout.CSRF)
		}
	}
}

func RecoverUsernamePage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(view.RecoverUsername, view.Bind(c), layout.Standalone)
	}
}

func RecoverUsernameSuccessPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := c.Query("t")
		if len(t) == 0 {
			c.Redirect(route.Home)
		}

		key := username.RecoverySuccessKey(t)
		address, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Redirect(route.Home)
		}

		b := view.Bind(c)
		b["EmailAddress"] = address
		return c.Render(view.RecoverUsernameSuccess, b, layout.Standalone)
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

				path := fmt.Sprintf("%s?t=%s", route.RecoverUsernameSuccess, rusid)
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

		path := fmt.Sprintf("%s?t=%s", route.RecoverUsernameSuccess, rusid)
		c.Append("HX-Reswap", "none")
		c.Append("HX-Redirect", path)
		return nil
	}
}
