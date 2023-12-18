package handlers

import (
	"context"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

func Register(i *shared.Interfaces) fiber.Handler {
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
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
		}

		u := username.Sanitize(p.Username)

		if !username.IsValid(u) {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("views/partials/register/err-invalid", c.Locals(bind.Name), "")
		}

		if !password.IsValid(p.Password) {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("views/partials/register/err-invalid", c.Locals(bind.Name), "")
		}

		if p.Password != p.ConfirmPassword {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("views/partials/register/err-invalid", c.Locals(bind.Name), "")
		}

		pwHash, err := password.Hash(p.Password)
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)
		result, err := qtx.CreatePlayer(
			context.Background(),
			queries.CreatePlayerParams{
				Username: u,
				PwHash:   pwHash,
			},
		)
		if err != nil {
			me, ok := err.(*mysql.MySQLError)
			if !ok {
				c.Append("HX-Retarget", "#register-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
			}
			if me.Number == mysqlerr.ER_DUP_ENTRY {
				c.Append("HX-Retarget", "#register-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusConflict)
				return c.Render("views/partials/register/err-conflict", c.Locals(bind.Name), "")
			}
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
		}

		pid, err := result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
		}

		username.Cache(i.Redis, pid, p.Username)

		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
		}

		sess.Set("pid", pid)
		if err = sess.Save(); err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/partials/register/err-internal", c.Locals(bind.Name), "")
		}

		c.Append("HX-Trigger-After-Swap", "ptrcr:register-success")
		c.Status(fiber.StatusCreated)
		return nil
	}
}
