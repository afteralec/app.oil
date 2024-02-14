package handler

import (
	"context"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/username"
)

func Register(i *interfaces.Shared) fiber.Handler {
	type Player struct {
		Username        string `form:"username"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirmPassword"`
	}

	return func(c *fiber.Ctx) error {
		p := new(Player)

		if err := c.BodyParser(p); err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
		}

		u := username.Sanitize(p.Username)

		if !username.IsValid(u) {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInvalidUsername, layout.None)
		}

		if !password.IsValid(p.Password) {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInvalidPassword, layout.None)
		}

		if p.Password != p.ConfirmPassword {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInvalidConfirmPassword, layout.None)
		}

		pwHash, err := password.Hash(p.Password)
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)

		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)
		result, err := qtx.CreatePlayer(
			context.Background(),
			query.CreatePlayerParams{
				Username: u,
				PwHash:   pwHash,
			},
		)
		if err != nil {
			me, ok := err.(*mysql.MySQLError)
			if !ok {
				c.Append("HX-Retarget", "#register-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusInternalServerError)
				return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
			}
			if me.Number == mysqlerr.ER_DUP_ENTRY {
				c.Append("HX-Retarget", "#register-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusConflict)
				return c.Render(partial.NoticeSectionError, partial.BindRegisterErrConflict, layout.None)
			}
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
		}

		pid, err := result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
		}

		if err := qtx.CreatePlayerSettings(context.Background(), query.CreatePlayerSettingsParams{
			PID:   pid,
			Theme: constant.ThemeDefault,
		}); err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
		}

		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
		}

		theme := sess.Get("theme")
		if theme != nil {
			if err := qtx.UpdatePlayerSettingsTheme(context.Background(), query.UpdatePlayerSettingsThemeParams{
				PID:   pid,
				Theme: theme.(string),
			}); err != nil {
				c.Append("HX-Retarget", "#register-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusInternalServerError)
				return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
			}
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
		}

		sess.Set("pid", pid)
		if err = sess.Save(); err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindRegisterErrInternal, layout.None)
		}

		username.Cache(i.Redis, pid, p.Username)

		c.Append("HX-Refresh", "true")
		c.Status(fiber.StatusCreated)
		return nil
	}
}
