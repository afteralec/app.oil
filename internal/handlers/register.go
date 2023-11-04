package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/roles"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

const RegisterRoute = "/player/new"

type Player struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Register(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		p := new(Player)

		if err := c.BodyParser(p); err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/register/err-internal", c.Locals("bind"), "")
		}

		u := username.Sanitize(p.Username)

		// TODO: Return the reason the username is invalid
		if !username.Validate(u) {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("web/views/partials/register/err-invalid", c.Locals("bind"), "")
		}

		// TODO: Return the reason the password is invalid
		if !password.Validate(p.Password) {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("web/views/partials/register/err-invalid", c.Locals("bind"), "")
		}

		pw_hash, err := password.Hash(p.Password)
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/register/err-internal", c.Locals("bind"), "")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/register/err-internal", c.Locals("bind"), "")
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)
		result, err := qtx.CreatePlayer(
			context.Background(),
			queries.CreatePlayerParams{
				Username: u,
				Role:     roles.Player,
				PwHash:   pw_hash,
			},
		)
		if err != nil {
			// TODO: Distinguish between "already exists" and a connection error
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render("web/views/partials/register/err-conflict", c.Locals("bind"), "")
		}

		pid, err := result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/register/err-internal", c.Locals("bind"), "")
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/register/err-internal", c.Locals("bind"), "")
		}

		username.Cache(i.Redis, pid, p.Username)

		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/register/err-internal", c.Locals("bind"), "")
		}

		sess.Set("pid", pid)
		if err = sess.Save(); err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/register/err-internal", c.Locals("bind"), "")
		}

		c.Append("HX-Trigger-After-Swap", "ptrcr:register-success")
		c.Status(fiber.StatusCreated)
		return nil
	}
}
