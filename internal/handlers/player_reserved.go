package handlers

import (
	"context"
	"log"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
)

type Input struct {
	Username string `form:"username"`
}

func PlayerReserved(c *fiber.Ctx) error {
	i := new(Input)
	if err := c.BodyParser(i); err != nil {
		return err
	}

	ctx := context.Background()
	u, err := queries.Q.GetPlayerUsername(ctx, i.Username)
	// TODO: Figure out getting the error code from these
	if err != nil {
		log.Print(err)
		return c.Render("web/views/htmx/player-free", fiber.Map{}, "")
	}

	if i.Username == u {
		return c.Render("web/views/htmx/player-reserved", fiber.Map{}, "")
	} else {
		return c.Render("web/views/htmx/player-free", fiber.Map{}, "")
	}
}
