package profile

import (
	"context"
	"log"
	"slices"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
)

func New(q *queries.Queries, r *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			return c.Redirect("/")
		}

		perms := c.Locals("perms")
		if perms == nil {
			return c.Redirect("/")
		}

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			log.Print(err)
			return c.Redirect("/")
		}

		if id != pid && !slices.Contains(perms.([]string), permissions.ViewPlayer) {
			return c.Redirect("/")
		}

		b := c.Locals("bind").(fiber.Map)
		b["ID"] = id

		return c.Render("web/views/profile", b)
	}
}

func NewWithoutParams(q *queries.Queries, r *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			return c.Redirect("/")
		}

		perms := c.Locals("perms")
		if perms == nil {
			return c.Redirect("/")
		}

		b := c.Locals("bind").(fiber.Map)

		emails, err := q.ListPlayerEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b["Emails"] = emails
		b["VerifiedEmails"] = email.Verified(emails)
		b["GravatarEmail"] = "othertest@quack.ninja"
		b["GravatarHash"] = email.GravatarHash("after.alec@gmail.com")

		return c.Render("web/views/profile", b)
	}
}
