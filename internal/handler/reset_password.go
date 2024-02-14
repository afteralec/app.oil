package handler

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/player/password"
	"petrichormud.com/app/internal/player/username"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/view"
)

func ResetPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tid := c.Query("t")
		if len(tid) == 0 {
			return c.Redirect("/")
		}

		b := view.Bind(c)
		b["ResetPasswordToken"] = tid

		return c.Render(view.ResetPassword, b, layout.Standalone)
	}
}

func ResetPasswordSuccessPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(view.ResetPasswordSuccess, view.Bind(c), layout.Standalone)
	}
}

func ResetPassword(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Username        string `form:"username"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirmPassword"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		vu := username.IsValid(in.Username)
		if !vu {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		if in.Password != in.ConfirmPassword {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		vp := password.IsValid(in.Password)
		if !vp {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		tid := c.Query("t")
		if len(tid) == 0 {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		key := password.RecoveryKey(tid)
		rpid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		pid, err := strconv.ParseInt(rpid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		p, err := i.Queries.GetPlayer(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusUnauthorized)
				c.Append(header.HXAcceptable, "true")
				return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
			}
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		if p.Username != in.Username {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		pwHash, err := password.Hash(in.Password)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		_, err = i.Queries.UpdatePlayerPassword(context.Background(), query.UpdatePlayerPasswordParams{
			ID:     pid,
			PwHash: pwHash,
		})
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindResetPasswordErr, layout.None)
		}

		c.Append("HX-Redirect", route.ResetPasswordSuccess)
		return nil
	}
}
