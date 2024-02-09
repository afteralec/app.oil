package handler

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/util"
)

// TODO: This and the markup allow entering your current password
// This should be a failure case
func ChangePassword(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Current         string `form:"current"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirm"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Something's gone terribly wrong.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}

		if in.Password != in.ConfirmPassword {
			c.Status(fiber.StatusBadRequest)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"The new password and password confirmation doesn't match.",
					"Please try again.",
				},
				"NoticeIcon": true,
			}
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Your session has expired.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Something's gone terribly wrong.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayer(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This is a catastrophic failure; a Player object doesn't exist for a logged-in player
				c.Status(fiber.StatusInternalServerError)
				b := fiber.Map{
					"NoticeSectionID": "profile-password-notice",
					"SectionClass":    "pt-4",
					"NoticeText": []string{
						"Something's gone terribly wrong.",
					},
					"NoticeIcon":    true,
					"RefreshButton": true,
				}
				return c.Render(partials.NoticeSectionError, b, layouts.None)
			}
			c.Status(fiber.StatusInternalServerError)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Something's gone terribly wrong.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}

		ok, err := password.Verify(in.Current, p.PwHash)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"Something's gone terribly wrong.",
				},
				"NoticeIcon":    true,
				"RefreshButton": true,
			}
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}
		if !ok {
			c.Status(fiber.StatusUnauthorized)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"The current password you entered isn't correct.",
				},
				"NoticeIcon": true,
			}
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}

		if !password.IsValid(in.Password) {
			c.Status(fiber.StatusBadRequest)
			b := fiber.Map{
				"NoticeSectionID": "profile-password-notice",
				"SectionClass":    "pt-4",
				"NoticeText": []string{
					"What you entered isn't a valid password.",
					"Please follow the prompts and try again.",
				},
				"NoticeIcon": true,
			}
			return c.Render(partials.NoticeSectionError, b, layouts.None)
		}

		b := fiber.Map{
			"NoticeSectionID": "profile-password-notice",
			"SectionClass":    "pt-4",
			"NoticeText": []string{
				"Success!",
				"Your password has been changed.",
			},
			"NoticeIcon": true,
		}

		return c.Render(partials.NoticeSectionSuccess, b, layouts.None)
	}
}
