package newemail

import (
	"context"
	"database/sql"
	"log"
	"net/mail"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/email"
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

		ec, err := q.CountPlayerEmails(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if ec >= MaxEmailCount {
			c.Append("HX-Retarget", "#add-email-error")
			return c.Render("web/views/partials/profile/email/err-too-many-emails", &fiber.Map{}, "")
		}

		i := new(NewEmailInput)
		if err := c.BodyParser(i); err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		e, err := mail.ParseAddress(i.Email)
		if err != nil {
			c.Append("HX-Retarget", "#add-email-error")
			return c.Render("web/views/partials/profile/email/err-invalid-email", &fiber.Map{}, "")
		}

		result, err := q.CreatePlayerEmail(
			context.Background(),
			queries.CreatePlayerEmailParams{Pid: pid.(int64), Email: e.Address},
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

		err = email.Verify(r, id, e.Address)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		return c.Render("web/views/partials/profile/email/new-email", &fiber.Map{
			"ID":    id,
			"Email": e.Address,
		}, "")
	}
}
