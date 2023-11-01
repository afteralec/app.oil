package verify

import (
	"context"
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
		token := c.Query("t")
		exists, err := r.Exists(context.Background(), token).Result()
		if err != nil {
			return c.Redirect("/")
		}
		if exists != 1 {
			return c.Redirect("/")
		}

		pid := c.Locals("pid")
		if pid == nil {
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		perms, err := permissions.List(q, r, pid.(int64))
		if err != nil {
			return c.Redirect("/")
		}
		if !slices.Contains(perms, permissions.AddEmail) {
			return c.Redirect("/")
		}

		b := c.Locals("bind").(fiber.Map)
		b["VerifyToken"] = c.Query("t")
		log.Print(b["VerifyToken"])

		return c.Render("web/views/verify", b, "web/views/layouts/standalone")
	}
}

func NewVerify(q *queries.Queries, r *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			// TODO: This should redirect them back to the login page for this token
			return nil
		}

		perms, err := permissions.List(q, r, pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if !slices.Contains(perms, permissions.AddEmail) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		key := c.Query("t")
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

		return c.Render("web/views/partials/verify/success", &fiber.Map{}, "")
	}
}
