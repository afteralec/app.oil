package verify

import (
	"context"
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
			// TODO: Here, the user isn't logged in - return a login page that redirects to or hx-boosts the verify page
			return c.Redirect("/")
		}

		perms, err := permissions.List(q, r, pid.(int64))
		if err != nil {
			return c.Redirect("/")
		}
		if !slices.Contains(perms, permissions.AddEmail) {
			return c.Redirect("/")
		}

		// TODO: Here, the user is logged in and has perms - send the page that has the function for verifying your email
		return c.Redirect("/")
	}
}

func NewVerify(q *queries.Queries, r *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := permissions.List(q, r, pid.(int64))
		if err != nil {
			return c.Redirect("/")
		}
		if !slices.Contains(perms, permissions.AddEmail) {
			return c.Redirect("/")
		}

		key := c.Params("t")
		if len(key) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		eid, err := r.Get(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = q.MarkEmailVerified(context.Background(), id)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = r.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return nil
	}
}
