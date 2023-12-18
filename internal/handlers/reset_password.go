package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

func ResetPasswordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tid := c.Query("t")
		if len(tid) == 0 {
			return c.Redirect("/")
		}

		b := c.Locals(bind.Name).(fiber.Map)
		b["ResetPasswordToken"] = tid

		return c.Render("views/reset/password", b, "views/layouts/standalone")
	}
}

func ResetPasswordSuccessPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("views/reset/password/success", c.Locals(bind.Name), "views/layouts/standalone")
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
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		vu := username.IsValid(r.Username)
		if !vu {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if r.Password != r.ConfirmPassword {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		vp := password.IsValid(r.Password)
		if !vp {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		tid := c.Query("t")
		if len(tid) == 0 {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		key := password.RecoveryKey(tid)
		rpid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		pid, err := strconv.ParseInt(rpid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		p, err := i.Queries.GetPlayer(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusUnauthorized)
				return nil
			}
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if p.Username != r.Username {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		pwHash, err := password.Hash(r.Password)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		_, err = i.Queries.UpdatePlayerPassword(context.Background(), queries.UpdatePlayerPasswordParams{
			ID:     pid,
			PwHash: pwHash,
		})
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		c.Append("HX-Redirect", routes.ResetPasswordSuccess)
		return nil
	}
}
