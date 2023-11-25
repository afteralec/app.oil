package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

const (
	ResetPasswordRoute        = "/reset/password"
	ResetPasswordSuccessRoute = "/reset/password/route"
)

func ResetPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tid := c.Query("t")
		if len(tid) == 0 {
			return c.Redirect("/")
		}

		b := c.Locals("bind").(fiber.Map)
		b["ResetPasswordToken"] = tid

		return c.Render("web/views/reset/password", b, "web/views/layouts/standalone")
	}
}

func ResetPasswordSuccessPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("web/views/reset/password/success", c.Locals("bind"), "web/views/layouts/standalone")
	}
}

func ResetPassword(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Username        string `form:"username"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirmPassword"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		vu := username.Validate(r.Username)
		if !vu {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if r.Password != r.ConfirmPassword {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		vp := password.Validate(r.Password)
		if !vp {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tid := c.Query("t")
		if len(tid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		key := password.RecoveryKey(tid)
		rpid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		pid, err := strconv.ParseInt(rpid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		p, err := i.Queries.GetPlayer(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if p.Username != r.Username {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		pwHash, err := password.Hash(r.Password)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Cache a new UUID here that represents a reference to an email to show on the success page

		_, err = i.Queries.UpdatePlayerPassword(context.Background(), queries.UpdatePlayerPasswordParams{
			ID:     pid,
			PwHash: pwHash,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", ResetPasswordSuccessRoute)
		return nil
	}
}
