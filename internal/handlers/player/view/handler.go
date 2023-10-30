package viewplayer

import (
	"log"
	"slices"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"

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

		return c.Render("web/views/player", b)
	}
}

type PlayerEmail struct {
	Email    string
	Verified bool
	ID       int64
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

		emails := []PlayerEmail{
			{ID: 1, Email: "test@test.com", Verified: false},
			{ID: 2, Email: "othertest@quack.ninja", Verified: true},
			{ID: 3, Email: "tests@testes.com", Verified: true},
			{ID: 4, Email: "quack@test.ninja", Verified: false},
			{ID: 5, Email: "ninja@quack.test", Verified: false},
		}
		b["Emails"] = emails

		if len(emails) == 0 {
			b["NoEmails"] = true
		}

		return c.Render("web/views/player", b)
	}
}
