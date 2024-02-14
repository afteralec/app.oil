package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"slices"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/username"
	"petrichormud.com/app/internal/view"
)

func RecoverPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(view.RecoverPassword, view.Bind(c), layout.Standalone)
	}
}

func RecoverPasswordSuccessPage(i *interfaces.Shared) fiber.Handler {
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

func RecoverPassword(i *interfaces.Shared) fiber.Handler {
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
