package handlers

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/util"
)

func ChangePassword(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Current         string `form:"current"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirm"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if in.Password != in.ConfirmPassword {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayer(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This is a catastrophic failure; a Player object doesn't exist for a logged-in player
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		ok, err := password.Verify(in.Current, p.PwHash)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if !ok {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if !password.IsValid(in.Password) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		b := fiber.Map{
			"NoticeSectionID": "profile-password-notice",
			"SectionClass":    "pt-4",
			"NoticeText": []string{
				"Success!",
				"Your password has been changed.",
			},
		}

		return c.Render(partials.NoticeSectionSuccess, b, layouts.None)
	}
}
