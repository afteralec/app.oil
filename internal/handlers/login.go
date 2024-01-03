package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
	"petrichormud.com/app/internal/views"
)

func Login(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	return func(c *fiber.Ctx) error {
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.PartialLoginErr, &fiber.Map{}, "")
		}

		p, err := i.Queries.GetPlayerByUsername(context.Background(), r.Username)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.PartialLoginErr, &fiber.Map{}, "")
		}

		v, err := password.Verify(r.Password, p.PwHash)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.PartialLoginErr, &fiber.Map{}, "")
		}
		if !v {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.PartialLoginErr, &fiber.Map{}, "")
		}

		pid := p.ID
		err = username.Cache(i.Redis, pid, p.Username)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.PartialLoginErr, &fiber.Map{}, "")
		}

		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.PartialLoginErr, &fiber.Map{}, "")
		}

		sess.Set("pid", pid)
		if err = sess.Save(); err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.PartialLoginErr, &fiber.Map{}, "")
		}

		c.Append("HX-Refresh", "true")
		return nil
	}
}

func LoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid != nil {
			return c.Redirect(routes.Home)
		}

		return c.Render(views.Login, c.Locals(constants.BindName), views.LayoutStandalone)
	}
}
