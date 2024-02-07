package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"slices"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/username"
	"petrichormud.com/app/internal/views"
)

func RecoverPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.RecoverPassword, views.Bind(c), layouts.Standalone)
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
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		b := views.Bind(c)
		b["EmailAddress"] = email
		return c.Render(views.RecoverPasswordSuccess, b, layouts.Standalone)
	}
}

func RecoverPassword(i *interfaces.Shared) fiber.Handler {
	type request struct {
		Username string `form:"username"`
		Email    string `form:"email"`
	}
	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
		}

		v := username.IsValid(r.Username)
		if !v {
			c.Status(fiber.StatusBadRequest)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInvalidUsername, layouts.None)
		}

		_, err := mail.ParseAddress(r.Email)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInvalidEmail, layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), r.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				id, err := password.SetupRecoverySuccess(i, r.Email)
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					c.Append(constants.HeaderHXAcceptable, "true")
					return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
				}
				var sb strings.Builder
				fmt.Fprintf(&sb, "%s?t=%s", routes.RecoverPasswordSuccess, id)
				path := sb.String()
				c.Append("HX-Redirect", path)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
		}

		emails, err := qtx.ListVerifiedEmails(context.Background(), p.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
		}
		if len(emails) == 0 {
			id, err := password.SetupRecoverySuccess(i, r.Email)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(constants.HeaderHXAcceptable, "true")
				return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
			}
			var sb strings.Builder
			fmt.Fprintf(&sb, "%s?t=%s", routes.RecoverPasswordSuccess, id)
			path := sb.String()
			c.Append("HX-Redirect", path)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
		}

		emailAddresses := []string{}
		for i := 0; i < len(emails); i++ {
			email := emails[i]
			emailAddresses = append(emailAddresses, email.Address)
		}

		if !slices.Contains(emailAddresses, r.Email) {
			id, err := password.SetupRecoverySuccess(i, r.Email)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(constants.HeaderHXAcceptable, "true")
				return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
			}
			var sb strings.Builder
			fmt.Fprintf(&sb, "%s?t=%s", routes.RecoverPasswordSuccess, id)
			path := sb.String()
			c.Append("HX-Redirect", path)
			return nil
		}

		err = password.SetupRecovery(i, p.ID, r.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
		}

		id, err := password.SetupRecoverySuccess(i, r.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(constants.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindRecoverPasswordErrInternal, layouts.None)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s?t=%s", routes.RecoverPasswordSuccess, id)
		path := sb.String()
		c.Append("HX-Redirect", path)
		return nil
	}
}
