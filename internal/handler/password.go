package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"slices"
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/player/password"
	"petrichormud.com/app/internal/player/username"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

// TODO: This and the markup allow entering your current password
// This should be a failure case
func ChangePassword(i *service.Interfaces) fiber.Handler {
	type input struct {
		Current         string `form:"current"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirm"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Something's gone terribly wrong.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}

		if in.Password != in.ConfirmPassword {
			c.Status(fiber.StatusBadRequest)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"The new password and password confirmation doesn't match.",
					"Please try again.",
				},
				"NoticeIcon": true,
			}
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Your session has expired.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Something's gone terribly wrong.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayer(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This is a catastrophic failure; a Player object doesn't exist for a logged-in player
				c.Status(fiber.StatusInternalServerError)
				b := fiber.Map{
					"NoticeSectionID": "profile-password-notice",
					"SectionClass":    "pt-4",
					"NoticeText": []string{
						"Something's gone terribly wrong.",
					},
					"NoticeIcon":    true,
					"RefreshButton": true,
				}
				return c.Render(partial.NoticeSectionError, b, layout.None)
			}
			c.Status(fiber.StatusInternalServerError)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Something's gone terribly wrong.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}

		ok, err := password.Verify(in.Current, p.PwHash)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Something's gone terribly wrong.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}
		if !ok {
			c.Status(fiber.StatusUnauthorized)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"The current password you entered isn't correct.",
				},
				"NoticeIcon": true,
			}
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}

		if !password.IsValid(in.Password) {
			c.Status(fiber.StatusBadRequest)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"What you entered isn't a valid password.",
					"Please follow the prompts and try again.",
				},
				"NoticeIcon": true,
			}
			return c.Render(partial.NoticeSectionError, b, layout.None)
		}

		b := fiber.Map{
			"NoticeSectionID": "profile-password-notice",
			"SectionClass":    "pt-4",
			"NoticeText": []string{
				"Success!",
				"Your password has been changed.",
			},
			"NoticeIcon": true,
		}

		return c.Render(partial.NoticeSectionSuccess, b, layout.None)
	}
}

func RecoverPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(view.RecoverPassword, view.Bind(c), layout.Standalone)
	}
}

func RecoverPasswordSuccessPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tid := c.Query("t")
		key := password.RecoverySuccessKey(tid)
		email, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			if err == redis.Nil {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		b := view.Bind(c)
		b["EmailAddress"] = email
		return c.Render(view.RecoverPasswordSuccess, b, layout.Standalone)
	}
}

func RecoverPassword(i *service.Interfaces) fiber.Handler {
	type input struct {
		Username string `form:"username"`
		Email    string `form:"email"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
		}

		v := username.IsValid(in.Username)
		if !v {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInvalidUsername, layout.None)
		}

		_, err := mail.ParseAddress(in.Email)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInvalidEmail, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), in.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				id, err := password.SetupRecoverySuccess(i, in.Email)
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					c.Append(header.HXAcceptable, "true")
					return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
				}
				var sb strings.Builder
				fmt.Fprintf(&sb, "%s?t=%s", route.RecoverPasswordSuccess, id)
				path := sb.String()
				c.Append("HX-Redirect", path)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
		}

		emails, err := qtx.ListVerifiedEmails(context.Background(), p.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
		}
		if len(emails) == 0 {
			id, err := password.SetupRecoverySuccess(i, in.Email)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
			}
			var sb strings.Builder
			fmt.Fprintf(&sb, "%s?t=%s", route.RecoverPasswordSuccess, id)
			path := sb.String()
			c.Append("HX-Redirect", path)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
		}

		emailAddresses := []string{}
		for i := 0; i < len(emails); i++ {
			email := emails[i]
			emailAddresses = append(emailAddresses, email.Address)
		}

		if !slices.Contains(emailAddresses, in.Email) {
			id, err := password.SetupRecoverySuccess(i, in.Email)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
			}
			var sb strings.Builder
			fmt.Fprintf(&sb, "%s?t=%s", route.RecoverPasswordSuccess, id)
			path := sb.String()
			c.Append("HX-Redirect", path)
			return nil
		}

		err = password.SetupRecovery(i, p.ID, in.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
		}

		id, err := password.SetupRecoverySuccess(i, in.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindRecoverPasswordErrInternal, layout.None)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s?t=%s", route.RecoverPasswordSuccess, id)
		path := sb.String()
		c.Append("HX-Redirect", path)
		return nil
	}
}

func ResetPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tid := c.Query("t")
		if len(tid) == 0 {
			return c.Redirect(route.Home)
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

func ResetPassword(i *service.Interfaces) fiber.Handler {
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
