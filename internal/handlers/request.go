package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
)

func CreateRequestComment(i *shared.Interfaces) fiber.Handler {
	type request struct {
		Text  string `form:"text"`
		Field string `form:"text"`
		CID   int64  `form:"cid"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(bind.Name), "views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		req, err := i.Queries.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.PID != pid {
			// TODO: Check permissions on the commenting player to see if they can comment on this application
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		// TODO: Sanitize and validate the contents of the comment
		r := new(request)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if r.CID == 0 {
			if len(r.Field) > 0 {
				_, err := i.Queries.AddCommentToRequestField(context.Background(), queries.AddCommentToRequestFieldParams{
					RID: rid,
					VID: req.VID,
					PID: pid.(int64),
					// TODO: Validate this field based on the type of the request
					// TODO: Keep a list of valid fields by request type?
					Field: r.Field,
					Text:  r.Text,
				})
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				c.Status(fiber.StatusCreated)
				// TODO: Return the comment markup here
				return nil
			} else {
				_, err := i.Queries.AddCommentToRequest(context.Background(), queries.AddCommentToRequestParams{
					RID:  rid,
					PID:  pid.(int64),
					Text: r.Text,
				})
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				c.Status(fiber.StatusCreated)
				// TODO: Return the comment markup here
				return nil
			}
		} else {
			rc, err := i.Queries.GetRequestComment(context.Background(), r.CID)
			if err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					return nil
				}
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			_, err = i.Queries.AddReplyToFieldComment(context.Background(), queries.AddReplyToFieldCommentParams{
				RID:   rid,
				PID:   pid.(int64),
				Field: rc.Field,
				Text:  r.Text,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusCreated)
			// TODO: Return the comment markup here
			return nil

		}
	}
}
