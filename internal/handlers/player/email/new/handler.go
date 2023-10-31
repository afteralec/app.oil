package newemail

import (
	"context"
	"database/sql"
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/queries"
)

const MaxEmailCount = 3

type NewEmailInput struct {
	Email string `form:"email"`
}

func New(db *sql.DB, s *session.Store, q *queries.Queries, r *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		emailCount, err := q.CountPlayerEmails(context.Background(), pid.(int64))
		if err != nil {
			log.Print(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if emailCount >= MaxEmailCount {
			c.Status(400)
			return nil
		}

		i := new(NewEmailInput)
		if err := c.BodyParser(i); err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		result, err := q.CreatePlayerEmail(
			context.Background(),
			queries.CreatePlayerEmailParams{Pid: pid.(int64), Email: i.Email},
		)
		if err != nil {
			log.Print(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		return c.Render("web/views/partials/profile/email/unverified-email", &fiber.Map{
			"ID":    id,
			"Email": i.Email,
		}, "")
	}
}
