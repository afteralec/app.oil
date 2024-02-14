package handler

import (
	"context"
	"database/sql"
	"net/mail"
	"strconv"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/player/username"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

func AddEmail(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Email string `form:"email"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrUnauthorized, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		ec, err := qtx.CountEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}

		if ec >= email.MaxCount {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrTooMany(), layout.None)
		}

		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInvalid, layout.None)
		}

		e, err := mail.ParseAddress(in.Email)
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInvalid, layout.None)
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}
		if err == nil && ve.Verified {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrConflict(e.Address), layout.None)
		}

		result, err := qtx.CreateEmail(
			context.Background(),
			query.CreateEmailParams{PID: pid.(int64), Address: e.Address},
		)
		if err != nil {
			if me, ok := err.(*mysql.MySQLError); ok {
				if me.Number == mysqlerr.ER_DUP_ENTRY {
					c.Append("HX-Retarget", "#add-email-error")
					c.Append("HX-Reswap", "outerHTML")
					c.Append(header.HXAcceptable, "true")
					c.Status(fiber.StatusConflict)
					return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrConflict(e.Address), layout.None)
				}
			}
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}

		if err = tx.Commit(); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}

		if err = email.SendVerificationEmail(i, id, e.Address); err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}

		c.Status(fiber.StatusCreated)
		return c.Render(partial.ProfileEmailNew, &fiber.Map{
			"ID":      id,
			"Address": e.Address,
			"Created": true,
		}, "")
	}
}

func EditEmail(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Email string `form:"email"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInvalid, layout.None)
		}

		ne, err := mail.ParseAddress(in.Email)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "innerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInvalid, layout.None)
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrUnauthorized, layout.None)
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInvalid, layout.None)
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInvalid, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "#profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
			}
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		if e.PID != pid {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		if !e.Verified {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		if e.Address == ne.Address {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrConflictSame(ne.Address), layout.None)
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), ne.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileAddEmailErrInternal, layout.None)
		}
		if err == nil && ve.Verified {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrConflict(ne.Address), layout.None)
		}

		if err := qtx.DeleteEmail(context.Background(), id); err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "#profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
			}
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		result, err := qtx.CreateEmail(context.Background(), query.CreateEmailParams{
			Address: ne.Address,
			PID:     pid,
		})
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		id, err = result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "#profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileEditEmailErrInternal, layout.None)
		}

		return c.Render(partial.ProfileEmailUnverified, &fiber.Map{
			"ID":       id,
			"Address":  ne.Address,
			"Verified": false,
		}, "")
	}
}

func DeleteEmail(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrUnauthorized, layout.None)
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
			}
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		if e.PID != pid.(int64) {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		if err := qtx.DeleteEmail(context.Background(), id); err != nil {
			if err == sql.ErrNoRows {
				c.Append("HX-Retarget", "profile-email-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
			}
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "profile-email-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partial.NoticeSectionError, partial.BindProfileDeleteEmailErrInternal, layout.None)
		}

		return c.Render(partial.ProfileEmailDeleteSuccess, &fiber.Map{
			"ID":      e.ID,
			"Address": e.Address,
		}, "")
	}
}

func VerifyEmailPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		token := c.Query("t")
		key := email.VerificationKey(token)

		exists, err := i.Redis.Exists(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		if exists != 1 {
			c.Status(fiber.StatusNotFound)
			b := view.Bind(c)
			b["NotFoundMessage"] = "Sorry, it looks like this link has expired."
			b["NotFoundButtonLink"] = route.Profile
			b["NotFoundButtonText"] = "Return to Profile"
			return c.Render(view.NotFound, b, layout.Standalone)
		}

		eid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			if err == redis.Nil {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if e.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		ve, err := qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil && err != sql.ErrNoRows {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		if err == nil && ve.Verified {
			c.Status(fiber.StatusConflict)
			b := view.Bind(c)
			b["ErrMessageConflict"] = "That email has already been verified."
			return c.Render(view.Conflict, b, layout.Standalone)
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		un, err := username.Get(i, pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		b := view.Bind(c)
		b["VerifyToken"] = c.Query("t")
		b["Address"] = e.Address
		b["Username"] = un
		return c.Render(view.VerifyEmail, b, layout.Standalone)
	}
}

func VerifyEmail(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Refresh", "true")
			return nil
		}

		token := c.Query("t")
		if len(token) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		key := email.VerificationKey(token)
		eid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			if err == redis.Nil {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		e, err := i.Queries.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if e.Verified {
			c.Status(fiber.StatusConflict)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = i.Queries.MarkEmailVerified(context.Background(), id); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return c.Render(partial.VerifyEmailSuccess, &fiber.Map{}, "")
	}
}
