package newplayer

import (
	"context"
	"database/sql"
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/username"
)

type Player struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func New(db *sql.DB, s *session.Store, q *queries.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		p := new(Player)

		if err := c.BodyParser(p); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		u := username.Sanitize(p.Username)

		if !username.Validate(u) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !password.Validate(p.Password) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		pw_hash, err := password.Hash(p.Password)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := db.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()

		qtx := q.WithTx(tx)
		ctx := context.Background()
		result, err := qtx.CreatePlayer(ctx, queries.CreatePlayerParams{
			Username: u,
			PwHash:   pw_hash,
		})
		if err != nil {
			c.Status(fiber.StatusConflict)
			return nil
		}
		err = tx.Commit()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		pid, err := result.LastInsertId()
		if err != nil {
			log.Print(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		sess, err := s.Get(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		sess.Set("pid", pid)
		if err = sess.Save(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		return nil
	}
}
