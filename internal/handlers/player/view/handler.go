package viewplayer

import (
	fiber "github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			return c.Redirect("/")
		}

		id := c.Params("id")
		b := c.Locals("bind").(fiber.Map)
		b["ID"] = id

		return c.Render("web/views/player", b)
	}
}
