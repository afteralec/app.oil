package handlers

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

type Player struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func NewPlayer(c *fiber.Ctx) error {
	p := new(Player)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	return c.Render("web/views/index", fiber.Map{
		"CopyrightYear": time.Now().Year(),
		"MetaContent":   "Petrichor MUD - a modern take on a classic MUD style of game.",
		"Title":         "Sup",
	})
}
