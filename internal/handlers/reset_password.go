package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/headers"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/username"
	"petrichormud.com/app/internal/views"
)

func ResetPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tid := c.Query("t")
		if len(tid) == 0 {
			return c.Redirect("/")
		}

		b := views.Bind(c)
		b["ResetPasswordToken"] = tid

		return c.Render(views.ResetPassword, b, layouts.Standalone)
	}
}

func ResetPasswordSuccessPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.ResetPasswordSuccess, views.Bind(c), layouts.Standalone)
	}
}

func ResetPassword(i *interfaces.Shared) fiber.Handler {
	type request struct {
		Username        string `form:"username"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirmPassword"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		vu := username.IsValid(r.Username)
		if !vu {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		if r.Password != r.ConfirmPassword {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		vp := password.IsValid(r.Password)
		if !vp {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		tid := c.Query("t")
		if len(tid) == 0 {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		key := password.RecoveryKey(tid)
		rpid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		pid, err := strconv.ParseInt(rpid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		p, err := i.Queries.GetPlayer(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusUnauthorized)
				c.Append(headers.HXAcceptable, "true")
				return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
			}
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		if p.Username != r.Username {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		pwHash, err := password.Hash(r.Password)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		_, err = i.Queries.UpdatePlayerPassword(context.Background(), queries.UpdatePlayerPasswordParams{
			ID:     pid,
			PwHash: pwHash,
		})
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindResetPasswordErr, layouts.None)
		}

		c.Append("HX-Redirect", routes.ResetPasswordSuccess)
		return nil
	}
}
